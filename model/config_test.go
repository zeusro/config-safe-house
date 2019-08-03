package model

import (
	"log"

	"testing"

	"path"

	"sigs.k8s.io/yaml"
	// "github.com/ghodss/yaml"
)

const (
	consulYAML = `
consul:
- instanceName: "prod"
  instanceURL: "https://consul-example.com"
  backup:
  - file:
      cron: "1m"
      path: "consul/"
      exclude:
      - "/leader/.+"
`
	TEST_CONSUL_Host = ""
	TEST_CONSUL_KEY  = ""
)

func CreateTestObject() Config {
	md := Config{
		Consul: []ConsulConfig{
			{InstanceName: "prod", InstanceURL: TEST_CONSUL_Host},
		},
	}
	return md
}

func TestUnmarshal(t *testing.T) {
	md := Config{}
	err := yaml.Unmarshal([]byte(consulYAML), &md)

	if err != nil {
		log.Fatalf("error: %v", err)
	}
	t.Logf("%#v", md)
}

func TestParseFromPath(t *testing.T) {
	p := path.Join("../", "config.yaml")
	md, err := ParseFromPath(p)

	if err != nil {
		log.Fatalf("error: %v", err)
	}
	t.Logf("%#v", md)
}
