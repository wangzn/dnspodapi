// @Author: wangzn04@gmail.com
// @Date: 2018-09-06 21:55:31

package dnspodapi

import (
	"testing"
)

func TestParseYaml(t *testing.T) {
	c, err := ParseYamlFile("dpctl/testdata/dnspod.yaml")
	if err != nil {
		t.Errorf("fail to parse yaml file, err: %s", err.Error())
	}
	if c == nil {
		t.Errorf("invalid config, config is nil")
	}
	if _, ok := c.Playbook["scene1"]; !ok {
		t.Errorf("invalid config, missing scene1")
	}
	scene1 := c.Playbook["scene1"]
	if _, ok := scene1[1].Params["record_file"]; !ok {
		t.Errorf("invalid config, missing record_file in Params")
	}
}
