package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/af83/edwig/api"
	"github.com/af83/edwig/siri"
)

func main() {
	uuidPtr := flag.Bool("testuuid", false, "use the test uuid generator")
	clockPtr := flag.String("testclock", "", "use a fake clock at time given. Format 20060102-1504")

	flag.Parse()

	if *uuidPtr {
		api.SetDefaultUUIDGenerator(api.NewFakeUUIDGenerator())
	}
	if *clockPtr != "" {
		testTime, err := time.Parse("20060102-1504", *clockPtr)
		if err != nil {
			panic(err)
		}
		api.SetDefaultClock(api.NewFakeClockAt(testTime))
	}

	if len(flag.Args()) == 0 {
		fmt.Println("usage: edwig [-testuuid] [-testclock=<time>] [-requestor-ref=<requestor>]")
		fmt.Println("\tcheck [-requestor-ref=<requestorRef>] <url>")
		fmt.Println("\tapi")
		os.Exit(1)
	}

	command := flag.Args()[0]

	var err error
	switch command {
	case "check":
		checkFlags := flag.NewFlagSet("check", flag.ExitOnError)
		requestorRefPtr := checkFlags.String("requestor-ref", "Edwig", "Specify requestorRef")
		checkFlags.Parse(flag.Args()[1:])

		err = checkStatus(checkFlags.Args()[0], *requestorRefPtr)
	case "api":
		err = api.NewServer("localhost:8080").ListenAndServe()
	}

	if err != nil {
		if _, ok := err.(*siri.SiriError); !ok {
			panic(err)
		}
		log.Println(err)
		os.Exit(2)
	}

	os.Exit(0)
}

func checkStatus(url string, requestorRef string) error {
	client := siri.NewSOAPClient(url)
	request := &siri.SIRICheckStatusRequest{
		RequestorRef:      requestorRef,
		RequestTimestamp:  api.DefaultClock().Now(),
		MessageIdentifier: "Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC",
	}

	startTime := time.Now()

	xmlResponse, err := client.CheckStatus(request)
	if err != nil {
		return err
	}

	responseTime := time.Since(startTime)

	// Log
	var logMessage []byte
	if xmlResponse.Status() {
		logMessage = []byte("SIRI OK - status true - ")
	} else {
		logMessage = []byte("SIRI CRITICAL: status false - ")
		if xmlResponse.ErrorType() == "OtherError" {
			logMessage = append(logMessage, fmt.Sprintf("%s %d %s - ", xmlResponse.ErrorType(), xmlResponse.ErrorNumber(), xmlResponse.ErrorText())...)
		} else {
			logMessage = append(logMessage, fmt.Sprintf("%s %s - ", xmlResponse.ErrorType(), xmlResponse.ErrorText())...)
		}
	}
	logMessage = append(logMessage, fmt.Sprintf("%.3f seconds response time", responseTime.Seconds())...)
	log.Println(string(logMessage[:]))

	return nil
}
