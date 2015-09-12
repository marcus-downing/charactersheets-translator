{{ define "summary" }}

<div class='panel panel-info'>
	<div class='panel-heading'><h3 class='panel-title'><span class="glyphicon glyphicon-globe"></span>&nbsp; 
		{{ if .CurrentUser.IsAdmin }}Translation summary{{ else }}Your translations{{ end }}</h3></div>

	<form class='navbar-form panel-view-form'>
		
			<img src='/images/flags/{{.CurrentUser.Language}}.png'> &nbsp;{{ index .LanguageNames .CurrentUser.Language }}
		
	</form>

	{{ $currentUser := .CurrentUser }}
	{{ range $profile := profileTranslations .CurrentUser }}
		<div class='panel-body profile'>
				<div class=''>
					<p>{{ $profile.LevelName }}: <b>{{ $profile.CompletedPercent }}%</b></p>
					<div class='progress'>
					<div class='progress-bar progress-bar' style='width: {{ $profile.CompletedPercent }}%'>
					</div>
					</div>
					<p class='lang-counts'>
						{{ if $profile.ByMeAlone }}
							<a class='lang-count lang-count-mine' href='/translate?level={{ $profile.Level }}&amp;show=mine'>
								<span class='badge'>{{ $profile.ByMeAlone }}</span> by me alone
							</a>
						{{ end }}
						{{ if $profile.ByOthersAlone }}
							<a class='lang-count lang-count-other' href='/translate?level={{ $profile.Level }}&amp;show=others'>
								<span class='badge'>{{ $profile.ByOthersAlone }}</span> by others
							</a>
						{{ end }}
						{{ if $profile.ByMeAndOthersNoConflict }}
							<a class='lang-count lang-count-both' href='/translate?level={{ $profile.Level }}&amp;show=both'>
								<span class='badge'>{{ $profile.ByMeAndOthersNoConflict }}</span> by me and others
							</a>
						{{ end }}
						{{ if $profile.ByMeAndOthersConflict }}
							<a class='lang-count lang-count-conflict' href='/translate?level={{ $profile.Level }}&amp;show=conflicts'>
								<span class='badge'>{{ $profile.ByMeAndOthersConflict }}</span> conflicts
							</a>
						{{ end }}
					</p>
					<div class="progress">
						<div class="progress-bar progress-bar-success" style="width: {{ $profile.ByMePercent }}%">
							<span class="sr-only">{{ $profile.ByMePercent }}% completed by me</span>
						</div>
						<div class="progress-bar progress-bar-info" style="width: {{ $profile.ByOthersPercent }}%">
							<span class="sr-only">{{ $profile.ByOthersPercent }}% completed by others</span>
						</div>
						<div class="progress-bar progress-bar-warning" style="width: {{ $profile.ByMeAndOthersNoConflictPercent }}%">
							<span class="sr-only">{{ $profile.ByMeAndOthersNoConflictPercent }}% complete by me and others</span>
						</div>
						<div class="progress-bar progress-bar-danger" style="width: {{ $profile.ByMeAndOthersConflictPercent }}%">
							<span class="sr-only">{{ $profile.ByMeAndOthersConflictPercent }}% conflicts</span>
						</div>
					</div>
				</div>
		</div>
	{{ end }}

	<div class="panel-footer">
		<a class='btn btn-success' href='/translate'><span class="glyphicon glyphicon-globe"></span>&nbsp; Translate</a> &nbsp;
		<a class='btn btn-link' href='/sources'><span class="glyphicon glyphicon-book"></span>&nbsp; Sources</a>
		<a class='btn btn-link' href='/entries'><span class="glyphicon glyphicon-certificate"></span>&nbsp; Entries</a>
	</div>
</div>

{{ end }}