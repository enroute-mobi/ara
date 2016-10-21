package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/af83/edwig/api"
	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

func main() {
	uuidPtr := flag.Bool("testuuid", false, "Use the test uuid generator")
	clockPtr := flag.String("testclock", "", "Use a fake clock at time given. Format 20060102-1504")
	pidPtr := flag.String("pidfile", "", "Write processus pid in given file")
	flag.BoolVar(&logger.Log.Debug, "debug", false, "Enable debug messages")
	flag.BoolVar(&logger.Log.Syslog, "syslog", false, "Redirect messages to syslog")

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("usage: edwig [-testuuid] [-testclock=<time>] [-pidfile=<filename>]")
		fmt.Println("\tcheck [-requestor-ref=<requestorRef>] <url>")
		fmt.Println("\tapi")
		os.Exit(1)
	}

	if *uuidPtr {
		model.SetDefaultUUIDGenerator(model.NewFakeUUIDGenerator())
	}
	if *clockPtr != "" {
		testTime, err := time.Parse("20060102-1504", *clockPtr)
		if err != nil {
			panic(err)
		}
		model.SetDefaultClock(model.NewFakeClockAt(testTime))
	}
	if *pidPtr != "" {
		f, err := os.Create(*pidPtr)
		if err != nil {
			logger.Log.Printf("Error: Unable to create a file at given path")
			os.Exit(2)
		}
		defer f.Close()
		_, err = f.WriteString(strconv.Itoa(os.Getpid()))
		if err != nil {
			panic(err)
		}
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
			logger.Log.Panicf("Error while running: %v", err)
		}
		// Siri errors
		logger.Log.Printf("%v", err)
		os.Exit(2)
	}

	os.Exit(0)
}

func checkStatus(url string, requestorRef string) error {
	client := siri.NewSOAPClient(url)
	request := &siri.SIRICheckStatusRequest{
		RequestorRef:      requestorRef,
		RequestTimestamp:  model.DefaultClock().Now(),
		MessageIdentifier: "Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC",
	}

	startTime := time.Now()

	xmlResponse, err := client.CheckStatus(request)
	if err != nil {
		return err
	}

	responseTime := time.Since(startTime)

	// Logstash
	logstashDatas := make(map[string]string)
	logstashDatas["requestXML"] = request.BuildXML()
	logstashDatas["responseXML"] = checkStatus.RawXML()
	logstashDatas["processingDuration"] = responseTime
	// ...
	err = audit.CurrentLogStash().WriteEvent(logstashDatas)
	if err != nil {
		logger.Log.Panicf("Error while sending datas to Logstash: %v", err)
	}

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
	logger.Log.Printf(string(logMessage))

	return nil
}
