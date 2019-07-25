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
	obj.Backup()
}
