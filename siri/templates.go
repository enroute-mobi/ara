package siri

import (
	"path/filepath"
	"text/template"

	"bitbucket.org/enroute-mobi/ara/config"
	"bitbucket.org/enroute-mobi/ara/logger"
)

var templates *template.Template

func init() {
	templatePath, err := config.GetTemplateDirectory()
	if err != nil {
		logger.Log.Panicf("Error while loading templates: %v", err)
	}

	templates = template.Must(template.ParseGlob(filepath.Join(templatePath, "*.template")))
}
