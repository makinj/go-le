package base_test

import (
	"encoding/json"
	"testing"

	"github.com/makinj/go-le/internal/testutils"
	"github.com/makinj/go-le/modules/base"
)

var TestBaseConfigJSON = []byte(`
	{
		"name":"test"
	}
`)

func UnmarshalConfig(json_conf []byte) (conf base.Config, err error) {
	err = json.Unmarshal(json_conf, &conf)
	return conf, err
}

func TestGetName(t *testing.T) {
	config, err := UnmarshalConfig(TestBaseConfigJSON)
	testutils.Ok(t, err)

	testutils.Equals(t, "test", config.GetName())

}
