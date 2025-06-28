package initializers

import (
	"log"
	"text/template"
)

var (
	tpl *template.Template
	err error
)

func InitParseFiles() {
	tpl = template.New("")
	tpl, err = tpl.ParseGlob("views/**/*.html")
	tpl, err = tpl.ParseGlob("views/*.html")
	if err != nil {
		log.Fatal("Error parsing templates: ", err)
	}
}
