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
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/wangzn/dnspodapi"
	"github.com/wangzn/goutils/mymap"
	"github.com/wangzn/goutils/structs"
	"github.com/wangzn/goutils/sys"
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
	records        string
	typ            string
	value          string
	recordAct      string
	clear          bool
	recordFile     string
	zone           string
	exportFileMode string
	forceDomain    bool
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

	recordCmd.PersistentFlags().StringVarP(&records, "record", "r", "",
		"records to operate")

	recordCmd.PersistentFlags().StringVarP(&typ, "type", "t", "",
		"record type, [ CNAME | A | MX | AAAA ]")

	recordCmd.PersistentFlags().StringVarP(&value, "value", "v", "",
		"record value")

	recordCmd.PersistentFlags().StringVarP(&recordAct, "action", "a", "list",
		"record action: [ create | remove | info | list | import | export | enable | disable ]")

	recordCmd.PersistentFlags().StringVar(&format, "format", "table",
		"output format: [ json | table ]")

	recordCmd.Flags().BoolVarP(&clear, "clear", "c", false,
		"clear existed record when importing records")

	recordCmd.Flags().BoolVar(&forceDomain, "force-domain", false,
		"force create new domain if not exist")

	recordCmd.Flags().StringVarP(&zone, "domain", "d", "",
		"domains for records, use ',' for multiple domains, e.g. 'abc.com,def.com'")

	recordCmd.Flags().StringVar(&exportFileMode, "export-file-mode", "",
		"export file mode: [ append | overwrite ], exit if file exists in default")

	recordCmd.Flags().StringVarP(&recordFile, "record_file", "f", "",
		"record file, each line contains 'record type value', e.g. 'www CNAME proxy'")
}

func runRecordCmd(cmd *cobra.Command, args []string) {
	r := RecordActionRunner{
		Record:   records,
		Action:   recordAct,
		APIID:    apiID,
		APIToken: apiToken,
		Params:   fillRecordParams(),
	}
	r.Run()
}

func fillRecordParams() map[string]string {
	m := make(map[string]string)
	m["domain"] = zone
	m["value"] = value
	m["type"] = typ
	m["clear"] = "off"
	if clear {
		m["clear"] = "on"
	}
	m["force_domain"] = "off"
	if forceDomain {
		m["force_domain"] = "on"
	}
	m["record_file"] = recordFile
	m["export_file_mode"] = exportFileMode
	return m
}

// RecordActionRunner defines the runner to run record action
type RecordActionRunner struct {
	Record   string
	Action   string
	Params   map[string]string
	APIID    int
	APIToken string
}

// Run starts to run action
func (r *RecordActionRunner) Run() {
	r.run()
}

// Name returns runner name
func (r *RecordActionRunner) Name() string {
	return fmt.Sprintf("%s record `%s` of domain `%s`",
		strings.Title(r.Action), r.Record,
		mymap.StringMustString(r.Params, "domain"))
}

// Detail returns detail information
func (r *RecordActionRunner) Detail() string {
	res := r.Name()
	switch r.Action {
	case "import":
		res = fmt.Sprintf("%s record into domain `%s` from file `%s`, "+
			"with clear flag `%s` and force-domain flag `%s`",
			strings.Title(r.Action),
			mymap.StringMustString(r.Params, "domain"),
			mymap.StringMustString(r.Params, "record_file"),
			mymap.StringMustString(r.Params, "clear"),
			mymap.StringMustString(r.Params, "force_domain"))
	case "export":
		res = fmt.Sprintf("%s record from domain `%s` into file `%s`, "+
			"with export-file-mode `%s`",
			strings.Title(r.Action),
			mymap.StringMustString(r.Params, "domain"),
			mymap.StringMustString(r.Params, "record_file"),
			mymap.StringMustString(r.Params, "export_file_mode"))
	}
	return res

}

func (r *RecordActionRunner) run() {
	dms := mymap.StringMustString(r.Params, "domain")
	if dms == "" {
		fmt.Println("domain is empty")
		os.Exit(1)
	}
	switch r.Action {
	case "create":
		r.doCreateRecord()
	case "remove":
		r.doRemoveRecord()
	case "info":
		r.doInfoRecord()
	case "export":
		r.doExportRecord()
	case "import":
		r.doImportRecord()
	case "enable", "disable":
		r.doStatusRecord()
	default:
		r.doListRecord()
	}
}

func (r *RecordActionRunner) toOPRecordEntry() []*OPRecordEntry {
	res := make([]*OPRecordEntry, 0)
	typ := mymap.StringMustString(r.Params, "type")
	val := mymap.StringMustString(r.Params, "value")
	for _, record := range strings.Split(r.Record, ",") {
		re := &OPRecordEntry{
			SubDomain: record,
			Type:      typ,
			Value:     val,
			Action:    RecordActionCreate,
		}
		res = append(res, re)
	}
	return res
}

