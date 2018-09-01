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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/facebookgo/errgroup"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/wangzn/dnspodapi"
	"github.com/wangzn/goutils/structs"
)

const (
	// RecordActionCreate defines OPCreate
	RecordActionCreate = "create"
	// RecordActionRemove defines OPRemove
	RecordActionRemove = "remove"
	// RecordActionModify defines OPModify
	RecordActionModify = "modify"
)

var (
	clear      bool
	recordFile string
	zone       string
	act        string
	format     string
)

var (
	// DefaultRecordLineID defines the default record_line_id, 0 is "默认"
	DefaultRecordLineID = "0"
)

// recordCmd represents the record command
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "record is ctl tool for record resource",
	Long:  ``,
	Run:   runRecordCmd,
}

func init() {
	rootCmd.AddCommand(recordCmd)

	recordCmd.Flags().BoolVarP(&clear, "clear", "c", false, "clear existed record")

	recordCmd.Flags().StringVarP(&act, "action", "a", "list",
		"record actiom: [ create | list ]")

	recordCmd.Flags().StringVar(&format, "format", "table",
		"output format: [ json | table ]")

	recordCmd.Flags().StringVarP(&zone, "zone", "z", "",
		"zones for records, use ',' for multiple zones, e.g. 'abc.com,def.com'")

	recordCmd.Flags().StringVarP(&recordFile, "record_file", "f", "",
		"record file, each line contains 'record type value', e.g. 'www CNAME proxy'")
}

func runRecordCmd(cmd *cobra.Command, args []string) {
	if zone == "" {
		fmt.Println("zone is empty")
		os.Exit(1)
	}
	switch act {
	case "create":
		doCreateRecord()
	default:
		doListRecord()
	}
}

func doListRecord() {
	//res := make([]dnspodapi.RecordEntry, 0)
	errs := make([]error, 0)
	for _, z := range strings.Split(zone, ",") {
		zinfo, err := dnspodapi.GetDomainInfo(z)
		if err != nil {
			errs = append(errs, err)
		}
		if zinfo != nil {
			//fmt.Println(dnspodapi.FormatDomains([]dnspodapi.DomainEntry{*zinfo},
			//	format))
			//fmt.Printf("Domain `%s`, ID: %s, list records: \n", zinfo.Name,
			//	zinfo.ID)
		}
		rs, err := dnspodapi.ListRecord(z, "")
		if err != nil {
			errs = append(errs, err)
		}
		if rs != nil && len(rs) > 0 {
			fmt.Println(dnspodapi.FormatRecords(rs, format, z))
		}
	}
	if len(errs) > 0 {
		pe(errs...)
	}
}

func doCreateRecord() {
	rs, err := loadRecordFile(recordFile)
	if err != nil {
		pe(err)
		return
	}
	res, cls, err := addRecords(rs, zone, clear)
	//	fmt.Println("Created records:")
	fmt.Println(FormatOps(res, format, "Created records"))
	if cls != nil && len(cls) > 0 {
		//		fmt.Println("Cleared records:")
		fmt.Println(FormatOps(cls, format, "Cleared recoreds"))
	}
	if err != nil {
		pe(err)
	}
}

