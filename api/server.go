package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/af83/edwig/siri"
)

type Server struct {
	UUIDConsumer
	ClockConsumer

	bind        string
	startedTime time.Time
}

func NewServer(bind string) *Server {
	server := Server{bind: bind}
	server.startedTime = server.Clock().Now()
	return &server
}

func (server *Server) ListenAndServe() error {
	http.HandleFunc("/siri", server.checkStatusHandler)
	fmt.Printf("Starting server on %s\n", server.bind)
	return http.ListenAndServe(server.bind, nil)
}

func (server *Server) checkStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Try to read and parse request body
	requestContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request: can't read content", 500)
		return
	}
	xmlRequest, err := siri.NewXMLCheckStatusRequestFromContent(requestContent)
	if err != nil {
		http.Error(w, "Invalid request: can't parse content", 500)
		return
	}

	fmt.Printf("CheckStatus %s\n", xmlRequest.MessageIdentifier())

	// Set Content-Type header and create a SIRICheckStatusResponse
	w.Header().Set("Content-Type", "text/xml")

	response := new(siri.SIRICheckStatusResponse)
	response.Address = strings.Join([]string{r.URL.Host, r.URL.Path}, "")
	response.ProducerRef = "Edwig"
	response.RequestMessageRef = xmlRequest.MessageIdentifier()
	response.ResponseMessageIdentifier = fmt.Sprintf("Edwig:ResponseMessage::%s:LOC", server.NewUUID())
	response.Status = true // Temp
	response.ResponseTimestamp = server.Clock().Now()
	response.ServiceStartedTime = server.startedTime

	fmt.Fprintf(w, siri.WrapSoap(response.BuildXML()))
}