func (r *RecordActionRunner) doCreateRecord() {
	data := r.Params
	cl := mymap.StringMustString(data, "clear")
	fd := mymap.StringMustString(data, "force_domain")
	dms := mymap.StringMustString(data, "domain")
	clb := false
	fdb := false
	if cl == "on" {
		clb = true
	}
	if fd == "on" {
		fdb = true
	}

	rs := r.toOPRecordEntry()

	res, cls, errs := addRecords(rs, dms, clb, fdb)
	fmt.Println(FormatOps(res, format, "Created records"))
	if cls != nil && len(cls) > 0 {
		fmt.Println(FormatOps(cls, format, "Cleared recoreds"))
	}
	pe(errs...)
}

func (r *RecordActionRunner) doRemoveRecord() {
	errs := make([]error, 0)
	dms := mymap.StringMustString(r.Params, "domain")
	res := make([]*OPRecordEntry, 0)
	for _, domain := range strings.Split(dms, ",") {
		// get all records with r.Record in domain
		zinfo, crs, ers := filterRecords(domain, r.Record)
		if ers != nil && len(ers) > 0 {
			errs = append(errs, ers...)
		}
		cls := clearSubdomain(zinfo.Name, zinfo.ID, crs)
		if cls != nil && len(cls) > 0 {
			res = append(res, cls...)
		}
	}
	fmt.Println(FormatOps(res, format, "Removed recoreds"))
	pe(errs...)
}

func (r *RecordActionRunner) doInfoRecord() {
	errs := make([]error, 0)
	dms := mymap.StringMustString(r.Params, "domain")
	for _, domain := range strings.Split(dms, ",") {
		// get all records with r.Record in domain
		_, crs, ers := filterRecords(domain, r.Record)
		if ers != nil && len(ers) > 0 {
			errs = append(errs, ers...)
		}
		fmt.Println(dnspodapi.FormatRecords(crs, format, domain))
	}
	pe(errs...)
}

