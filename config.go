// @Author: wangzn04@gmail.com
// @Date: 2018-08-31 14:57:05

package dnspodapi

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// Global defines the `global` section in config file
type Global struct {
	LogOutput string `json:"log_output" yaml:"log_output"`
	Format    string `json:"format" yaml:"format"`
	Auth      string `json:"auth" yaml:"auth"`
}

// AuthEntry defines the auth entry
type AuthEntry struct {
	APIID    int    `json:"api_id" yaml:"api_id"`
	APIToken string `json:"api_token" yaml:"api_token"`
}

// Auth defines the `auth` seciton in config file
type Auth map[string]AuthEntry

// ActionEntry defines the struct of a action
type ActionEntry struct {
	Auth     string            `json:"auth" yaml:"auth"`
	Category string            `json:"category" yaml:"category"`
	Action   string            `json:"action" yaml:"action"`
	Subject  string            `json:"subject" yaml:"subject"`
	Params   map[string]string `json:"params" yaml:"params"`
}

// Scene defines a series entry
type Scene []ActionEntry

// Playbook defines the `playbook` section in config file
type Playbook map[string]Scene

// Config defines the token info for api
type Config struct {
	Global   Global   `json:"global" yaml:"global"`
	Auth     Auth     `json:"auth" yaml:"auth"`
	Playbook Playbook `json:"playbook" yaml:"playbook"`
}

// ParseYamlFile parses yaml config file into Config struct
func ParseYamlFile(f string) (*Config, error) {
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return ParseYamlBytes(bs)
}

// ParseYamlBytes parse yaml config bytes into Config struct
func ParseYamlBytes(bs []byte) (*Config, error) {
	t := Config{}
	err := yaml.Unmarshal(bs, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
