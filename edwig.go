package main

import (
	"flag"
	"fmt"

	"github.com/af83/edwig/api"
	"github.com/af83/edwig/siri"
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
		RequestTimestamp:  api.DefaultClock().Now(),
		MessageIdentifier: "CheckStatus:Test:0",
	}
	response, err := client.CheckStatus(request)
	if err != nil {
		return err
	}
	fmt.Println(response)
	return nil
}
