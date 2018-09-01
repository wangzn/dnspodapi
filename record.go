// Copyright 2018 TED@Sogou, Inc. All rights reserved.

// @Author: wangzhongning@sogou-inc.com
// @Date: 2018-08-31 14:55:46

package dnspodapi

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
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

// RecordCreateOrModifyResult defines the API result of `create` or `Modify`
type RecordCreateOrModifyResult struct {
	Status RespCommon  `json:"status"`
	Record RecordEntry `json:"record"`
}

// RecordRemoveResult defines the API result of `remove`
type RecordRemoveResult struct {
	Status RespCommon `json:"status"`
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
	for _, v := range []string{"domain", "record_id", "domain_id", "sub_domain",
		"record_type", "record_line", "record_line_id", "value", "mx", "ttl",
		"status", "weight"} {
		if dn, ok := data[v]; ok {
			pd.Add(v, dn.(string))
		}
	}
	return callReflectFunc(reflect.ValueOf(record), RecordModuleName, action,
		nil, pd)
}

// List returns record list
func (r Record) List(bs []byte) *RecordListResult {
	data := new(RecordListResult)
	err := json.Unmarshal(bs, data)
	if err != nil {
		return nil
	}
	return data
}

// Info returns record info
func (r Record) Info(bs []byte) *RecordInfoResult {
	data := new(RecordInfoResult)
	err := json.Unmarshal(bs, data)
	if err != nil {
		return nil
	}
	return data
}

// Create returns created record
func (r Record) Create(bs []byte) *RecordCreateOrModifyResult {
	data := new(RecordCreateOrModifyResult)
	err := json.Unmarshal(bs, data)
	if err != nil {
		return nil
	}
	return data
}

// Modify returns modified record
func (r Record) Modify(bs []byte) *RecordCreateOrModifyResult {
	return r.Create(bs)
}

// Remove returns removed record
func (r Record) Remove(bs []byte) *RecordRemoveResult {
	data := new(RecordRemoveResult)
	err := json.Unmarshal(bs, data)
	if err != nil {
		return nil
	}
	return data
}

// CreateRecord creates a new record
func CreateRecord(domain string, domainID int, data Params) (
	*RecordEntry, error) {
	res := recordReflectFunc("create", data)
	if res.Err != nil {
		return nil, res.Err
	}
	if ret, ok := res.Data.(RecordCreateOrModifyResult); ok {
		return GetRecordInfo(ret.Record.ID)
	}
	return nil, Err(ErrInvalidTypeAssertion, "RecordCreateResult")
}

// ModifyRecord modifies a record
func ModifyRecord(domain string, domainID int, recordID string, data Params) (
	*RecordEntry, error) {
	data["record_id"] = recordID
	if domain != "" {
		data["domain"] = domain
	}
	if domainID > 0 {
		data["domain_id"] = strconv.Itoa(domainID)
	}
	res := recordReflectFunc("modify", data)
	if res.Err != nil {
		return nil, res.Err
	}
	if ret, ok := res.Data.(RecordCreateOrModifyResult); ok {
		return GetRecordInfo(ret.Record.ID)
	}
	return nil, Err(ErrInvalidTypeAssertion, "RecordModifyResult")
}

// RemoveRecord removes a record
func RemoveRecord(domain string, domainID int, recordID string) (bool, error) {
	data := make(map[string]interface{})
	data["record_id"] = recordID
	if domain != "" {
		data["domain"] = domain
	}
	if domainID > 0 {
		data["domain_id"] = domainID
	}
	res := recordReflectFunc("remove", data)
	if res.Err != nil {
		return false, res.Err
	}
	if ret, ok := res.Data.(RecordRemoveResult); ok {
		if ret.Status.Code == "1" {
			return true, nil
		}
		return false, Err(ErrInvalidStatus, ret.Status.Code, ret.Status.Message)
	}
	return false, Err(ErrInvalidTypeAssertion, "RecordRemoveResult")
}

// GetRecordInfo returns record entry
func GetRecordInfo(recordID string) (*RecordEntry, error) {
	data := make(map[string]interface{})
	data["record_id"] = recordID
	res := recordReflectFunc("info", data)
	if res.Err != nil {
		return nil, res.Err
	}
	if ret, ok := res.Data.(RecordInfoResult); ok {
		return &ret.Record, nil
	}
	return nil, Err(ErrInvalidTypeAssertion, "RecordInfoResult")
}
