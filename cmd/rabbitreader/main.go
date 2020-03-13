package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/streadway/amqp"
)

var (
	rabbitDial    string
	rabbitChannel string
	quantity      int
)

var rootCmd = &cobra.Command{
	Use:   "rabbitreader",
	Short: "read from rabbit write on stdout",
	Run:   rabbitreader,
}

func rabbitreader(cmd *cobra.Command, args []string) {
	conn, err := amqp.Dial(rabbitDial)
	if err != nil {
		os.Stderr.WriteString("Cant connect")
		flag.PrintDefaults()
		os.Exit(1)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		os.Stderr.WriteString("No channel")
		flag.PrintDefaults()
		os.Exit(1)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		rabbitChannel, // name
		false,         // durable
		false,         // delete when usused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		os.Stderr.WriteString("No Queue")
		flag.PrintDefaults()
		os.Exit(1)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		msg := fmt.Sprintf("Can't create reader: %v", err)
		os.Stderr.WriteString(msg)
		flag.PrintDefaults()
		os.Exit(1)
	}
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Printf("%s\n", string(d.Body))
			if quantity > 0 {
				quantity--
			}
			if quantity == 0 {
				forever <- false
				break
			}
		}
	}()

	_ = <-forever
}

func main() {
	rootCmd.Flags().StringVar(&rabbitDial, "dial", "", "Dial config, the rabbit from we read data")
	rootCmd.Flags().StringVar(&rabbitChannel, "channel", "", "Channel to read data")
	rootCmd.Flags().IntVar(&quantity, "quantity", -1, "how many read after command kill itself")

	rootCmd.MarkFlagRequired("dial")
	rootCmd.MarkFlagRequired("channel")
	rootCmd.Flags().MarkHidden("quantity")
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}
