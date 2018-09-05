// @Author: wangzn04@gmail.com
// @Date: 2018-08-31 14:55:30

package dnspodapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/wangzn/goutils/structs"
)

const (
	// DomainStatusEnable defines the string of `enable` of a domain
	DomainStatusEnable = "enable"
	// DomainStatusDisable defines the string of `disable` of a domain
	DomainStatusDisable = "disable"
	// DomainStatusUnkown defines the string of `unkown` of a domain
	DomainStatusUnkown = ""
)

// DomainEntry defines the API result struct of domain line
type DomainEntry struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	Grade            string `json:"grade"`
	GroupID          string `json:"group_id"`
	SearchenginePush string `json:"searchengine_push"`
	IsMark           string `json:"is_marK"`
	TTL              string `json:"ttl"`
	CnameSpeedup     string `json:"cname_speedup"`
	Remark           string `json:"remark"`
	CreatedOn        string `json:"created_on"`
	UpdatedOn        string `json:"updated_on"`
	Punycode         string `json:"punycode"`
	ExtStatus        string `json:"ext_status"`
	SrcFlag          string `json:"src_flag"`
	Name             string `json:"name"`
	GradeTitle       string `json:"grade_title"`
	IsVIP            string `json:"is_ip"`
	Owner            string `json:"owner"`
	Records          string `json:"records"`
}

// DomainEntryIDInt defines same for DomainEntry only id is int
type DomainEntryIDInt struct {
	ID               int    `json:"id"`
	Status           string `json:"status"`
	Grade            string `json:"grade"`
	GroupID          string `json:"group_id"`
	SearchenginePush string `json:"searchengine_push"`
	IsMark           string `json:"is_marK"`
	TTL              string `json:"ttl"`
	CnameSpeedup     string `json:"cname_speedup"`
	Remark           string `json:"remark"`
	CreatedOn        string `json:"created_on"`
	UpdatedOn        string `json:"updated_on"`
	Punycode         string `json:"punycode"`
	ExtStatus        string `json:"ext_status"`
	SrcFlag          string `json:"src_flag"`
	Name             string `json:"name"`
	GradeTitle       string `json:"grade_title"`
	IsVIP            string `json:"is_ip"`
	Owner            string `json:"owner"`
	Records          string `json:"records"`
}

// DomainListInfo defines the API struct of `info` field
type DomainListInfo struct {
	DomainTotal   int    `json:"domain_total"`
	AllTotal      int    `json:"all_total"`
	MineTotal     int    `json:"mine_total"`
	ShareTotal    string `json:"share_total"`
	VIPTotal      int    `json:"vip_total"`
	IsmarkTotal   int    `json:"ismark_total"`
	PauseTotal    int    `json:"pause_total"`
	ErrorTotal    int    `json:"error_total"`
	LockTotal     int    `json:"lock_total"`
	SpamTotal     int    `json:"spam_total"`
	VIPExpire     int    `json:"vip_expire"`
	ShareOutTotal int    `json:"share_out_total"`
}

// DomainListResult defines the API result of `list`
type DomainListResult struct {
	Status  RespCommon         `json:"status"`
	Info    DomainListInfo     `json:"info"`
	Domains []DomainEntryIDInt `json:"domains"`
}

// DomainInfoResult defines the API result of `info``
type DomainInfoResult struct {
	Status RespCommon  `json:"status"`
	Domain DomainEntry `json:"domain"`
}

// DomainCreateInfo defines the struct of field `domain` in create result
type DomainCreateInfo struct {
	ID        string `json:"id"`
	Punnycode string `json:"punnycode"`
	Domain    string `json:"domain"`
}

// DomainCreateResult defines the API result of `create`
type DomainCreateResult struct {
	Status RespCommon       `json:"status"`
	Domain DomainCreateInfo `json:"domain"`
}

// DomainRemoveResult defines the API result of `remove`
type DomainRemoveResult struct {
	Status RespCommon `json:"status"`
}

// DomainStatusResult defines the API result of `status`
type DomainStatusResult struct {
	Status RespCommon `json:"status"`
}

