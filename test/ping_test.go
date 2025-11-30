package test

import (
	"context"
	"encoding/json"
	"escobra/cmd"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestPing(t *testing.T) {

	es, _ := cmd.NewElasticTypedClient()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := es.Ping().Do(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r)

}

func TestCluster(t *testing.T) {
	es, _ := cmd.NewElasticTypedClient()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// r, err := es.Cluster.Info("http").Do(ctx)

	r, err := es.Cluster.Health().Do(ctx)

	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.MarshalIndent(r, "", " ")

	fmt.Printf("%s", b)

}
