// @Author: wangzn04@gmail.com
// @Date: 2018-08-31 15:15:18

package dnspodapi

import (
	"encoding/json"
	"log"
	"reflect"
	// "github.com/wangzn/dnspodapi"
)

// InfoModuleName defines the module name used in API
const InfoModuleName = "info"

var info Info

// Info defines the struct
type Info struct {
	A string
}

// Version returns version info
func (i Info) Version(bs []byte) interface{} {
	var d struct {
		Status RespCommon `json:"status"`
	}
	err := json.Unmarshal(bs, &d)
	if err != nil {
		return err
	}
	return d
}

func init() {
	info = Info{}
	Register(InfoModuleName, infoReflectFunc)
}

// infoReflectFunc defines the reflect func
func infoReflectFunc(action string, data Params) ActionResult {
	log.Println("in info reflect func, aciton:", action)
	return callReflectFunc(reflect.ValueOf(info), InfoModuleName, action, nil,
		nil)
}
