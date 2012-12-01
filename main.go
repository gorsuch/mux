package main

import "bufio"
import "flag"
import "fmt"
import "io"
import "net/http"
import "os"
import "net/url"

func write(ch string, muxdUrl url.URL) {
	rdr := bufio.NewReader(os.Stdin)
	muxdUrl.RawQuery = "channel=" + ch
	for {
		switch line, err := rdr.ReadString('\n'); err {
		case nil:
			_, err := http.PostForm(muxdUrl.String(), url.Values{"data": {line[:len(line)-1]}})
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

func read(ch string, muxdUrl url.URL) {
	muxdUrl.RawQuery = "channel=" + ch
	resp, err := http.Get(muxdUrl.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	// TODO do something with the status code
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

	muxdUrl, err := url.Parse(os.Getenv("MUXD_URL"))
	if err != nil {
		muxdUrl, _ = url.Parse("http://localhost:8080")
		fmt.Fprintln(os.Stderr, "setting muxdUrl to http://localhost:8080")
	}
	
	if *r == true {
		read(*c, *muxdUrl)
	} else {
		write(*c, *muxdUrl)
	}
}
