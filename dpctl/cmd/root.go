// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/wangzn/dnspodapi"
)

const (
	// APIIDField defines the field raw string of `api_id` in config file
	APIIDField = "api_id"
	// APITokenField defines the field raw string of `api_token` in config file
	APITokenField = "api_token"
	// LogOutput defines the output of writer of log
	LogOutput = "log_output"
	// GlobalNamespace defines the default namespace raw string
	GlobalNamespace = "global"
)

var (
	cfgFile   string
	format    string
	logOutput string
	cfg       *dnspodapi.Config
	err       error
)

var (
	apiID    int
	apiIDStr string
	apiToken string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pdctl",
	Short: "pdctl is a tool for operate DNS record of DNSPOD",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dnspod.yaml)")

	rootCmd.PersistentFlags().StringVar(&format, "format", "table",
		"output format: [ json | table ]")

	rootCmd.PersistentFlags().StringVar(&logOutput, "log", "null",
		"where to output log: [ null | stdout | stderr ]")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cfgFile = filepath.Join(home, ".dnspod.yaml")
	}
	cfg, err = dnspodapi.ParseYamlFile(cfgFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	initFormat(cfg.Global.Format)
	initLog(cfg.Global.LogOutput)
	initAuth(cfg.Global.Auth, cfg.Auth)
}

func initFormat(f string) {
	// cli format first
	if format == "" {
		format = f
	}
}

func initLog(l string) {
	if logOutput != "" {
		// cli log output first
		l = logOutput
	}
	switch l {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		log.SetOutput(ioutil.Discard)
	}
}

func initAuth(key string, auths dnspodapi.Auth) {
	if key == "" {
		key = "default"
	}
	if ae, ok := auths[key]; ok {
		apiID = ae.APIID
		apiToken = ae.APIToken
	}
	dnspodapi.SetAPIToken(apiID, apiToken)
}

// APIID returns apiID
func APIID() int {
	return apiID
}

// APIToken returns apiToken
func APIToken() string {
	return apiToken
}
