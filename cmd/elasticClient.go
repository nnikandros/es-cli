package cmd

import (
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticTypedClient() (*elasticsearch.TypedClient, error) {

	cert, _ := os.ReadFile(os.Getenv("ES_CERT_PATH"))

	cfg := elasticsearch.Config{Addresses: checkHosts(), APIKey: os.Getenv("ES_API_KEY"), CACert: cert}

	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, err
	}

	return es, nil

}

func checkHosts() []string {
	hosts := os.Getenv("ES_HOSTS")

	if b := strings.Contains(hosts, ","); b {
		return strings.Split(hosts, ",")
	}

	return []string{hosts}

}
