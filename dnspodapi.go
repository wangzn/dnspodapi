// @Author: wangzn04@gmail.com
// @Date: 2018-08-31 14:53:59

package dnspodapi

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/wangzn/goutils/myhttp"
)

const (
	// DefaultEndpoint defines the default endpoint URL of dnspodapi
	DefaultEndpoint = "https://dnsapi.cn"
	// DefaultFormat defines the default format of response
	DefaultFormat = "json"
)

var (
	endpoint = DefaultEndpoint
	// json only currently
	format   = DefaultFormat
	apiToken = ""
)

var (
	mu           sync.Mutex
	reflectFuncs map[string]ReflectFunc
)

// Params defines the common data type used in dnspodapi
type Params struct {
	url.Values
}

// P init a Params
func P() Params {
	return Params{url.Values{}}
}

// ActionResult defines the result struct of a action
type ActionResult struct {
	Code int `json:"code"`
	// Data map[string]string `json:"data"`
	Data interface{} `json:"data"`
	Err  error
}

// ReflectFunc declares the func format to reflect action to different modules
type ReflectFunc func(string, Params) ActionResult

// DPA defines the struct for api related info
type DPA struct {
	config *Config
}

func init() {
	reflectFuncs = make(map[string]ReflectFunc)
}

// SetEndpoint set the URL endpoint
func SetEndpoint(ep string) {
	endpoint = ep
}

// Register registeres a func with module
func Register(module string, rf ReflectFunc) error {
	mu.Lock()
	defer mu.Unlock()
	if reflectFuncs == nil {
		reflectFuncs = make(map[string]ReflectFunc)
	}
	if _, ok := reflectFuncs[module]; ok {
		return fmt.Errorf("module %s is registered", module)
	}
	reflectFuncs[module] = rf
	return nil
}

// New returns a new DPA pointer
func New(appid int, token string) *DPA {
	c := &Config{
		appid: appid,
		token: token,
	}
	SetAPIToken(appid, token)
	return &DPA{
		config: c,
	}
}

// SetAPIToken set appid and token
func SetAPIToken(appid int, token string) {
	apiToken = fmt.Sprintf("%d,%s", appid, token)
}

// Action runs a command
func Action(module, action string, data Params) ActionResult {
	// TODO: impl
	mu.Lock()
	defer mu.Unlock()
	if f, ok := reflectFuncs[module]; ok {
		return f(action, data)
	}
	return ActionResult{
		Code: ErrUnkownModule,
		Err:  Err(ErrUnkownModule, module),
	}
}

// HTTPResp returns the http response
func HTTPResp(module, action string, vs url.Values, data url.Values) (int, []byte, error) {
	ep := fmt.Sprintf("%s/%s", endpoint, ModuleActionString(module, action))
	if data == nil {
		data = url.Values{}
	}
	if vs == nil {
		vs = url.Values{}
	}
	data.Add("login_token", apiToken)
	data.Add("format", format)
	log.Println("send post request", ep, vs, data.Encode())
	return myhttp.HPost(ep, vs, []byte(data.Encode()), "")
}

// ModuleActionString returns the string to join module and action for URL
func ModuleActionString(module, action string) string {
	return fmt.Sprintf("%s.%s", strings.Title(module), strings.Title(action))
}
