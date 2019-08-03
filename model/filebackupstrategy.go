package model

type FileBackupStrategy struct {
	Cron        string   `yaml:"cron"`
	Exclude     []string `yaml:"exclude"`
	Path        string   `yaml:"path"`
	CleanPolicy string   `yaml:"cleanPolicy"`
}
