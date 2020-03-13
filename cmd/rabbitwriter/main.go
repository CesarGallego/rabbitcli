package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/streadway/amqp"
)

var (
	rabbitDial    string
	rabbitChannel string
)

var rootCmd = &cobra.Command{
	Use:   "rabbitwriter",
	Short: "write to rabbit from stdin",
	Run:   rabbitwriter,
}

func rabbitwriter(cmd *cobra.Command, args []string) {
	conn, err := amqp.Dial(rabbitDial)
	if err != nil {
		os.Stderr.WriteString("No connection")
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

	var body []byte

	snr := bufio.NewScanner(os.Stdin)
	for snr.Scan() {
		line := snr.Text()
		if len(line) == 0 {
			break
		}

		body = []byte(line)
		pubErr := ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if pubErr != nil {
			e := fmt.Sprintf("%v\n", pubErr)
			os.Stderr.WriteString(e)
		}
	}
	if err := snr.Err(); err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func main() {
	rootCmd.Flags().StringVar(&rabbitDial, "dial", "", "Dial config, the rabbit to we write data")
	rootCmd.Flags().StringVar(&rabbitChannel, "channel", "", "Channel to write data")
	rootCmd.MarkFlagRequired("dial")
	rootCmd.MarkFlagRequired("channel")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
