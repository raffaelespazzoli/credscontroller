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

package vaultcontroller

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"io"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var vaultClient *api.Client
var kubernetesClientSet *kubernetes.Clientset

var log = logrus.New()

func initLog() {
	var err error
	log.Level, err = logrus.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		log.Fatalln(err)
	}
}

func validateConfig() {
	log.Debugln("validating configuration")
	if viper.GetString("vault-token") == "" {
		log.Fatalln("vault-token must be set and non-empty")
	}
	if viper.GetString("vault-controller-key") == "" {
		log.Fatalln("vault-controller-key must be set and non-empty")
	}
	if viper.GetString("vault-controller-cert") == "" {
		log.Fatalln("vault-controller-cert must be set and non-empty")
	}
	if viper.GetString("vault-controller-port") == "" {
		log.Fatalln("vault-controller-port must be set and non-empty")
	}
	if viper.GetString("vault-cacert") == "" {
		log.Fatalln("vault-cacert must be set and non-empty")
	}
	if viper.GetString("vault-addr") == "" {
		log.Fatalln("vault-addr must be set and non-empty")
	}
	log.Debugln("configuration is valid")
}

func RunVaultController() {
	initLog()

	log.Infoln("Starting vault-controller app...")

	validateConfig()
	var err error

	// creates the in-cluster config
	log.Debugln("creating in cluster default kube client config")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalln(err)
	}
	log.WithFields(logrus.Fields{
		"config-host": config.Host,
	}).Debugln("kube client config created")

	//	 creates the clientset
	log.Debugln("creating kube client set")
	kubernetesClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}
	log.Debugln("kube client set created")

	//creates the vault config
	log.Debugln("creating vault config")
	vConfig := api.Config{
		Address: viper.GetString("vault-addr"),
	}
	tlsConfig := api.TLSConfig{
		CACert: viper.GetString("vault-cacert"),
	}
	err = vConfig.ConfigureTLS(&tlsConfig)
	if err != nil {
		log.Fatalln(err)
	}
	log.Debugln("created vault config")

	//creates the vault client
	log.Debugln("creating vault client")
	vaultClient, err = api.NewClient(&vConfig)
	if err != nil {
		log.Fatalln(err)
	}
	log.Debugln("created vault client")

	http.Handle("/token", handler{tokenRequestHandler})
	go func() {
		log.Fatalln(http.ListenAndServeTLS(viper.GetString("vault-controller-addr"), viper.GetString("vault-controller-cert"), viper.GetString("vault-controller-key"), nil))
	}()

	log.Infoln("Listening for token requests.")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infoln("Shutdown signal received, exiting...")
}

type handler struct {
	f func(io.Writer, *http.Request) (int, error)
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code, err := h.f(w, r)
	w.WriteHeader(code)
	if err != nil {
		log.Printf("%v", err)
		fmt.Fprintf(w, "%v", err)
	}
}
