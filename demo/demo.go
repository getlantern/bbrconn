package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/getlantern/bbrconn"
)

func main() {
	c, err := net.Dial("tcp", "osnews.com:80")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	c.Write([]byte("HTTP/1.1 GET /\r\n"))

	tc, err := bbrconn.Wrap(c)
	if err != nil {
		log.Fatal(err)
	}
	info, err := tc.Info()
	if err != nil {
		log.Fatal(err)
	}

	j, _ := json.MarshalIndent(info, "", "  ")
	fmt.Println(string(j))
}
