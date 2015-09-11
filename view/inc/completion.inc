{{ define "completion" }}

<div class='panel panel-primary'>
	<div class="panel-heading"><h3 class="panel-title"><span class="glyphicon glyphicon-stats"></span>&nbsp; Completion summary</h3></div>
		
	<table class='table panel-view-form'>
		<tr class='wizard'><th class='active text-center'>Core</th>
			<th class='text-center'>Advanced</th>
			<th class='active text-center'>Further</th>
			<th class='text-center'>Third Party</th></tr>
	</table>

	<table class='table'>
		{{ $languageNames := .LanguageNames }}
		{{ $languageCompletions := .LanguageCompletion }}
		{{ range $lang := .DisplayLanguages }}
		{{ $languageName := index $languageNames $lang }}
		{{ $languageCompletion := index $languageCompletions $lang }}
		<tr>
			<td colspan='4'><img src='/images/flags/{{$lang}}.png'>&nbsp; {{ $languageName }}</td>
		</tr>
		<tr>
			{{ range $complete := $languageCompletion }}
			{{ $col := percentColour $complete }}
			<td>
				<div class="progress">
				  <div class="progress-bar progress-bar-{{$col}}" role="progressbar" aria-valuenow="{{$complete}}" aria-valuemin="0" aria-valuemax="100" style="width: {{$complete}}%;">
				    <span class="sr-only" data-toggle="tooltip" data-placement="bottom" title="{{$complete}}% complete">{{$complete}}% complete</span>
				  </div>
				</div>
			</td>
			{{ end }}
		</tr>
		{{ end }}
	</table>
</div>

{{ end }}