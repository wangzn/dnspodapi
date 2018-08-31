// Copyright 2018 TED@Sogou, Inc. All rights reserved.

// @Author: wangzhongning@sogou-inc.com
// @Date: 2018-08-31 14:55:46

package dnspodapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
)

// RecordEntry defines the API result struct of record line
type RecordEntry struct {
	ID            string `json:"id"`
	TTL           string `json:"ttl"`
	Value         string `json:"value"`
	Enabled       string `json:"enabled"`
	Status        string `json:"status"`
	UpdatedOn     string `json:"updated_on"`
	Name          string `json:"name"`
	Line          string `json:"line"`
	LineID        string `json:"line_id"`
	Type          string `json:"type"`
	Weight        string `json:"weight"`
	MonitorStatus string `json:"monitor_status"`
	Remark        string `json:"remark"`
	UseAqb        string `json:"use_aqb"`
	MX            string `json:"mx"`
	Hold          string `json:"hold"`
}

// RecordDomainEntry defines the `domain field in `list`
type RecordDomainEntry struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	TTL      int      `json:"ttl"`
	MinTTL   int      `json:"min_ttl"`
	Status   string   `json:"status"`
	DnspodNS []string `json:"dnspod_ns"`
}

// RecordListInfo defines the `info` field in `list`
type RecordListInfo struct {
	SubDomains  string `json:"sub_domains"`
	RecordTotal string `json:"record_total"`
	RecordsNum  string `json:"records_num"`
}

// RecordListResult defiens the API result of `list`
type RecordListResult struct {
	Status  RespCommon        `json:"status"`
	Domain  RecordDomainEntry `json:"domain"`
	Info    RecordListInfo    `json:"info"`
	Records []RecordEntry     `json:"records"`
}

// RecordInfoResult defines the API result of `info`
type RecordInfoResult struct {
	Status RespCommon        `json:"status"`
	Domain RecordDomainEntry `json:"domain"`
	Record RecordEntry       `json:"record"`
}

const (
	// RecordModuleName defines the const value of record module
	RecordModuleName = "record"
)

// Record defines the dummy struct
type Record struct{}

var record Record

func init() {
	record = Record{}
	Register(RecordModuleName, recordReflectFunc)
}

// recordReflectFunc defines the reflect func
// record
//			list
//			get $record_id
func recordReflectFunc(action string, data Params) ActionResult {
	pd := url.Values{}
	for _, v := range []string{"domain", "record_id"} {
		if dn, ok := data[v]; ok {
			pd.Add(v, dn.(string))
		}
	}
	return callReflectFunc(reflect.ValueOf(record), RecordModuleName, action,
		nil, pd)
}

// List returns record list
func (r Record) List(bs []byte) *RecordListResult {
	fmt.Println(string(bs))
	data := new(RecordListResult)
	err := json.Unmarshal(bs, data)
	if err != nil {
		return nil
	}
	return data
}
