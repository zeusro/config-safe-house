package consul

import (
	"testing"
)

func TestBackup(t *testing.T) {
	obj := ConsulBackup{
		Host:"",
	}
	obj.Backup()
	// t.Log(string(text))
	// assert.True(t, m != nil)
}