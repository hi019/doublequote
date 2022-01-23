//go:build ignore
// +build ignore

package main

import (
	"log"
	"strings"
	"text/template"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	log.Println("Generating code...")

	err := entc.Generate("./schema", &gen.Config{
		Templates: []*gen.Template{
			gen.MustParse(gen.NewTemplate("with").
				Funcs(template.FuncMap{"title": strings.ToTitle}).
				ParseFiles("template/with.tmpl")),
		},
	})
	if err != nil {
		log.Fatal("running ent codegen:", err)
	}
}
