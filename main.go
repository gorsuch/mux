package main

import "bufio"
import "fmt"
import "io"
import "github.com/fzzbt/radix/redis"
import "os"

func main() {
	conf := redis.DefaultConfig()
	c := redis.NewClient(conf)
	defer c.Close()

	rdr := bufio.NewReader(os.Stdin)

	for {
		switch line, err := rdr.ReadString('\n'); err {
		case nil:
			reply := c.Publish("stream", line[:len(line) -1])
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