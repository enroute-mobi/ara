package main

import (
	"flag"
	"fmt"
	"github.com/af83/edwig/siri"
	"time"
)

func main() {
	flag.Parse()

	command := flag.Args()[0]

	var err error
	switch command {
	case "check":
		err = checkStatus(flag.Args()[1])
	}

	if err != nil {
		panic(err)
	}
}

func checkStatus(url string) error {
	client := siri.NewSOAPClient(url)
	request := &siri.SIRICheckStatusRequest{
		RequestorRef:      "NINOXE:default",
		RequestTimestamp:  time.Now(), // FIXME Use clock. See #1735
		MessageIdentifier: "CheckStatus:Test:0",
	}
	response, err := client.CheckStatus(request)
	if err != nil {
		return err
	}
	fmt.Println(response)
	return nil
}
