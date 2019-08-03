package model

import (
	"errors"
	"strings"
)

type ConsulConfig struct {
	InstanceName string           `yaml:"instanceName"`
	InstanceURL  string           `yaml:"instanceURL"`
	Backup       []BackupStrategy `yaml:"backup"`
}

func (md *ConsulConfig) CheckSelf() (err error) {
	if strings.TrimSpace(md.InstanceName) == "" {
		err = errors.New("InstanceURL should not be null")
		return err
	}
	if strings.TrimSpace(md.InstanceName) == "" {
		err = errors.New("InstanceName should not be null")
		return err
	}
	if len(md.Backup) == 0 {
		err = errors.New("Backup should not be null")
		return err
	}
	return nil
}
