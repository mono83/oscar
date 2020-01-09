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

	{{ if .Events }}
	<section class="logs">
		<h3>Trace</h3>
		<table class="table-logs" border="0">
			<thead>
				<tr><td>Time</td><td>Event</td><td colspan="3">Payload</td></tr>
			</thead>
			<tbody>
			{{range .Events}}
			{{ if eq .TypeString "RemoteRequest" }}
			<tr class="{{ .TypeString }}">
				<td class="time">{{ .TimeString }}</td>
				<td class="type success-{{ .Data.Success }}">{{ .Data.Type }}</td>
				<td class="elapsed">{{ .Data.ElapsedString }}</td>
				<td class="ray">{{ .Data.Ray }}</td>
				<td class="url">{{ .Data.URI }}</td>
			{{ else if eq .TypeString "AssertDone" }}
			<tr class="{{ .TypeString }}">
				<td class="time">{{ .TimeString }}</td>
				<td class="type success-{{ .Data.Success }}">{{ .Data.Operation }}</td>
				{{ if .Data.Success }}
				<td class="expected" title="Expected value" colspan="2">{{ .Data.Expected }}</td>
				{{ else }}
				<td class="expected" title="Expected value">{{ .Data.Expected }}</td>
				<td class="actual" title="Actual value">{{ .Data.Actual }}</td>
				{{ end }}
				<td class="qualifier">{{ .Data.Qualifier }} {{ .Data.Doc }}&nbsp;</td>
			{{ else if eq .TypeString "LogEvent" }}
			<tr class="{{ .TypeString }}-{{ .Data.LevelString }}">
				<td class="time">{{ .TimeString }}</td>
				<td class="type">{{ .Data.LevelString }}</td>
				<td class="message" colspan="3">{{ .Data.Pattern }}</td>
			{{ end }}
			</tr>
			{{ end }}
			</tbody>
		</table>
	</section>
	{{ end }}

	{{ if .Events.RemoteRequests }}
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
			{{range .Events.RemoteRequests}}<tr>
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
