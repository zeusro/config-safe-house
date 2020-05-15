package consul

import (
	"testing"
)

var k = ConsulKiller{
	Host:       TEST_CONSUL_Host,
	PrefixPath: "../../backup/consul",
}

func TestCleanOld(t *testing.T) {
	cron := "5 m"
	k.CleanOld(cron)
}
