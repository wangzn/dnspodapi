// @Author: wangzn04@gmail.com
// @Date: 2018-08-13 15:35:51

package mymap

import (
	"strconv"
	"strings"
)

// ParseMaps parses a string to a map[string]string
func ParseMaps(str, s1, s2 string) map[string]string {
	res := make(map[string]string)
	p1 := strings.Split(str, s1)
	for _, e := range p1 {
		p2 := strings.Split(e, s2)
		if len(p2) == 2 {
			res[p2[0]] = p2[1]
		} else {
			res[p2[0]] = ""
		}
	}
	return res
}

// StringMustString returns a must value of a map
func StringMustString(m map[string]string, k string) string {
	def := ""
	if v, ok := m[k]; ok {
		return v
	}
	return def
}

// StringMustInt returns a must value of int
func StringMustInt(m map[string]string, k string) int {
	def := 0
	str := StringMustString(m, k)
	if str == "" {
		return def
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return def
	}
	return i
}

// StringMustFloat64 returns a must value of float64
func StringMustFloat64(m map[string]string, k string) float64 {
	def := 0.0
	str := StringMustString(m, k)
	if str == "" {
		return def
	}
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return def
	}
	return f
}
