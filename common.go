// @Author: wangzn04@gmail.com
// @Date: 2018-08-31 14:55:56

package dnspodapi

import (
	"log"
	"net/url"
	"reflect"
	"strings"
)

const (
	// StatusEnable defines the string of `enable` of a domain
	StatusEnable = "enable"
	// StatusDisable defines the string of `disable` of a domain
	StatusDisable = "disable"
	// StatusUnkown defines the string of `unkown` of a domain
	StatusUnkown = ""
)

// RespCommon defines common struct in http endpoint
type RespCommon struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

func callReflectFunc(v reflect.Value, module, action string, vs url.Values, data Params) ActionResult {
	action = strings.Title(action)
	log.Println("in reflect func", module, action)

	code, bs, err := HTTPResp(module, action, vs, data.Values)
	if err != nil {
		return ActionResult{
			Code: code,
			Err:  err,
		}
	}
	in := []reflect.Value{reflect.ValueOf(bs)}
	if !v.IsValid() {
		return ErrActionResult(ErrReflectStructIsNil, module)
	}
	fpp := v.MethodByName(action)
	if !fpp.IsValid() || fpp.IsNil() {
		return ErrActionResult(ErrReflectFuncIsNil, module, action)
	}
	out := fpp.Call(in)
	if len(out) == 0 {
		return ErrActionResult(ErrReflectFuncInvalidReturnValue)
	}
	return ActionResult{
		Code: 1000,
		Err:  nil,
		Data: out[0].Interface(),
	}
}

// ErrActionResult returns ActionResult with error
func ErrActionResult(en int, args ...interface{}) ActionResult {
	return ActionResult{
		Code: en,
		Err:  Err(en, args...),
	}
}

// verifyStatus translate un-formatted string into enable or disable
// enable: enable, 1, online, on
// disable: disable, 0, offline, off
func verifyStatus(st string) string {
	st = strings.ToLower(st)
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
			return StatusEnable
		}
		return StatusDisable
	}
	return StatusUnkown
}
