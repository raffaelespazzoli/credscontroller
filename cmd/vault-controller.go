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
	"github.com/raffaelespazzoli/credscontroller/vaultcontroller"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

// vault-controllerCmd represents the vault-controller command
var vaultcontrollerCmd = &cobra.Command{
	Use:   "vault-controller",
	Short: "starts the creds controller with vault integration",
	Long:  "starts creds controller to receive credential requests with integration with vault",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Print: " + strings.Join(args, " "))
		fmt.Println("vaultcontroller called")
		vaultcontroller.RunVaultController()
	},
}

func init() {
	//fmt.Println("vaultcontroller.init")
	RootCmd.AddCommand(vaultcontrollerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// vault-controllerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// vault-controllerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	vaultcontrollerCmd.Flags().String("vault-token", "", "vault token")
	vaultcontrollerCmd.Flags().String("vault-wrap-ttl", "120", "vault wrap token time to live")
	vaultcontrollerCmd.Flags().String("vault-controller-key", "/var/run/secrets/kubernetes.io/certs/tls.key", "Vault-controller's private key for TLS connections")
	vaultcontrollerCmd.Flags().String("vault-controller-cert", "/var/run/secrets/kubernetes.io/certs/tls.crt", "Vault-controller's certificate for TLS connections")
	vaultcontrollerCmd.Flags().String("vault-controller-port", "8443", "Port at which the vault-controller will listen")
	viper.BindPFlag("vault-token", vaultcontrollerCmd.Flags().Lookup("vault-token"))
	viper.BindPFlag("vault-wrap-ttl", vaultcontrollerCmd.Flags().Lookup("vault-wrap-ttl"))
	viper.BindPFlag("vault-controller-key", vaultcontrollerCmd.Flags().Lookup("vault-controller-key"))
	viper.BindPFlag("vault-controller-cert", vaultcontrollerCmd.Flags().Lookup("vault-controller-cert"))
	viper.BindPFlag("vault-controller-port", vaultcontrollerCmd.Flags().Lookup("vault-controller-port"))

}