const (
	// DomainModuleName defines the const value of domain module
	DomainModuleName = "domain"
)

// Domain defines the dummy struct
type Domain struct{}

var domain Domain

func init() {
	domain = Domain{}
	Register(DomainModuleName, domainReflectFunc)
}

// domainReflectFunc defines the reflect func
// domain
//			list
//			get $domain
func domainReflectFunc(action string, data Params) ActionResult {
	return callReflectFunc(reflect.ValueOf(domain), DomainModuleName, action,
		nil, data)
}

// List returns domain list
func (d Domain) List(bs []byte) *DomainListResult {
	data := new(DomainListResult)
	err := json.Unmarshal(bs, data)
	if err != nil {
		log.Println("fail to unmarshal domainlistresult:", err.Error())
		return nil
	}
	return data
}

// Info returns domain info
func (d Domain) Info(bs []byte) *DomainInfoResult {
	data := new(DomainInfoResult)
	err := json.Unmarshal(bs, data)
	if err != nil {
		log.Println("fail to unmarshal domaininforesult:", err.Error())
		return nil
	}
	return data
}

// Create returns domain created result
func (d Domain) Create(bs []byte) *DomainCreateResult {
	data := new(DomainCreateResult)
	err := json.Unmarshal(bs, data)
	if err != nil {
		log.Println("fail to unmarshal domaincreateresult:", err.Error())
		return nil
	}
	return data
}

// Remove removes a domain
func (d Domain) Remove(bs []byte) *DomainRemoveResult {
	data := new(DomainRemoveResult)
	err := json.Unmarshal(bs, data)
	if err != nil {
		log.Println("fail to unmarshal domainremoveresult:", err.Error())
		return nil
	}
	return data
}

// Status sets status of a domain
func (d Domain) Status(bs []byte) *DomainStatusResult {
	data := new(DomainStatusResult)
	err := json.Unmarshal(bs, data)
	if err != nil {
		log.Println("fail to unmarshal domainstatusresult:", err.Error())
		return nil
	}
	return data
}

// CreateDomain creates a domain entry
// optional param: group_id, is_mark
func CreateDomain(domain string) (*DomainEntry, error) {
	return CreateDomainDetail(domain, "", "")
}

// CreateDomainDetail with optional params
func CreateDomainDetail(domain, groupID, isMark string) (*DomainEntry, error) {
	data := P()
	data.Add("domain", domain)
	if groupID != "" {
		data.Add("group_id", groupID)
	}
	if isMark != "" {
		data.Add("is_mark", isMark)
	}
	res := domainReflectFunc("create", data)
	if res.Err != nil {
		return nil, res.Err
	}
	if ret, ok := res.Data.(*DomainCreateResult); ok {
		if ret != nil {
			if ret.Status.Code != "1" {
				return nil, fmt.Errorf("invalid resp status, msg: %s",
					ret.Status.Message)
			}
			id := ret.Domain.ID
			if id != "" {
				return GetDomainInfo(domain)
			}
		}
		return nil, Err(ErrInvalidTypeAssertion, "DomainCreateResult")
	}
	return nil, Err(ErrInvalidTypeAssertion, "DomainCreateResult")
}

// GetDomainInfo returns domain entry
func GetDomainInfo(name string) (*DomainEntry, error) {
	data := P()
	data.Add("domain", name)
	res := domainReflectFunc("info", Params(data))
	if res.Err != nil {
		return nil, res.Err
	}
	if ret, ok := res.Data.(*DomainInfoResult); ok {
		if ret != nil {
			if ret.Status.Code != "1" {
				return nil, fmt.Errorf("invalid resp status, msg: %s",
					ret.Status.Message)
			}
			return &ret.Domain, nil
		}
		return nil, Err(ErrInvalidTypeAssertion, "DomainInfoResult")
	}
	return nil, Err(ErrInvalidTypeAssertion, "DomainInfoResult")
}

