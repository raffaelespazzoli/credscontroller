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
	"encoding/json"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
)

type tokenHandler struct {
	vaultAddr string
}

func (h tokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat(viper.GetString("creds-file"))
	if !os.IsNotExist(err) {
		log.Println("Token file already exists")
		w.WriteHeader(409)
		return
	}

	var swi api.SecretWrapInfo
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	r.Body.Close()

	err = json.Unmarshal(data, &swi)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

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
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	log.Debugln("created vault config")

	//creates the vault client
	log.Debugln("creating vault client")
	client, err := api.NewClient(&vConfig)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	log.Debugln("created vault client")

	client.SetToken(swi.Token)
	client.SetAddress(h.vaultAddr)

	// Vault knows to unwrap the client token if the token to unwrap is empty.
	secret, err := client.Logical().Unwrap("")
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	f, err := os.Create(viper.GetString("creds-file"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(&secret)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	log.Printf("wrote %s", viper.GetString("creds-file"))
	w.WriteHeader(200)
}
