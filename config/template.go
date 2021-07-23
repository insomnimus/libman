package config

const HelpTemplate = `{{.Name}}: {{.Usage}}
Usage:
  {{.Name}} [OPTIONS]
or
  {{.Name}} <SUBCOMMAND>

Options for the base command:
  {{range $index, $option := .VisibleFlags}}{{if $index}}
  {{end}}{{$option}}{{end}}{{if .VisibleCommands}}

Subcommands:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
  {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{end}}
`