// GetDomainList returns domain entry list
func GetDomainList() ([]DomainEntryIDInt, error) {
	res := domainReflectFunc("list", Params{})
	if res.Err != nil {
		return nil, res.Err
	}
	if ret, ok := res.Data.(*DomainListResult); ok {
		if ret != nil {
			if ret.Status.Code != "1" {
				return nil, fmt.Errorf("invalid resp status, msg: %s",
					ret.Status.Message)
			}
			return ret.Domains, nil
		}
	}
	return nil, Err(ErrInvalidTypeAssertion, "DomainListResult")
}

// RemoveDomain removes domain
func RemoveDomain(name string) (bool, error) {
	data := P()
	data.Add("domain", name)
	res := domainReflectFunc("remove", data)
	if res.Err != nil {
		return false, res.Err
	}
	if ret, ok := res.Data.(*DomainRemoveResult); ok {
		if ret != nil {
			if ret.Status.Code == "1" {
				return true, nil
			}
			return false, fmt.Errorf("remove fail, code: %s, msg: %s",
				ret.Status.Code, ret.Status.Message)
		}
	}
	return false, Err(ErrInvalidTypeAssertion, "DomainRemoveResult")
}

// SetDomainStatus set domain status
// enable: enable, 1, online, on
// disable: disable, 0, offline, off
func SetDomainStatus(name string, status string) error {
	st := verifyStatus(status)
	if st == DomainStatusUnkown {
		return fmt.Errorf("invalid target status, accept `enable` or `disable`")
	}
	data := P()
	data.Add("domain", name)
	data.Add("status", st)
	res := domainReflectFunc("status", data)
	if res.Err != nil {
		return res.Err
	}
	if ret, ok := res.Data.(*DomainStatusResult); ok {
		if ret != nil {
			if ret.Status.Code == "1" {
				return nil
			}
			return fmt.Errorf("status fail, code: %s, msg: %s",
				ret.Status.Code, ret.Status.Message)
		}
	}
	return Err(ErrInvalidTypeAssertion, "DomainStatusResult")
}

// FormatDomains returns output string
func FormatDomains(rs []DomainEntry, format string) string {
	res := ""
	switch format {
	case "json":
		bs, _ := json.Marshal(rs)
		res = string(bs)
	default:
		b := new(bytes.Buffer)
		table := tablewriter.NewWriter(b)
		dummy := DomainEntry{}
		header := structs.StructKeys(dummy, true)
		table.SetHeader(header)
		for _, r := range rs {
			table.Append(structs.StructValues(r, true))
		}
		if len(header) > 2 {
			total := make([]string, len(header))
			total[len(total)-1] = strconv.Itoa(len(rs))
			total[len(total)-2] = "TOTAL"
			table.SetFooter(total)
		}
		table.Render()
		res = b.String()
	}
	return res
}

// FormatDomainIDInts returns output string
func FormatDomainIDInts(rs []DomainEntryIDInt, format string) string {
	res := ""
	switch format {
	case "json":
		bs, _ := json.Marshal(rs)
		res = string(bs)
	default:
		b := new(bytes.Buffer)
		table := tablewriter.NewWriter(b)
		dummy := DomainEntryIDInt{}
		header := structs.StructKeys(dummy, true)
		table.SetHeader(header)
		for _, r := range rs {
			table.Append(structs.StructValues(r, true))
		}
		if len(header) > 2 {
			total := make([]string, len(header))
			total[len(total)-1] = strconv.Itoa(len(rs))
			total[len(total)-2] = "TOTAL"
			table.SetFooter(total)
		}
		table.Render()
		res = b.String()
	}
	return res
}

// enable: enable, 1, online, on
// disable: disable, 0, offline, off
func verifyStatus(st string) string {
	en := map[string]bool{
		"enable": true,
		"online": true,
		"on":     true,
		"1":      true,

		"disable": false,
		"offline": false,
		"off":     false,
		"0":       false,
	}
	if v, ok := en[st]; ok {
		if v {
			return DomainStatusEnable
		}
		return DomainStatusDisable
	}
	return DomainStatusUnkown
}
