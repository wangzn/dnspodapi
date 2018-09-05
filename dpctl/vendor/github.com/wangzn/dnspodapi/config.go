// @Author: wangzn04@gmail.com
// @Date: 2018-08-31 14:57:05

package dnspodapi

import (
	"fmt"
)

// Config defines the token info for api
type Config struct {
	token   string
	appid   int
	appname string
}

// LoginToken returns the string used in api URL parameters
func (c *Config) LoginToken() string {
	return fmt.Sprintf("%d,%s", c.appid, c.token)
}

// Token returns token
func (c *Config) Token() string {
	return c.token
}

// Appid returns appid
func (c *Config) Appid() int {
	return c.appid
}
