package utils

import (
	"bytes"
	"fmt"
	"text/template"
)

func RenderTemplate(tmpl *template.Template, data any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.Bytes(), nil
}
