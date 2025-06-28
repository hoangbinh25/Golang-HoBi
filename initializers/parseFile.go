package initializers

import (
	"html/template"
	"log"
)

var (
	Tpl *template.Template
	err error
)

func InitParseFiles() {
	Tpl, err = Tpl.ParseGlob("views/**/*.html")
	if err != nil {
		log.Fatal("Error parsing templates: ", err)
	}
}
