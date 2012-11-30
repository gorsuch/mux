package main

import "bufio"
import "flag"
import "fmt"
import "io"
import "github.com/fzzbt/radix/redis"
import "os"
import "os/signal"
import "syscall"

func stream() {
	conf := redis.DefaultConfig()
	c := redis.NewClient(conf)
	defer c.Close()

	rdr := bufio.NewReader(os.Stdin)
	for {
		switch line, err := rdr.ReadString('\n'); err {
		case nil:
			reply := c.Publish("stream", line[:len(line)-1])
			if reply.Err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}

		case io.EOF:
			os.Exit(0)

		default:
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}
}

func tap() {
	done := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)

	conf := redis.DefaultConfig()
	c := redis.NewClient(conf)
	defer c.Close()

	h := func(msg *redis.Message) {
		switch msg.Type {
		case redis.MessageMessage:
			fmt.Println(msg.Payload)
		}
	}

	sub, err := c.Subscription(h)
	if err != nil {
		panic(err)
	}

	defer sub.Close()
	sub.Subscribe("stream")

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done
}

func main() {
	s := flag.Bool("s", false, "stream mode")
	t := flag.Bool("t", false, "tap mode")
	flag.Parse()

	if (*s == true && *t == true) || (*s == false && *t == false) {
		fmt.Fprintln(os.Stderr, "You us either -s or -t")
		os.Exit(1)
	}

	if *s == true {
		stream()
	} else {
		tap()
	}
}
