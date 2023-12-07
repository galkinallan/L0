package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"

	"github.com/nats-io/stan.go"
)

//go:embed model.txt
var data []byte

func main() {
	sc, err := stan.Connect("test-cluster", "publisher", stan.NatsURL("0.0.0.0:4223"))
	if err != nil {
		log.Fatalf("Cannot connect to nats %v\n", err)
	}

	defer sc.Close()

	sData := bytes.Split(data, []byte("\n\n"))

	for i, data := range sData {
		err = sc.Publish("order", data)
		if err != nil {
			fmt.Printf("cannot publish: %v\n", err)
		}
		fmt.Printf("published data %d\n", i+1)
	}

}
