package main

import (
	"github.com/nats-io/nats.go"
	"github.com/nguyenhuuluan434/nats-pub-sub/model"
	"github.com/nguyenhuuluan434/nats-pub-sub/utils"
	"log"
	"strconv"
)

const (
	streamName     = "ORDERS"
	streamSubjects = "ORDERS.*"
	subjectName    = "ORDERS.created"
)

func main() {
	conn, _ := nats.Connect(nats.DefaultURL)
	ctx, err := conn.JetStream()
	utils.CheckErr(err)
	err = createStream(ctx)
	utils.CheckErr(err)
	err = createOrder(ctx)
	utils.CheckErr(err)

}

func createOrder(ctx nats.JetStreamContext) error {
	var order *model.Order
	for i := 0; i < 10; i++ {
		order = model.NewOrder(i, "cust-"+strconv.Itoa(i), "create")
		orderJson, _ := order.ToJson()
		_, err := ctx.Publish(subjectName, orderJson)
		if err != nil {
			log.Println(err)
		}
		log.Printf("order with id %d has been publish", i)
	}
	return nil
}

func createStream(ctx nats.JetStreamContext) error {
	stream, err := ctx.StreamInfo(streamName)
	if err != nil {
		log.Println(err)
	}
	if stream != nil {
		return nil
	}
	log.Printf("create stream %q and subject %q", streamName, subjectName)
	_, err = ctx.AddStream(&nats.StreamConfig{Name: streamName, Subjects: []string{streamSubjects}})
	if err != nil {
		return err
	}
	return nil
}
