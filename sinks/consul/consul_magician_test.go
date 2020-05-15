package consul

import (
	"testing"
)


var m =ConsulMagician{
	Host:       TEST_CONSUL_Host,
}


func TestReplaceAllKeys(t *testing.T) {
	m.ReplaceAllKeys()
}