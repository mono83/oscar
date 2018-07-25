package tpl

var suite = `
{{ template "header"}}

<section class="card">
	<h1>{{ .Name }}</h1>
</section>

{{ range .Elements }}
<section class="case">
	<h2><a name="{{ .ID }}">{{ .Name }}</a></h3>

	{{ if .Error }}
	<section class="error">
		<h3>Error</h3><pre>{{ .Error }}</pre>
	</section>
	{{ end }}

	{{ if .Logs }}
	<section class="logs">
		<h3>Logs</h3>
		<table class="table-logs" border="0">
			<thead>
				<tr><td>Time</td><td>Level</td><td>Message</td></tr>
			</thead>
			<tbody>
			{{range .Logs}}<tr class="log-{{ .LevelString }}">
			<td class="time">{{ .TimeString }}</td>
			<td class="type">{{ .LevelString }}</td>
			<td class="message">{{ .Message }}</td>
			</tr>
			{{ end }}
			</tbody>
		</table>
	</section>
	{{ end }}

	{{ if .Remotes }}
	<section class="remotes">
		<h3>Remote requests</h2>
		<table class="table-remotes" border="0">
			<thead>
			<tr>
				<td colspan="2">URI</td>
				<td>Success</td>
				<td>Spent time</td>
			</tr>
			</thead>
			<tbody>
			{{range .Remotes}}<tr>
				<td class="type">{{ .Type }}</td>
				<td class="url">{{ .URI }}</td>
				<td>{{ .Success }}</td>
				{{StdElapsed .Elapsed }}
			</tr>
			{{end}}
			</tbody>
		</table>
	</section>
	{{ end }}

	
	{{ if .Variables }}
	<section class="vars">
		<h3>Variables</h2>
		<table class="table-vars" border="0">
		{{range $k, $v := .Variables}}<tr>
		<td class="var">{{ $k }}</td>
		<td class="val">{{ $v }}</td>
		</tr>{{end}}
		</table>
	</section>
	{{ end }}

</section>
{{ end }}

{{ template "footer"}}
`
