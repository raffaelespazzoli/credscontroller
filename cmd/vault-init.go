// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/raffaelespazzoli/credscontroller/vaultinit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// init-credsCmd represents the init-creds command
var initcredsCmd = &cobra.Command{
	Use:   "vault-init",
	Short: "vault-init retrieves credentials managed by the creds controller",
	Long:  `starts the credential retrieval process, currently supports integration with vault`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("calling vault-init with")
		fmt.Println("vault-addr: " + viper.GetString("vault-addr"))
		fmt.Println("vault-cacert: " + viper.GetString("vault-cacert"))
		fmt.Println("log-level: " + viper.GetString("log-level"))
		fmt.Println("pod-name: " + viper.GetString("pod-name"))
		fmt.Println("pod-ip: " + viper.GetString("pod-ip"))
		fmt.Println("pod-namespace: " + viper.GetString("pod-namespace"))
		fmt.Println("vault-controller-addr: " + viper.GetString("vault-controller-addr"))
		fmt.Println("vault-controller-cacert: " + viper.GetString("vault-controller-cacert"))
		fmt.Println("creds-init-port: " + viper.GetString("creds-init-port"))
		vaultinit.RunInitCreds()
	},
}

func init() {
	//fmt.Println("initcreds.init")
	RootCmd.AddCommand(initcredsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// init-credsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// init-credsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initcredsCmd.Flags().String("pod-name", "", "name of the pod running the creds initialization")
	initcredsCmd.Flags().String("pod-ip", "", "ip of the pod running the creds initialization")
	initcredsCmd.Flags().String("pod-namespace", "", "namespace of the pod running the creds initialization")
	initcredsCmd.Flags().String("vault-controller-addr", "https://vault-controller:8443", "address of the vault controller")
	initcredsCmd.Flags().String("vault-controller-cacert", "/var/run/secrets/kubernetes.io/serviceaccount/service-ca.crt", "ca certificate to be used when connecting to the vault controller")
	initcredsCmd.Flags().String("creds-init-port", "8443", "port at which to listen for vault controller incoming connections")
	initcredsCmd.Flags().String("tmp-cert-dir", "/tmp", "directory where temp certs will be created")
	initcredsCmd.Flags().String("creds-file", "/var/run/secrets/vaultproject.io/secret.json", "file where the credentials will be stored")
	viper.BindPFlag("pod-name", initcredsCmd.Flags().Lookup("pod-name"))
	viper.BindPFlag("pod-ip", initcredsCmd.Flags().Lookup("pod-ip"))
	viper.BindPFlag("pod-namespace", initcredsCmd.Flags().Lookup("pod-namespace"))
	viper.BindPFlag("vault-controller-addr", initcredsCmd.Flags().Lookup("vault-controller-addr"))
	viper.BindPFlag("vault-controller-cacert", initcredsCmd.Flags().Lookup("vault-controller-cacert"))
	viper.BindPFlag("creds-init-port", initcredsCmd.Flags().Lookup("creds-init-port"))
	viper.BindPFlag("tmp-cert-dir", initcredsCmd.Flags().Lookup("tmp-cert-dir"))
	viper.BindPFlag("creds-file", initcredsCmd.Flags().Lookup("creds-file"))

}
