// Copyright 2016 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vaultinit

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

var roots *x509.CertPool
var log = logrus.New()
var keyfile string
var certfile string

func initLog() {
	fmt.Println("entering initLog")
	level, err := logrus.ParseLevel(viper.GetString("log-level"))
	fmt.Println("parse level: " + level.String())
	log.Level = level
	if err != nil {
		log.Fatalln(err)
	}
}

func createTempCerts() {

	keyfile = viper.GetString("tmp-cert-dir") + "/key.tls"
	certfile = viper.GetString("tmp-cert-dir") + "/cert.tls"

	//init rootca
	roots = x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(parsePEMcert(viper.GetString("vault-controller-cacert")))
	if !ok {
		log.Fatalln("failed to parse root certificate from $VAULT_CONTROLLER_CACERT: %s", viper.GetString("vault-controller-cacert"))
	}

	//init temp certificates
	log.Debugln("initializing temp certificates")
	createSelfSignedCertificate(viper.GetString("pod-ip"), keyfile, certfile)
	log.Debugln("temp certificates initialized")

}

func validateConfig() {
	log.Debugln("validating configuration")
	if viper.GetString("pod-name") == "" {
		log.Fatalln("pod-name must be set and non-empty")
	}
	if viper.GetString("pod-namespace") == "" {
		log.Fatalln("pod-namespace must be set and non-empty")
	}
	if viper.GetString("pod-ip") == "" {
		log.Fatalln("pod-ip must be set and non-empty")
	}
	if viper.GetString("creds-init-port") == "" {
		log.Fatalln("vault-controller-port must be set and non-empty")
	}
	if viper.GetString("vault-cacert") == "" {
		log.Fatalln("vault-cacert must be set and non-empty")
	}
	if viper.GetString("vault-addr") == "" {
		log.Fatalln("vault-addr must be set and non-empty")
	}
	if viper.GetString("vault-controller-addr") == "" {
		log.Fatalln("vault-controller-addr must be set and non-empty")
	}
	if viper.GetString("vault-controller-cacert") == "" {
		log.Fatalln("vault-controller-cacert must be set and non-empty")
	}
	log.Debugln("configuration is valid")
}

func RunInitCreds() {
	initLog()
	log.Infoln("Starting vault-init...")

	validateConfig()
	createTempCerts()

	http.Handle("/", tokenHandler{viper.GetString("vault-addr")})
	go func() {

		u := fmt.Sprintf("0.0.0.0:%s", viper.GetString("vault-init-port"))
		log.Debugf("starting to listen at: %s", u)
		log.Fatalln(http.ListenAndServeTLS(u, certfile, keyfile, nil))
	}()

	// Ensure the token handler is ready.
	time.Sleep(time.Millisecond * 300)

	tokenFile := viper.GetString("creds-file")

	// Remove exiting token files before requesting a new one.
	if err := os.Remove(tokenFile); err != nil {
		log.WithFields(logrus.Fields{
			"tokenFile": tokenFile,
			"err":       err,
		}).Infoln("could not remove token file")
	}

	// Set up a file watch on the wrapped vault token.
	tokenWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Fatalln("could not create watcher")
	}
	err = tokenWatcher.Add(path.Dir(tokenFile))
	if err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Fatalln("could not add watcher")
	}

	done := make(chan bool)
	retryDelay := 5 * time.Second
	go func() {
		for {
			log.Debugf("requesting token at: %s with name: %s and namespace %s", viper.GetString("vault-controller-addr"), viper.GetString("pod-name"), viper.GetString("pod-namespace"))
			err := requestToken(viper.GetString("vault-controller-addr"), viper.GetString("pod-name"), viper.GetString("pod-namespace"))
			if err != nil {
				log.Infof("token request: Request error %v; retrying in %v", err, retryDelay)
				time.Sleep(retryDelay)
				continue
			}
			log.Infoln("Token request complete; waiting for callback...")
			select {
			case <-time.After(time.Second * 30):
				log.Infoln("token request: Timeout waiting for callback")
				break
			case <-tokenWatcher.Events:
				tokenWatcher.Close()
				time.Sleep(200 * time.Millisecond)
				close(done)
				return
			case err := <-tokenWatcher.Errors:
				log.Infof("token request: error watching the token file", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		log.Infoln("Shutdown signal received, exiting...")
	case <-done:
		log.Infoln("Successfully obtained and unwrapped the vault token, exiting...")
	}
}

func requestToken(vaultControllerAddr, name, namespace string) error {
	u := fmt.Sprintf("%s/token?name=%s&namespace=%s", vaultControllerAddr, name, namespace)
	log.Infof("Requesting a new wrapped token from %s", vaultControllerAddr)

	//crete tls client configuration
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: roots,
			},
		},
	}

	resp, err := client.Post(u, "", nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 202 {
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("%s", data)
}

func parsePEMcert(filename string) []byte {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("unable to read file %s; error %s", filename, err)
	}
	return bytes
}

// helper function to create a cert template with a serial number and other required fields
func certTemplate() (*x509.Certificate, error) {
	// generate a random serial number (a real cert authority would have some logic behind this)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.New("failed to generate serial number: " + err.Error())
	}

	tmpl := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"vaul init"}},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour), // valid for an hour
		BasicConstraintsValid: true,
	}
	return &tmpl, nil
}

func createSelfSignedCertificate(ip string, keyfile string, certfile string) {

	// code from here: https://ericchiang.github.io/post/go-tls/
	// generate a new key-pair
	rootKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("generating random key: %v", err)
	}

	rootCertTmpl, err := certTemplate()
	if err != nil {
		log.Fatalf("creating cert template: %v", err)
	}
	// describe what the certificate will be used for
	rootCertTmpl.IsCA = true
	rootCertTmpl.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	rootCertTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	rootCertTmpl.IPAddresses = []net.IP{net.ParseIP(ip)}

	rootCert, rootCertPEM, err := createCert(rootCertTmpl, rootCertTmpl, &rootKey.PublicKey, rootKey)
	if err != nil {
		log.Fatalf("error creating cert: %v", err)
	}
	log.Debugf("created temporary cert: %s", rootCertPEM)
	log.Debugf("created temporary cert signature: %#x", rootCert.Signature)

	rootKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rootKey),
	})

	err = ioutil.WriteFile(keyfile, rootKeyPEM, 0644)
	if err != nil {
		log.Fatalf("error writing key to file: ", keyfile, err)
	}
	err = ioutil.WriteFile(certfile, rootCertPEM, 0644)
	if err != nil {
		log.Fatalf("error writing cert to file: ", certfile, err)
	}

}

func createCert(template, parent *x509.Certificate, pub interface{}, parentPriv interface{}) (
	cert *x509.Certificate, certPEM []byte, err error) {

	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentPriv)
	if err != nil {
		return
	}
	// parse the resulting certificate so we can use it again
	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		return
	}
	// PEM encode the certificate (this is a standard TLS encoding)
	b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	certPEM = pem.EncodeToMemory(&b)
	return
}
