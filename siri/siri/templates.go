package siri

import (
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"bitbucket.org/enroute-mobi/ara/config"
	"bitbucket.org/enroute-mobi/ara/logger"
)

var (
	stringLt  = []byte("&lt;")
	stringGt  = []byte("&gt;")
	stringAmp = []byte("&amp;")
)

var templates *template.Template

func stringEscape(w io.Writer, b []byte) {
	last := 0
	for i, c := range b {
		var string []byte
		switch c {
		case '<':
			string = stringLt
		case '>':
			string = stringGt
		case '&':
			string = stringAmp
		default:
			continue
		}
		w.Write(b[last:i])
		w.Write(string)
		last = i + 1
	}
	w.Write(b[last:])
}

// see https://github.com/golang/go/blob/8488309192b0ed4b393e2f7b2a93491139ff8ad0/src/text/template/funcs.go#L611
func StringEscape(s string) string {
	if !strings.ContainsAny(s, "<>") {
		return s
	}
	var b strings.Builder
	stringEscape(&b, []byte(s))
	return b.String()
}

func init() {
	templatePath, err := config.GetTemplateDirectory()
	if err != nil {
		logger.Log.Panicf("Error while loading templates: %v", err)
	}

	templates = template.Must(template.New("").
		Funcs(template.FuncMap{"stringEscape": StringEscape}).
		ParseGlob(filepath.Join(templatePath, "*.template")))
}
