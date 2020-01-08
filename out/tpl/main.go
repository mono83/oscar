package tpl

import (
	"html/template"
)

// Load parses and loads all HTML templates
func Load() (*template.Template, error) {
	root := template.New("root")

	root.Funcs(map[string]interface{}{
		"StdElapsed": tableCellElapsed,
		"StdCount":   tableCellCount,
	})

	if _, err := root.New("header").Parse(header); err != nil {
		return nil, err
	}
	if _, err := root.New("footer").Parse(footer); err != nil {
		return nil, err
	}
	if _, err := root.New("summary").Parse(summary); err != nil {
		return nil, err
	}
	if _, err := root.New("remotes").Parse(remotes); err != nil {
		return nil, err
	}
	if _, err := root.New("suite").Parse(suite); err != nil {
		return nil, err
	}
	if _, err := root.New("css").Parse(css); err != nil {
		return nil, err
	}

	return root, nil
}

var header = `
<html>
<head>
	<title>Oscar test results</title>
	<link rel="stylesheet" href="main.css">
</head>
<body>
<section class="navigation">
<a href="summary.html">Summary</a>
Remote requests: <a href="remotes-avg.html">avg</a> <a href="remotes-max.html">max</a> <a href="remotes-sum.html">sum</a>
</section>

`

var footer = `
</body>
`

var css = `
body { margin:0; padding: 1em; }

table {font-size: 11pt}
table thead td {text-align: center; vertical-align: middle; font-weight: bold}
table td {border: 1px solid #666; border-spacing:0; border-collapse: collapse; margin: 0; padding: 2px;}

.table-summary tr.suite td {background-color: #ddd;}
.table-summary tbody tr:hover td {border-color: #9AF; background-color: #79D;}

.elapsed, .count {font: normal 9pt "PT Mono", monospace; text-align: right; padding-right: 0.5em; padding-left: 0.5em;}
.elapsed.empty {text-align: center}

.table-summary .count.failed {font-weight: bold; color: #E30; background-color: #FD9;}
.table-vars td {font: normal 9pt "PT Mono", monospace;}
.table-logs td {font: normal 9pt "PT Mono", monospace;}
.table-logs td.message {white-space:pre;}

.table-logs thead td {font-weight: normal; border: 0; background-color: #03A9F4; color: white;}
.table-logs td {border: 0; margin: 0; padding: 2px; vertical-align: text-top;}
.table-logs .type {text-align:center;}
.table-logs .log-info .type {background-color: #81D4FA;}
.table-logs .log-error .type {background-color: #FABF81;}
.table-logs .log-debug .type {background-color: #EEEEEE;}
.table-logs .log-error .message {font-weight: bold;}
.table-logs .log-info .message {font-weight: bold;}
.table-logs .log-trace td { color: #B0BEC5;}
`
