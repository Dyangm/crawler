package config

import (
	"fmt"
	"testing"
)

func TestNewConfig(t *testing.T) {
	c, err := GetConfig()
	if err != nil {
		t.Fail()
		return
	}

	fmt.Println(c)
}
