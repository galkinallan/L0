package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/galkinallan/L0/internal/cache"
	"github.com/galkinallan/L0/internal/db"
	"github.com/galkinallan/L0/internal/models"
	"github.com/galkinallan/L0/pkg/service"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
)

func main() {
	ctx := context.Background()

	storage := db.MustCreateDb(ctx)
	defer storage.Close()

	cache := cache.CreateCache()

	fmt.Println("Filling cache from db")
	storage.FillCache(context.Background(), cache)
	cache.PrintKeys()

	service := service.CreateOrderService(storage, cache)

	sc, err := stan.Connect("test-cluster", "subscriber", stan.NatsURL("0.0.0.0:4223"))
	defer sc.Close()

	if err != nil {
		log.Fatal(err)
	}

	sub, err := sc.Subscribe("order", func(m *stan.Msg) {
		var order models.Order
		err := json.Unmarshal(m.Data, &order)
		if err != nil {
			fmt.Printf("not valid json %v\n", err)
			return
		}
		fmt.Printf("Order recieved: %#v\n", order.OrderUID)

		err = validator.New().Struct(&order)
		if err != nil {
			fmt.Printf("Not valid json\n")
			return
		}

		err = service.CreateOrder(ctx, order.OrderUID, order)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	})

	defer sub.Unsubscribe()

	http.HandleFunc("/", service.GetHttpHandle())

	http.ListenAndServe(":8080", nil)

}
