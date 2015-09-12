{{ define "issues-main" }}

<div class='panel panel-warning'>
	<div class="panel-heading"><h3 class="panel-title"><span class="glyphicon glyphicon-pencil"></span>&nbsp; {{ .NumIssues }} open issues</h3></div>
	<table class="table">
		{{ $issues := .Issues }}
		{{ range $issue := $issues }}
		<tr class="{{ $issue.CssClass }}">
			<th scope="row">{{ $issue.Number }}</th>
			<td style="padding-bottom: 0;"><a href="{{ $issue.URL }}">{{ $issue.Name }}</a>
				<div class="summary">
					{{ $issue.SummaryHTML }}
				</div>
			</td>
			<td><img src="{{ $issue.User.Avatar }}" width="24" height="24"></td>
		</tr>
		{{ end }}
	</table>
	<div class="panel-footer">
		<a href='https://github.com/marcusatbang/charactersheets/issues/new'><span class='glyphicon glyphicon-fire'></span>&nbsp; File an issue</a> &nbsp; 
		<a href='https://github.com/marcusatbang/charactersheets/issues'><span class='glyphicon glyphicon-folder-open'></span>&nbsp; Issues</a>
	</div>
</div>

{{ end }}