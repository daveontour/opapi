package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
)

var trace bool
var debug bool
var log bool

const Layout = "2006-01-02T15:04:05"

var rootCmd = &cobra.Command{
	Use:   "rabbitclient.exe {Connection String} {Exchange} {Subscription Topic} {Change Topic} {Log, Debug or Trace}",
	Short: `Various set of utilities`,
	Run: func(cmds *cobra.Command, args []string) {
		if len(args) >= 4 {
			if strings.ToLower(args[4]) == "trace" {
				trace = true
			}
			if strings.ToLower(args[4]) == "debug" {
				debug = true
			}
			if strings.ToLower(args[4]) == "log" {
				log = true
			}
		}
		var forever chan struct{}

		go ListenChanges(args)
		go ListenSub(args)
		<-forever
	},
}

func InitCobra() {

	rootCmd.CompletionOptions.DisableDefaultCmd = false

}
func ExecuteCobra() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func ListenChanges(args []string) {

	connString := args[0]
	exchange := args[1]
	topic := args[3]
	conn, err := amqp.Dial(connString)

	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare the Exchange which the system will try to match if it exists or create if it doesn't
	err = ch.ExchangeDeclare(
		exchange, // exchange
		"topic",  // routing key
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	//the session queue that will receive the messages from the Topic publisher
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare the listening queue")

	fmt.Printf("Binding Changes queue %s to exchange %s with routing key %s\n", q.Name, exchange, topic)

	// Bind the seession queue to the Publisher
	err = ch.QueueBind(
		q.Name,   // queue name
		topic,    // routing key
		exchange, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	// Read the messages from the queue
	go func() {
		i := 1
		for d := range msgs {

			fmt.Println("Rabbit Message Received ", i)
			message := string(d.Body[:])
			changePush(message)
		}
	}()

	fmt.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func ListenSub(args []string) {

	connString := args[0]
	exchange := args[1]
	topic := args[2]
	conn, err := amqp.Dial(connString)

	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare the Exchange which the system will try to match if it exists or create if it doesn't
	err = ch.ExchangeDeclare(
		exchange, // exchange
		"topic",  // routing key
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	//the session queue that will receive the messages from the Topic publisher
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare the listening queue")

	fmt.Printf("Binding Subscription queue %s to exchange %s with routing key %s\n", q.Name, exchange, topic)

	// Bind the seession queue to the Publisher
	err = ch.QueueBind(
		q.Name,   // queue name
		topic,    // routing key
		exchange, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	// Read the messages from the queue
	go func() {
		i := 1
		for d := range msgs {

			fmt.Println("Rabbit Message Received ", i)
			message := string(d.Body[:])
			subPush(message)
		}
	}()

	fmt.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
	}
}
func changePush(jsonData string) {

	path := "./RabbitHooksLogs"
	_ = os.MkdirAll(path, os.ModePerm)

	if debug {
		fmt.Println(string(jsonData[:500]))
	} else if trace {
		fmt.Println(string(jsonData))
	}
	if log {
		file, errs := os.CreateTemp("./RabbitHooksLogs", "changelog-*.json")
		if errs != nil {
			fmt.Println(fmt.Errorf("Error logging to file %s", errs))
			return
		}
		defer file.Close()
		_, err := file.WriteString(string(jsonData))
		if err != nil {
			fmt.Println(fmt.Errorf("Error logging to file %s", err))
			return
		}
	}
	fmt.Printf("Change message received at %s\n", time.Now().Format(Layout))
}

func subPush(jsonData string) {

	path := "./RabbitHooksLogs"
	_ = os.MkdirAll(path, os.ModePerm)

	if debug {
		fmt.Println(string(jsonData[:500]))
	} else if trace {
		fmt.Println(string(jsonData))
	}
	if log {
		file, errs := os.CreateTemp("./RabbitHooksLogs", "pushlog-*.json")
		if errs != nil {
			fmt.Println(fmt.Errorf("Error logging to file %s", errs))
			return
		}
		defer file.Close()
		_, err := file.WriteString(string(jsonData))
		if err != nil {
			fmt.Println(fmt.Errorf("Error logging to file %s", err))
			return
		}
	}
	fmt.Printf("Subcription message received at %s\n", time.Now().Format(Layout))
}
