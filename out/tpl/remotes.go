package tpl

var remotes = `
{{ template "header"}}

<section class="card">
	<h1>Top remote requests</h1>
	<p>This page contains information about remote requests, that consumes most time</p>
</section>

<table class="table-remotes" border="0">
	<thead>
	<tr>
		<td colspan="2" rowspan="2">URI</td>
		<td rowspan="2">Count</td>
		<td colspan="4">Spent time</td>
	</tr>
	<tr>
		<td>Min</td>
		<td>Avg</td>
		<td>Max</td>
		<td>Total</td>
	</td>
	</thead>
	<tbody>
	{{range .}}<tr>
		<td class="type">{{ .Type }}</td>
		<td class="url">{{ .URI }}</td>
		{{StdCount .Count false }}
		{{StdElapsed .Min }}
		{{StdElapsed .Avg }}
		{{StdElapsed .Max }}
		{{StdElapsed .Total }}
	</tr>
	{{end}}
	</tbody>
</table>
{{ template "footer"}}`
