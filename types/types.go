package types

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/elweday/subtitles-go/helpers"
)

type Args struct {
	Words      []helpers.Item
	Index      int
	Percentage float64
}


func NewStyle(name string, t string, funcs map[string]any) func(args Args) string{
	return func(args Args) string {
		tmpl, err := template.New(name).Funcs(funcs).Parse(t)
		if err != nil {
			return fmt.Sprintf("Error parsing template: %v", err)
		}

		var tplOutput bytes.Buffer
		err = tmpl.Execute(&tplOutput, args)
		if err != nil {
			return fmt.Sprintf("Error executing template: %v", err)
		}
		return tplOutput.String()
	}

}