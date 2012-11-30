package main

import "bufio"
import "flag"
import "fmt"
import "io"
import "github.com/fzzbt/radix/redis"
import "net/http"
import "os"
import "os/signal"
import "syscall"
import "net/url"

func write(ch string) {
	//resp, err := http.Post("http://localhost:8080/?channel=mux", url.Values{:data })
	rdr := bufio.NewReader(os.Stdin)
	for {
		switch line, err := rdr.ReadString('\n'); err {
		case nil:
			_, err := http.PostForm("http://localhost:8080/?channel=mux", url.Values{"data": {line[:len(line)-1]}})
			if err != nil {
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

func read(ch string) {
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
	sub.Subscribe(ch)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println()
		done <- true
	}()

	<-done
}

func main() {
	r := flag.Bool("r", false, "read mode")
	w := flag.Bool("w", false, "write mode")
	c := flag.String("c", "mux", "channel")
	flag.Parse()

	if (*r == true && *w == true) || (*r == false && *w == false) {
		fmt.Fprintln(os.Stderr, "You need either -r or -w")
		os.Exit(1)
	}

	if *r == true {
		read(*c)
	} else {
		write(*c)
	}
}
