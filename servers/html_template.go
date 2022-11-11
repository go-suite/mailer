package servers

import "html/template"

var html_index = template.Must(template.New("html_index").Parse(`
<html>
	<head>
	</head>
	<body>
		<center>
			<img src="assets/mailer.svg" alt="{{ .title }}">
			<h1>{{ .title }}</h1>
		</center>
	</body>
</html>
`))
