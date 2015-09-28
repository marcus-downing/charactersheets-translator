{{ define "me" }}

<div class='panel panel-primary'>
	<div class="panel-heading"><h3 class="panel-title"><span class="glyphicon glyphicon-user"></span>&nbsp; You</h3></div>
	<div class='panel-body'>
		<h2 style="margin-top: 0"><img src='http://www.gravatar.com/avatar/{{ md5 .CurrentUser.Email }}?s=64' class="gravatar-64 gravatar-shadow">&nbsp; {{ .CurrentUser.Name }}</h2>
		{{/*<p>Preferred language: &nbsp; <img src='/images/flags/{{.CurrentUser.Language}}.png'> {{ index .LanguageNames .CurrentUser.Language }}</p>*/}}
		{{/*
		{{ if .CurrentUser.IsLanguageLead }}<p><span class='glyphicon glyphicon-chevron-up'></span>&nbsp; You are language lead for {{ index .LanguageNames .CurrentUser.Language }}</p>{{ end }}
		*/}}
		{{ if .CurrentUser.IsAdmin }}<p><span class='glyphicon glyphicon-star'></span>&nbsp; You are administrator</p>{{ end }}
	</div>
	<div class="panel-footer">
		<div><a href='http://www.gravatar.com' class='btn btn-link' id='user-avatar'><span class="glyphicon glyphicon-link"></span>&nbsp; Set my avatar at Gravatar.com</a></div>
		<a class='btn btn-link' href='/account'><span class="glyphicon glyphicon-user"></span>&nbsp; Account settings</a> &nbsp; 
		<a class='btn btn-warning btn-sm' href='/logout'><span class="glyphicon glyphicon-log-out"></span>&nbsp; Log out</a>
	</div>
</div>

{{ end }}