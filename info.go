// Copyright 2018 TED@Sogou, Inc. All rights reserved.

// @Author: wangzhongning@sogou-inc.com
// @Date: 2018-08-31 15:15:18

package dnspodapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	// "github.com/wangzn/dnspodapi"
)

// ModuleName defines the module name used in API
const ModuleName = "Info"

var info Info

// Info defines the struct
type Info struct{}

func init() {
	Register(ModuleName, infoReflectFunc)
}

// infoReflectFunc defines the reflect func
func infoReflectFunc(action string, data Params) ActionResult {
	code, bs, err := HTTPResp(ModuleName, action, nil, nil)
	if err != nil {
		return ActionResult{
			Code: code,
			Err:  err,
		}
	}
	in := []reflect.Value{reflect.ValueOf(bs)}
	out := reflect.ValueOf(info).Call(in)
	if len(out) == 0 {
		return ActionResult{
			Code: 1001,
			Err:  fmt.Errorf("invalid call func result: not enough result"),
		}
	}
	return ActionResult{
		Code: 1000,
		Err:  nil,
		Data: out[0].Interface(),
	}
}

// Version returns version info
func (i Info) version(bs []byte) interface{} {
	//{
	//	"status": {
	//		"code": "1",
	//		"message": "4.6",
	//		"created_at": "2012-09-10 11:20:39"
	//	}
	//}
	var d struct {
		Status RespCommon `json:"status"`
	}
	err := json.Unmarshal(bs, &d)
	if err != nil {
		return err
	}
	return d
}
