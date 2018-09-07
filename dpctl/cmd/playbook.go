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
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/wangzn/dnspodapi"
)

var (
	playbookAct string
	scene       string
)

// playbookCmd represents the playbook command
var playbookCmd = &cobra.Command{
	Use:   "playbook",
	Short: "run pre-defined playbook",
	Long:  ``,
	Run:   runPlaybookCmd,
}

func init() {
	rootCmd.AddCommand(playbookCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playbookCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playbookCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	playbookCmd.PersistentFlags().StringVarP(&playbookAct, "action", "a", "",
		"playbook action: [ run | preview ]")

	playbookCmd.PersistentFlags().StringVarP(&scene, "scene", "s", "",
		"run the playbook scene")

}

func runPlaybookCmd(cmd *cobra.Command, args []string) {
	if pb, ok := cfg.Playbook[scene]; ok {
		pr := NewPlaybookRunner(playbookAct, scene, pb, cfg.Auth)
		pr.Run()
	} else {
		fmt.Printf("scene %s not found", scene)
		os.Exit(1)
	}
}

// PlaybookRunner defines the runner to run scene
type PlaybookRunner struct {
	Scene         string
	Action        string
	ActionRunners []ActionRunner
}

// NewPlaybookRunner returns a new PlaybookRunner pointer
func NewPlaybookRunner(action, scene string, cfg dnspodapi.Scene, auths dnspodapi.Auth,
) *PlaybookRunner {
	pbr := PlaybookRunner{
		Scene:  scene,
		Action: action,
	}
	ars := make([]ActionRunner, 0)
	for _, a := range cfg {
		var ar ActionRunner
		if _, ok := auths[a.Auth]; !ok {
			continue
		}
		auth := auths[a.Auth]
		switch a.Category {
		case "domain":
			ar = newDomainRunner(a, auth)
		case "record":
			ar = newRecordRunner(a, auth)
		case "info":
			ar = newInfoRunner(a, auth)
		default:
			log.Printf("invalid runner :%s", a.Subject)
		}
		if ar != nil {
			ars = append(ars, ar)
		}
	}
	pbr.ActionRunners = ars
	return &pbr
}

// Run starts to run playbook
func (p *PlaybookRunner) Run() {
	switch p.Action {
	case "run":
		p.run()
	case "preview":
		p.preview()
	default:
		fmt.Println("invalid playbook action")
		os.Exit(1)
	}
}

func (p *PlaybookRunner) run() {
	log.Printf("Playbook `%s` starts to run...", p.Scene)
	for _, ar := range p.ActionRunners {
		log.Printf("Action `%s` starts to run...", ar.Name())
		ar.Run()
	}
}

func (p *PlaybookRunner) preview() {
	fmt.Printf("\nPlaybook scene `%s` steps: \n\n", p.Scene)
	for i, ar := range p.ActionRunners {
		//spew.Dump(ar)
		//spew.Sdump(ar)
		// fmt.Printf("%#v", ar)
		fmt.Printf("\t%d: %s\n", i, ar.Detail())
	}
}

func newDomainRunner(a dnspodapi.ActionEntry, auth dnspodapi.AuthEntry,
) *DomainActionRunner {
	r := DomainActionRunner{
		Domain:   a.Subject,
		Action:   a.Action,
		Params:   a.Params,
		APIID:    auth.APIID,
		APIToken: auth.APIToken,
	}
	return &r
}

func newRecordRunner(a dnspodapi.ActionEntry, auth dnspodapi.AuthEntry,
) *RecordActionRunner {
	r := RecordActionRunner{
		Record:   a.Subject,
		Action:   a.Action,
		Params:   a.Params,
		APIID:    auth.APIID,
		APIToken: auth.APIToken,
	}
	return &r
}

func newInfoRunner(a dnspodapi.ActionEntry, auth dnspodapi.AuthEntry,
) *InfoActionRunner {
	r := InfoActionRunner{
		Action:   a.Action,
		Params:   a.Params,
		APIID:    auth.APIID,
		APIToken: auth.APIToken,
	}
	return &r
}
