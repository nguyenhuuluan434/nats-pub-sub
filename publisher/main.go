package main

import (
	nats "github.com/nats-io/nats.go"
	"log"
)

const (
	streamName     = "ORDERS"
	streamSubjects = "ORDERS.*"
	subjectName    = "ORDERS.created"
)

func main() {
	/*conn, _ := nats.Connect(nats.DefaultURL)
	ctx, err := conn.JetStream()
	checkErr
*/

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
