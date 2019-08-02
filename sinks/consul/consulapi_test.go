package consul

import (
	"testing"
)

const (
	TEST_CONSUL_Host = ""
	TEST_CONSUL_KEY  = ""
)

func TestKey(t *testing.T) {
	obj := ConsulAPI{
		Host: TEST_CONSUL_Host,
	}
	value := obj.Key(TEST_CONSUL_KEY)
	t.Log(value)
}

func TestKeys(t *testing.T) {
	obj := ConsulAPI{
		Host: TEST_CONSUL_Host,
	}
	t.Log(obj.Keys())
}
