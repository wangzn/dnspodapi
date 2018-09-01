// Copyright 2018 TED@Sogou, Inc. All rights reserved.

// @Author: wangzhongning@sogou-inc.com
// @Date: 2018-08-31 14:55:30

package dnspodapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/wangzn/goutils/structs"
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
	pd := url.Values{}
	if dn, ok := data["domain"]; ok {
		pd.Add("domain", dn.(string))
	}
	return callReflectFunc(reflect.ValueOf(domain), DomainModuleName, action,
		nil, pd)
}

// List returns domain list
func (d Domain) List(bs []byte) *DomainListResult {
	data := new(DomainListResult)
	err := json.Unmarshal(bs, data)
	if err != nil {
		return nil
	}
	return data
}

// Info returns domain info
func (d Domain) Info(bs []byte) *DomainInfoResult {
	data := new(DomainInfoResult)
	err := json.Unmarshal(bs, data)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}

// GetDomainInfo returns domain entry
func GetDomainInfo(name string) (*DomainEntry, error) {
	data := make(map[string]interface{})
	data["domain"] = name
	res := domainReflectFunc("info", data)
	if res.Err != nil {
		return nil, res.Err
	}
	if ret, ok := res.Data.(*DomainInfoResult); ok {
		if ret != nil {
			return &ret.Domain, nil
		}
		return nil, Err(ErrInvalidTypeAssertion, "DomainInfoResult")
	}
	return nil, Err(ErrInvalidTypeAssertion, "DomainInfoResult")
}

// GetDomainList returns domain entry list
func GetDomainList() ([]DomainEntryIDInt, error) {
	res := domainReflectFunc("list", nil)
	if res.Err != nil {
		return nil, res.Err
	}
	if ret, ok := res.Data.(DomainListResult); ok {
		return ret.Domains, nil
	}
	return nil, Err(ErrInvalidTypeAssertion, "DomainListResult")
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