func (r *RecordActionRunner) doExportRecord() {
	errs := make([]error, 0)
	dms := mymap.StringMustString(r.Params, "domain")
	fn := mymap.StringMustString(r.Params, "record_file")
	fmode := mymap.StringMustString(r.Params, "export_file_mode")
	fp, err := getFilePointer(fn, fmode)
	defer fp.Close()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, domain := range strings.Split(dms, ",") {
		_, rs, ers := filterRecords(domain, r.Record)
		if ers != nil && len(ers) > 0 {
			errs = append(errs, ers...)
		}
		_, err = fp.Write([]byte(fmt.Sprintf("# export domain `%s` with "+
			"%d records at %s \n",
			domain,
			len(rs),
			time.Now().String())))
		if err != nil {
			errs = append(errs, err)
		}
		for _, r := range rs {
			_, err = fp.Write([]byte(r.ExportLine() + "\n"))
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
}

func (r *RecordActionRunner) doListRecord() {
	//res := make([]dnspodapi.RecordEntry, 0)
	data := r.Params
	dms := mymap.StringMustString(data, "domain")
	errs := make([]error, 0)
	for _, z := range strings.Split(dms, ",") {
		zinfo, err := dnspodapi.GetDomainInfo(z)
		if err != nil {
			errs = append(errs, err)
		}
		if zinfo != nil {
			if zinfo.ID == "" {
				errs = append(errs, fmt.Errorf("Fail to get domain info: `%s`",
					z))
				continue
			}
		} else {
			errs = append(errs, fmt.Errorf("Fail to get domain info: `%s`", z))
			continue
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

func (r *RecordActionRunner) doImportRecord() {
	data := r.Params
	fn := mymap.StringMustString(data, "record_file")
	cl := mymap.StringMustString(data, "clear")
	fd := mymap.StringMustString(data, "force_domain")
	dms := mymap.StringMustString(data, "domain")
	clb := false
	fdb := false
	if cl == "on" {
		clb = true
	}
	if fd == "on" {
		fdb = true
	}

	rs, err := loadRecordFile(fn)
	if err != nil {
		pe(err)
		return
	}

	res, cls, errs := addRecords(rs, dms, clb, fdb)
	//	fmt.Println("Created records:")
	fmt.Println(FormatOps(res, format, "Imported records"))
	if cls != nil && len(cls) > 0 {
		//		fmt.Println("Cleared records:")
		fmt.Println(FormatOps(cls, format, "Cleared recoreds"))
	}
	pe(errs...)
}

func (r *RecordActionRunner) doStatusRecord() {
	errs := make([]error, 0)
	res := make([][]string, 0)
	dms := mymap.StringMustString(r.Params, "domain")
	header := []string{"domain", "record", "status", "msg"}
	for _, d := range strings.Split(dms, ",") {
		zinfo, rs, ers := filterRecords(d, r.Record)
		if ers != nil && len(ers) > 0 {
			errs = append(errs, ers...)
		}
		for _, record := range rs {
			err := dnspodapi.SetRecordStatus(zinfo.Name, zinfo.ID, record.ID,
				r.Action)
			msg := ""
			ok := true
			if err != nil {
				errs = append(errs, err)
				msg = err.Error()
				ok = false
			}
			res = append(res, []string{d, record.Name, fmt.Sprintf("%v", ok),
				msg})
		}
	}
	fmt.Println(mymap.FormatSlices(header, res, format))
	pe(errs...)
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
	RealValue string `json:"real_value"`
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
		line := strings.Trim(scanner.Text(), " ")
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		parts := strings.Fields(line)
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
func addRecords(rs []*OPRecordEntry, zs string, clear, autoDomain bool) (
	[]*OPRecordEntry, []*OPRecordEntry, []error,
) {
	clret := make([]*OPRecordEntry, 0)
	rsret := make([]*OPRecordEntry, 0)
	errs := make([]error, 0)
	for i, z := range strings.Split(zs, ",") {
		zinfo, err := dnspodapi.GetDomainInfo(z)
		if err != nil {
			if autoDomain {
				zinfo, err = ensureDomain(z)
				if err != nil {
					// fail to create domain, continue anyway
					errs = append(errs, err)
					continue
				}
			} else {
				// no force auto create domain, just continue
				errs = append(errs, err)
				continue
			}
		}
		if zinfo == nil || zinfo.ID == "" {
			errs = append(errs, fmt.Errorf("Fail to get info of domain `%s`",
				z))
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
	return rsret, clret, errs
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
func checkSubdomain(r *OPRecordEntry, zrs []dnspodapi.RecordEntry,
) []dnspodapi.RecordEntry {
	ret := make([]dnspodapi.RecordEntry, 0)
	for _, v := range zrs {
		if r.SubDomain == v.Name && r.Type == v.Type {
			ret = append(ret, v)
		}
	}
	return ret
}

func clearSubdomain(domain string, domainID string, zrs []dnspodapi.RecordEntry,
) []*OPRecordEntry {
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

func genRecordParams(r *OPRecordEntry, zinfo *dnspodapi.DomainEntry,
) dnspodapi.Params {
	res := dnspodapi.P()
	res.Add("record_line_id", DefaultRecordLineID)
	if r.Type == "CNAME" && strings.HasSuffix(r.Value, ".") {
		// is local zone
		r.RealValue = fmt.Sprintf("%s%s", r.Value, zinfo.Name)
	} else {
		r.RealValue = r.Value
	}
	res.Add("sub_domain", r.SubDomain)
	res.Add("record_type", r.Type)
	res.Add("value", r.RealValue)
	res.Add("domain", zinfo.Name)
	res.Add("domain_id", zinfo.ID)
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
	if ers != nil && len(ers) > 0 {
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

// ensureDomain add a domain if not exist
func ensureDomain(domain string) (*dnspodapi.DomainEntry, error) {
	return dnspodapi.CreateDomain(domain)
}

func getFilePointer(fn, fm string) (*os.File, error) {
	switch fm {
	case "append":
		return os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	case "overwrite":
		err := os.Remove(fn)
		if err != nil {
			return nil, err
		}
		return os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
	default:
		if sys.FNExist(fn) {
			return os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
		}
		return nil, fmt.Errorf("file exist")
	}
}

func filterRecords(domain, records string) (*dnspodapi.DomainEntry,
	[]dnspodapi.RecordEntry, []error) {
	errs := make([]error, 0)
	zinfo, err := dnspodapi.GetDomainInfo(domain)
	if err != nil {
		errs = append(errs, err)
		return nil, nil, errs
	}
	rs, err := dnspodapi.ListRecord(zinfo.Name, zinfo.ID)
	if err != nil {
		errs = append(errs, err)
		return nil, nil, errs
	}
	crs := make([]dnspodapi.RecordEntry, 0)
	if records != "" {
		for _, record := range strings.Split(records, ",") {
			for _, r := range rs {
				if r.Name == record {
					crs = append(crs, r)
					// can not break here, because record could have different types
				}
			}
		}
	} else {
		crs = rs
	}
	return zinfo, crs, errs
}
