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
type Info struct{}

// Version returns version info
func (i Info) Version(bs []byte) interface{} {
	d := new(InfoVersionResult)
	err := json.Unmarshal(bs, d)
	if err != nil {
		return err
	}
	return d
}

// InfoVersionResult defines the API result of info.version
type InfoVersionResult struct {
	Status RespCommon `json:"status"`
}

func init() {
	info = Info{}
	Register(InfoModuleName, infoReflectFunc)
}

// infoReflectFunc defines the reflect func
func infoReflectFunc(action string, data Params) ActionResult {
	log.Println("in info reflect func, aciton:", action)
	return callReflectFunc(reflect.ValueOf(info), InfoModuleName, action, nil,
		data)
}

// GetVersion returns API version
func GetVersion() (string, error) {
	res := infoReflectFunc("version", P())
	if res.Err != nil {
		return "", res.Err
	}
	if ret, ok := res.Data.(*InfoVersionResult); ok {
		if ret != nil {
			if ret.Status.Code == "1" {
				return ret.Status.Message, nil
			}
			return "", Err(ErrInvalidStatus, "Info.Version", ret.Status.Code,
				ret.Status.Message)
		}
		return "", Err(ErrInvalidTypeAssertion, "InfoVersionResult")
	}
	return "", Err(ErrInvalidTypeAssertion, "InfoVersionResult")
}