// OPRecordEntry defines the struct for adding record and result info
type OPRecordEntry struct {
	Action    string `json:"action"`
	Domain    string `json:"domain"`
	DomainID  string `json:"domain_id"`
	RecordID  string `json:"record_id"`
	SubDomain string `json:"sub_domain"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Err       error  `json:"err"`
	Message   string `json:"message"`
}

func loadRecordFile(f string) ([]*OPRecordEntry, error) {
	res := make([]*OPRecordEntry, 0)
	fp, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 3 {
			// maybe changed later
			continue
		}
		re := &OPRecordEntry{
			SubDomain: parts[0],
			Type:      parts[1],
			Value:     parts[2],
			Action:    RecordActionCreate,
		}
		res = append(res, re)
	}
	if err := scanner.Err(); err != nil {
		return res, err
	}
	return res, nil
}

// addRecords add all record into multi-zones
// returns all cleared records info if clear is true
func addRecords(rs []*OPRecordEntry, zs string, clear bool) (
	[]*OPRecordEntry, []*OPRecordEntry, error,
) {
	clret := make([]*OPRecordEntry, 0)
	rsret := make([]*OPRecordEntry, 0)
	errs := make([]error, 0)
	for i, z := range strings.Split(zs, ",") {
		zinfo, err := dnspodapi.GetDomainInfo(z)
		if err != nil {
			log.Println(err.Error())
			errs = append(errs, err)
			continue
		}
		zrs, err := dnspodapi.ListRecord(z, zinfo.ID)
		if err != nil {
			log.Println(err.Error())
			errs = append(errs, err)
			continue
		}
		var rss []*OPRecordEntry
		if i == 0 {
			rss = rs
		} else {
			rss = cprs(rs)
		}

		crs := addRecordsInZone(rss, zrs, zinfo, clear)

		clret = append(clret, crs...)
		rsret = append(rsret, rss...)
	}
	return rsret, clret, errgroup.NewMultiError(errs...)
}

// addRecordsInZone add all records into zone,
// return all cleared record info if clear is true
func addRecordsInZone(rs []*OPRecordEntry, zrs []dnspodapi.RecordEntry,
	zinfo *dnspodapi.DomainEntry, clear bool) []*OPRecordEntry {
	clearret := make([]*OPRecordEntry, 0)
	for _, r := range rs {
		rinfos := checkSubdomain(r, zrs)
		if rinfos != nil && len(rinfos) > 0 {
			// if exist
			if clear {
				// if clear, then first clear and add
				cls := clearSubdomain(zinfo.Name, zinfo.ID, rinfos)
				if cls != nil && len(cls) > 0 {
					clearret = append(clearret, cls...)
				}
				addRecordInZone(r, zinfo)
			} else {
				// if not clear, just do nothing, no adding and no clear
				r.Message = "record exists, don't clear"
				r.Domain = zinfo.Name
				r.DomainID = zinfo.ID
			}
		} else {
			// not exist, just create a record
			addRecordInZone(r, zinfo)
		}
	}
	return clearret
}

// checkSubdomain returns all records with same subdomain and type
func checkSubdomain(r *OPRecordEntry, zrs []dnspodapi.RecordEntry) []dnspodapi.RecordEntry {
	ret := make([]dnspodapi.RecordEntry, 0)
	for _, v := range zrs {
		if r.SubDomain == v.Name && r.Type == v.Type {
			ret = append(ret, v)
		}
	}
	return ret
}

func clearSubdomain(domain string, domainID string, zrs []dnspodapi.RecordEntry) []*OPRecordEntry {
	ret := make([]*OPRecordEntry, 0)
	for _, v := range zrs {
		or := &OPRecordEntry{
			Domain:    domain,
			DomainID:  domainID,
			SubDomain: v.Name,
			Type:      v.Type,
			Value:     v.Value,
			RecordID:  v.ID,
			Action:    RecordActionRemove,
		}
		ok, err := dnspodapi.RemoveRecord(domain, domainID, v.ID)
		if err != nil {
			or.Err = err
			or.Message = "remove record error"
		} else if ok {
			or.Message = "remove record succ"
		} else {
			or.Message = "remove record fail"
		}
		ret = append(ret, or)
	}
	return ret
}

func addRecordInZone(r *OPRecordEntry, zinfo *dnspodapi.DomainEntry) {
	data := genRecordParams(r, zinfo)
	re, err := dnspodapi.CreateRecord(zinfo.Name, zinfo.ID, data)
	if err != nil {
		r.Err = err
		r.Message = "add record in zone error"
	} else {
		r.RecordID = re.ID
		r.Message = "add record in zone succ"
	}
	r.Domain = zinfo.Name
	r.DomainID = zinfo.ID
}

func genRecordParams(r *OPRecordEntry, zinfo *dnspodapi.DomainEntry) dnspodapi.Params {
	res := make(dnspodapi.Params)
	res["record_line_id"] = DefaultRecordLineID
	if r.Type == "CNAME" && strings.HasSuffix(r.Value, ".") {
		// is local zone
		r.Value = fmt.Sprintf("%s%s", r.Value, zinfo.Name)
	}
	res["sub_domain"] = r.SubDomain
	res["record_type"] = r.Type
	res["value"] = r.Value
	res["domain"] = zinfo.Name
	res["domain_id"] = zinfo.ID
	return res
}

// FormatOps returns string of Ops slice
func FormatOps(rs []*OPRecordEntry, format string, cap string) string {
	res := ""
	switch format {
	case "json":
		bs, _ := json.Marshal(rs)
		res = string(bs)
	default:
		b := new(bytes.Buffer)
		table := tablewriter.NewWriter(b)
		dummy := &OPRecordEntry{}
		header := structs.StructKeys(dummy, true)
		table.SetCaption(true, cap)
		table.SetHeader(header)
		for _, r := range rs {
			table.Append(structs.StructValues(r, true))
		}
		table.Render()
		res = b.String()
	}
	return res
}

func pe(ers ...error) {
	if len(ers) > 0 {
		fmt.Printf("\nError:\n")
		for _, e := range ers {
			fmt.Printf("\t%s\n", e.Error())
		}
	}
}

func cprs(rs []*OPRecordEntry) []*OPRecordEntry {
	ret := make([]*OPRecordEntry, 0)
	for _, r := range rs {
		v := OPRecordEntry{}
		v = *r
		ret = append(ret, &v)
	}
	return ret
}
