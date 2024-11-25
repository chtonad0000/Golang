//go:build !solution

package ciletters

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed template.txt
var templateContent string

func add(a, b int) int {
	return a + b
}

func LastLines(log string) string {
	lines := strings.Split(log, "\n")
	if len(lines) > 10 {
		lines = lines[len(lines)-10:]
	}
	return strings.Join(lines, "\n            ")
}

func MakeLetter(n *Notification) (string, error) {
	t := template.Must(template.New("PipelineTemplate").Funcs(template.FuncMap{"lastLines": LastLines}).Funcs(template.FuncMap{"add": add}).Parse(templateContent))
	var buffer bytes.Buffer
	err := t.Execute(&buffer, n)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
