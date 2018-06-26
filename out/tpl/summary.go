package tpl

var summary = `
{{ template "header"}}

<section class="card">
	<h1>Test case summary</h1>
	<p>This page contains summary data for all test cases</p>
</section>

<table class="table-summary" border="0">
	<thead>
	<tr>
		<td rowspan="2" colspan="2">Name</td>
		<td colspan="2">Assertions</td>
		<td rowspan="2">Remote requests</td>
		<td colspan="3">Spent time</td>
	</tr>
	<tr>
		<td>Failed</td>
		<td>Success</td>
		<td>Total</td>
		<td>HTTP</td>
		<td>Sleep</td>
	</tr>
	</thead>
	<tbody>
	{{range .Flatten}}{{if eq .Type "TestSuite"}}<tr class="suite">{{else if eq .Type "Report"}}<tr class="total">{{else}}<tr>{{end}}
		<td>{{ .Type }}</td>
		<td>
			{{if eq .Type "TestSuite"}}
			<a href="suite-{{ .ID }}.html">{{ .Name }}</a>
			{{else if eq .Type "TestCase"}}
			<a href="suite-{{ .GetParentID }}.html#{{ .ID }}">{{ .Name }}</a>
			{{else}}
			{{ .Name }}
			{{end}}
		</td>
		{{StdCount .CountAssertionsFailedRecursive true}}
		{{StdCount .CountAssertionsSuccessRecursive false}}
		{{StdCount .CountRemoteRequestsRecursive false}}
		{{StdElapsed .Elapsed }}
		{{StdElapsed .ElapsedRemoteRecursive }}
		{{StdElapsed .ElapsedSleepRecursive }}
	</tr>
	{{end}}
	</tbody>
</table>
{{ template "footer"}}
`
