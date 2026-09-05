package main

import (
	"html/template"
	"os"
)

// The {{.}} placeholder means "render the entire context value here."
// If we passed a struct instead of a plain string, we could reach into
// specific fields with {{.FieldName}}.
var x = `
<html>
<body>
Hello {{.}}
</body>
</html>
`

func main() {
	t, err := template.New("hello").Parse(x)
	if err != nil {
		panic(err)
	}
	t.Execute(os.Stdout, "<script>alert('world')</script>")
}
