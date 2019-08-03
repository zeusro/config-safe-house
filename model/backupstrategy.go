package model

type BackupStrategy struct {
	File FileBackupStrategy `yaml:"file"`
}
