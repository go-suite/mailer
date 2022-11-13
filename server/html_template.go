package server

import "html/template"

var html_index = template.Must(template.New("html_index").Parse(`
<!doctype html>
<html lang="en">
	<head>
  		<meta charset="utf-8">
  		<title>{{ .title }}</title>
	</head>
	<body>
		<center>
			<img src="assets/mailer.svg" alt="{{ .title }}">
			<h1>{{ .title }}</h1>
			<p>{{ .secure }}</p>
		</center>
	</body>
</html>
`))
