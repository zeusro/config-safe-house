package consul

import (
	"testing"
	"time"
)

var obj = ConsulBackup{
	Host:      TEST_CONSUL_Host,
	Exclude:   []string{"goms/leader/.+"},
	StartDate: time.Now(),
}

func TestBackup(t *testing.T) {
	TestCleanOld(t)
	obj.Backup("backup/consul")
}

func TestCleanOld(t *testing.T) {
	today := time.Now()
	obj.CleanOld(today)
}
