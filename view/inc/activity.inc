{{ define "activity" }}

<div class='panel panel-info'>
	<div class="panel-heading"><h3 class="panel-title"><span class="glyphicon glyphicon-user"></span>&nbsp; Recent users</h3></div>
	<table class='table'>
	{{ range $recent := .RecentUsers }}
	<tr>
		<td style="width: 32px;"><img src='http://www.gravatar.com/avatar/{{ md5 $recent.User.Email }}?s=32' class="gravatar-32"></td>
		<td>
			<b>{{ $recent.User.Name }}</b>
			<div>{{ $recent.LoggedInFor }}</div>
		</td>
		<td style="width: 32px;"><img src='/images/flags/{{ $recent.User.Language }}.png'></td>
	</tr>
	{{ end }}
	</table>
	<div class="panel-footer">
		{{ if .IsAdmin }}
			<a class='btn btn-primary btn-sm' href='/users/add'><span class="glyphicon glyphicon-plus"></span> <span class="glyphicon glyphicon-user"></span>&nbsp; Add a user</a> 
		{{ end }}
		<a class='btn btn-link btm-sm' href='/users'><span class="glyphicon glyphicon-user"></span>&nbsp; Users</a>
	</div>
</div>

{{ end }}