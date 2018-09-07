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

	"github.com/spf13/cobra"
	"github.com/wangzn/dnspodapi"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "info is for some basic test",
	Long:  ``,
	Run:   runInfoCmd,
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("version", "", "print the api version")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runInfoCmd(cmd *cobra.Command, args []string) {
	r := InfoActionRunner{}
	r.Run()
}

// InfoActionRunner defines the runner to run info action
type InfoActionRunner struct {
	Action   string
	Params   map[string]string
	APIID    int
	APIToken string
}

// Run starts to run action
func (r *InfoActionRunner) Run() {
	r.run()
}

// Name returns runner name
func (r *InfoActionRunner) Name() string {
	return fmt.Sprintf("Info running...")
}

// Detail returns detail information
func (r *InfoActionRunner) Detail() string {
	return fmt.Sprintf("Get API version info")
}

func (r *InfoActionRunner) run() {
	ver, err := dnspodapi.GetVersion()
	if err != nil {
		pe(err)
	} else {
		fmt.Printf("\n\tOK. API version: %s \n", ver)
	}
}
