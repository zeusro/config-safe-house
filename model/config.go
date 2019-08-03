package model

import (
	"io/ioutil"

	"sigs.k8s.io/yaml"
)

type Config struct {
	Consul []ConsulConfig
}

func ParseFromPath(path string) (*Config, error) {
	md := Config{}
	content, err := ioutil.ReadFile(path)
	// fmt.Println(string(content))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, &md)
	if err != nil {
		return nil, err
	}
	return &md, err
}
