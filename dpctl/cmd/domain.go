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
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/wangzn/dnspodapi"
	"github.com/wangzn/goutils/mymap"
)

const (
	// DomainActionCreate defines OPCreate
	DomainActionCreate = "create"
	// DomainActionRemove defines OPRemove
	DomainActionRemove = "remove"
	// DomainActionInfo defines OPInfo
	DomainActionInfo = "info"
	// DomainActionStatus defines OPStatus
	DomainActionStatus = "status"
	// DomainActionList defines OPList
	DomainActionList = "list"
)

var (
	domains   string
	domainAct string
)

// domainCmd represents the domain command
var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "domain is ctl tool for domain resources",
	Long:  ``,
	Run:   runDomainCmd,
}

func init() {
	rootCmd.AddCommand(domainCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// domainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	domainCmd.PersistentFlags().StringVarP(&domainAct, "action", "a", "list",
		"domain action: [ create | list | remove | info | enable | disable ]")

	domainCmd.Flags().StringVarP(&domains, "domain", "d", "",
		"domains to operate")

}

func runDomainCmd(cmd *cobra.Command, args []string) {
	r := DomainActionRunner{
		Domain:   domains,
		Action:   domainAct,
		APIID:    apiID,
		APIToken: apiToken,
	}
	r.Run()
}

// DomainActionRunner defines the runner to run domain action
type DomainActionRunner struct {
	Domain   string
	Action   string
	Params   map[string]string
	APIID    int
	APIToken string
}

// Run starts to run action
func (r *DomainActionRunner) Run() {
	r.run()
}

// Name returns runner name
func (r *DomainActionRunner) Name() string {
	return fmt.Sprintf("%s domain `%s`", strings.Title(r.Action), r.Domain)
}

// Detail returns detail information
func (r *DomainActionRunner) Detail() string {
	return r.Name()
}

func (r *DomainActionRunner) run() {
	dnspodapi.SetAPIToken(r.APIID, r.APIToken)
	r.checkDomainParams()
	switch r.Action {
	case "list":
		doListDomain()
	case "create":
		doCreateDomain(r.Domain)
	case "remove":
		doRemoveDomain(r.Domain)
	case "info":
		doInfoDomain(r.Domain)
	case "enable":
		doStatusDomain(r.Domain, "enable")
	case "disable":
		doStatusDomain(r.Domain, "disable")
	default:
		doListDomain()
	}
}

func (r *DomainActionRunner) checkDomainParams() {
	if r.Action != "list" {
		if r.Domain == "" {
			fmt.Println("domains is empty")
			os.Exit(1)
		}
	}
}

func doListDomain() {
	res, err := dnspodapi.GetDomainList()
	if err != nil {
		pe(err)
	}
	fmt.Println(dnspodapi.FormatDomainIDInts(res, format))
}

func doCreateDomain(dms string) {
	errs := make([]error, 0)
	ds := make([]dnspodapi.DomainEntry, 0)
	for _, d := range strings.Split(dms, ",") {
		de, err := dnspodapi.CreateDomain(d)
		if err != nil {
			errs = append(errs, err)
		}
		if de != nil {
			ds = append(ds, *de)
		}
	}
	if len(ds) > 0 {
		fmt.Println(dnspodapi.FormatDomains(ds, format))
	}
	pe(errs...)
}

func doRemoveDomain(dms string) {
	errs := make([]error, 0)
	res := make([][]string, 0)
	header := []string{"domain", "status", "msg"}
	for _, d := range strings.Split(dms, ",") {
		ok, err := dnspodapi.RemoveDomain(d)
		msg := ""
		if err != nil {
			errs = append(errs, err)
			msg = err.Error()
		}
		res = append(res, []string{d, fmt.Sprintf("%v", ok), msg})
	}
	fmt.Println(mymap.FormatSlices(header, res, format))
	pe(errs...)
}

func doInfoDomain(dms string) {
	errs := make([]error, 0)
	ds := make([]dnspodapi.DomainEntry, 0)
	for _, d := range strings.Split(dms, ",") {
		de, err := dnspodapi.GetDomainInfo(d)
		if err != nil {
			errs = append(errs, err)
		}
		if de != nil {
			ds = append(ds, *de)
		}
	}
	if len(ds) > 0 {
		fmt.Println(dnspodapi.FormatDomains(ds, format))
	}
	pe(errs...)
}

func doStatusDomain(dms, st string) {
	errs := make([]error, 0)
	res := make([][]string, 0)
	header := []string{"domain", "status", "msg"}
	for _, d := range strings.Split(dms, ",") {
		err := dnspodapi.SetDomainStatus(d, st)
		msg := ""
		ok := true
		if err != nil {
			errs = append(errs, err)
			msg = err.Error()
			ok = false
		}
		res = append(res, []string{d, fmt.Sprintf("%v", ok), msg})
	}
	fmt.Println(mymap.FormatSlices(header, res, format))
	pe(errs...)
}
