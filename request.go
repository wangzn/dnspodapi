// @Author: wangzn04@gmail.com
// @Date: 2018-08-31 14:55:56

package dnspodapi

import (
	"log"
	"net/url"
	"reflect"
	"strings"
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
