{{ define "me" }}

<div class='panel panel-primary'>
	<div class="panel-heading"><h3 class="panel-title"><span class="glyphicon glyphicon-user"></span>&nbsp; Your Account</h3></div>
	<div class='panel-body'>
		<h2 style="margin-top: 0"><img src='http://www.gravatar.com/avatar/{{ md5 .CurrentUser.Email }}?s=64' style='box-shadow: 0 3px 3px rgba(0,0,0,0.5);'>&nbsp; {{ .CurrentUser.Name }}</h2>
		{{/*<p>Preferred language: &nbsp; <img src='/images/flags/{{.CurrentUser.Language}}.png'> {{ index .LanguageNames .CurrentUser.Language }}</p>*/}}
		{{/*
		{{ if .CurrentUser.IsLanguageLead }}<p><span class='glyphicon glyphicon-chevron-up'></span>&nbsp; You are language lead for {{ index .LanguageNames .CurrentUser.Language }}</p>{{ end }}
		*/}}
		{{ if .CurrentUser.IsAdmin }}<p><span class='glyphicon glyphicon-star'></span>&nbsp; You are administrator</p>{{ end }}
	</div>
	<div class="panel-footer">
		<a class='btn btn-link btn-sm' href='/account'><span class="glyphicon glyphicon-user"></span>&nbsp; Account settings</a> &nbsp; 
		<a class='btn btn-warning btn-sm' href='/logout'><span class="glyphicon glyphicon-log-out"></span>&nbsp; Log out</a>
	</div>
</div>

{{ end }}