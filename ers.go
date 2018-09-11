// @Author: wangzn04@gmail.com
// @Date: 2018-08-31 16:43:21

package dnspodapi

import (
	"fmt"
)

const (
	// ErrUnkownError for unkown
	ErrUnkownError = iota + 1000
	// ErrUnkownModule for unregisterd module
	ErrUnkownModule
	// ErrReflectStructIsNil for nil reflect struct value
	ErrReflectStructIsNil
	// ErrReflectFuncIsNil for nil reflect struct func
	ErrReflectFuncIsNil
	// ErrReflectFuncInvalidReturnValue for invalid return value for reflect func
	ErrReflectFuncInvalidReturnValue
	// ErrInvalidTypeAssertion for invalid type assertion
	ErrInvalidTypeAssertion
	// ErrInvalidStatus for invalid status code
	ErrInvalidStatus
)

var ers map[int]string

func init() {
	ers = make(map[int]string)
	ers[ErrUnkownModule] = "unkown module name `%s`"
	ers[ErrReflectFuncIsNil] = "reflect valueof func is nil, module: %s, action: %s"
	ers[ErrReflectStructIsNil] = "reflect valueof struct is nil, module: %s"
	ers[ErrReflectFuncInvalidReturnValue] = "invalid return value, got 0"
	ers[ErrInvalidTypeAssertion] = "invalid type assertion to `%s`"
	ers[ErrInvalidStatus] = "status fail, action: %s, code: %s, msg: %s"
}

// Err returns error
func Err(code int, args ...interface{}) error {
	if f, ok := ers[code]; ok {
		f = fmt.Sprintf("%d | %s", code, f)
		return fmt.Errorf(f, args...)
	}
	return fmt.Errorf("%d | unkown error info", code)
}
