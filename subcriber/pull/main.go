package main

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/nguyenhuuluan434/nats-pub-sub/model"
	"github.com/nguyenhuuluan434/nats-pub-sub/utils"
	"log"
	"time"
)

const (
	subSubjectName = "ORDERS.created"
	pubSubjectName = "ORDERS.approved"
	approveStatus  = "approved"
)

func main() {
	conn, err := nats.Connect(nats.DefaultURL)
	utils.CheckErr(err)
	ctx, err := conn.JetStream()
	utils.CheckErr(err)
	subscription, _ := ctx.PullSubscribe(subSubjectName, "order-review", nats.PullMaxWaiting(128))
	timeOutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for {
		select {
		case <-timeOutCtx.Done():
			return
		default:
		}
		//pull batch message for NATS
		msgs, err := subscription.Fetch(5, nats.Context(timeOutCtx))
		if err != nil {
			log.Println(err)
			continue
		}
		var order model.Order
		for _, msg := range msgs {
			//manual acknowledgement to the server
			err := msg.Ack()
			if err != nil {
				log.Println(err)
				continue
			}
			err = json.Unmarshal(msg.Data, &order)
			if err != nil {
				log.Println(err)
				continue
			}
			reviewOrder(ctx, order)
		}
	}
}

func reviewOrder(js nats.JetStreamContext, order model.Order) {
	order.Status = approveStatus
	orderJSON, _ := order.ToJson()
	_, err := js.Publish(pubSubjectName, orderJSON)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Order with OrderID:%d has been %s\n", order.OrderID, order.Status)
}
