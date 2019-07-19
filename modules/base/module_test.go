package base_test

import (
	"testing"

	"github.com/makinj/go-le/internal/testutils"
)

func TestMakeModule(t *testing.T) {
	_, err := UnmarshalConfig(TestBaseConfigJSON)
	testutils.Ok(t, err)
	//TKTK finish
}
