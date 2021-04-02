package api

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type ImportController struct {
	referential   *core.Referential
	importRequest ImportRequest
	csvReader     *bytes.Buffer
}

type ImportRequest struct {
	Force bool
}

func NewImportController(referential *core.Referential) ControllerInterface {
	return &ImportController{
		referential: referential,
		csvReader:   new(bytes.Buffer),
	}
}

func (controller *ImportController) serve(response http.ResponseWriter, request *http.Request, requestData *RequestData) {
	mediaType, params, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil || !strings.HasPrefix(mediaType, "multipart/") {
		http.Error(response, "Expected multipart content", http.StatusUnsupportedMediaType)
		return
	}
	mr := multipart.NewReader(request.Body, params["boundary"])
	for i := 0; i < 2; i++ {
		p, err := mr.NextPart()
		if err != nil {
			http.Error(response, "Can't parse multipart content", http.StatusBadRequest)
			return
		}

		switch p.FormName() {
		case "request":
			jsonDecoder := json.NewDecoder(p)
			jsonDecoder.Decode(&controller.importRequest)
			if err != nil {
				http.Error(response, "Can't parse JSON multipart content", http.StatusBadRequest)
				return
			}
		case "data":
			io.Copy(controller.csvReader, p)
		default:
			http.Error(response, "Wrong multipart content", http.StatusBadRequest)
			return
		}
	}

	stime := controller.referential.Clock().Now()

	result := model.NewLoader(string(controller.referential.Slug()), controller.importRequest.Force, false).Load(controller.csvReader)
	logger.Log.Debugf("ImportController Load time : %v", controller.referential.Clock().Since(stime))

	jsonBytes, _ := json.Marshal(result)
	logger.Log.Debugf("ImportController Json Marshal time : %v ", controller.referential.Clock().Since(stime))
	response.Write(jsonBytes)
}
