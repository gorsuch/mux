package main

import "bufio"
import "flag"
import "fmt"
import "io"
import "net/http"
import "os"
import "net/url"

func write(ch string) {
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
	fmt.Println("reading")
	resp, err := http.Get("http://localhost:8080?channel=mux")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	rdr := bufio.NewReader(resp.Body)
	for {
		switch line, err := rdr.ReadString('\n'); err {
		case nil:
			fmt.Print(line)
		case io.EOF:
			os.Exit(0)
		default:
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}
	fmt.Println("done")
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
