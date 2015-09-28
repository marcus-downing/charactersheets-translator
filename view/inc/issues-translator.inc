{{ define "issues-translator" }}

<div class='panel panel-warning'>
	<div class="panel-heading"><h3 class="panel-title"><span class="glyphicon glyphicon-pencil"></span>&nbsp; Translator: {{ .NumTranslatorIssues }} open issues</h3></div>
	<table class="table issues-table">
		{{ $issues := .TranslatorIssues }}
		{{ range $issue := $issues }}
		<tr class="{{ $issue.CssClass }}">
			<th scope="row">{{ $issue.Number }}</th>
			<td style="padding-bottom: 0;"><a href="{{ $issue.URL }}">{{ $issue.Name }}</a>
				<div class="summary">
					{{ $issue.SummaryHTML }}
				</div>
			</td>
			<td><img src="{{ $issue.User.Avatar }}" class="gravatar-32"></td>
		</tr>
		{{ end }}
	</table>
	<div class="panel-footer">
		<a href='https://github.com/marcusatbang/charactersheets-translator/issues/new' class='btn btn-danger btn-sm'><span class='glyphicon glyphicon-fire'></span>&nbsp; File an issue</a> &nbsp; 
		<a href='https://github.com/marcusatbang/charactersheets-translator/issues' class='btn btn-link'><span class='glyphicon glyphicon-folder-open'></span>&nbsp; Issues</a>
	</div>
</div>

{{ end }}