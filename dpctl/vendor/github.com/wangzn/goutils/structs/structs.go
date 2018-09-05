// @Author: wangzn04@gmail.com
// @Date: 2017-08-10 21:22:52

package structs

import (
	"fmt"

	"github.com/fatih/structs"
)

// StructKeys returns all field names in struct
// b is a boolean for print embedded or not
func StructKeys(v interface{}, b bool) []string {
	s := structs.New(v)
	ks := make([]string, 0)
	for _, f := range s.Fields() {
		if f.IsEmbedded() {
			if b {
				tmp := StructKeys(f.Value(), b)
				ks = append(ks, tmp...)
			}
		} else {
			ks = append(ks, f.Name())
		}
	}
	return ks
}

// StructValues returns all values in struct
// b is a boolean for print embedded or not
func StructValues(v interface{}, b bool) []string {
	vs := make([]string, 0)
	s := structs.New(v)
	for _, f := range s.Fields() {
		if f.IsEmbedded() {
			if b {
				tmp := StructValues(f.Value(), b)
				vs = append(vs, tmp...)
			}
		} else {
			vs = append(vs, fmt.Sprintf("%v", f.Value()))
		}
	}
	return vs
}
