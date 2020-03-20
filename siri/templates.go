package siri

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var templates *template.Template

func init() {
	// Small hack to make tests work, otherwise the relative path don't work
	wd, _ := os.Getwd()
	for !strings.HasSuffix(wd, "ara") {
		wd = filepath.Dir(wd)
	}

	templates = template.Must(template.ParseGlob(filepath.Join(wd, "siri/templates/*.template")))
}
