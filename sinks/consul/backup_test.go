package consul

import (
	"testing"
	"time"
)

var obj = ConsulBackup{
	Host:       TEST_CONSUL_Host,
	Exclude:    []string{"goms/leader/.+"},
	StartDate:  time.Now(),
	PrefixPath: "../../backup/consul",
}

func TestBackup(t *testing.T) {
	TestCleanOld(t)
	obj.Backup()
}

func TestReplaceAllKeys(t *testing.T) {
	obj.ReplaceAllKeys()
}

func TestCleanOld(t *testing.T) {
	cron := "5 m"
	obj.CleanOld(cron)
}
