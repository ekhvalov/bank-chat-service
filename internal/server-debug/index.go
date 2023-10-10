package serverdebug

import (
	"html/template"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/ekhvalov/bank-chat-service/internal/logger"
)

type page struct {
	Path        string
	Description string
}

type indexPage struct {
	pages []page
}

func newIndexPage() *indexPage {
	return &indexPage{}
}

func (i *indexPage) addPage(path string, description string) {
	i.pages = append(i.pages, page{
		Path:        path,
		Description: description,
	})
}

func (i *indexPage) handler(eCtx echo.Context) error {
	return template.Must(template.New("index").
		Funcs(template.FuncMap{
			"ToUpper": strings.ToUpper,
		}).Parse(`<html>
	<title>Chat Service Debug</title>
<body>
	<h2>Chat Service Debug</h2>
	<ul>
	{{range .Pages}}
		<li><a href="{{.Path}}">{{.Path}}</a> {{.Description}}</li>
	{{end}}
	</ul>

	<h2>Log Level</h2>
	<form onSubmit="putLogLevel()">
		<select id="log-level-select">
			{{range .LogLevels}}
			<option {{if eq . $.LogLevel}} selected="selected"{{end}} value={{.}}>{{. | ToUpper}}</option>
			{{end}}
		</select>
		<input type="submit" value="Change"></input>
	</form>
	
	<script>
		function putLogLevel() {
			const data = 'level='+document.getElementById('log-level-select').value
			const req = new XMLHttpRequest();
			req.open('PUT', '/log/level', false);
			req.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
			req.setRequestHeader("Content-length", data.length);
			req.setRequestHeader("Connection", "close");
			req.onload = function() { window.location.reload(); };
			req.send(data);
		};
	</script>
</body>
</html>
`)).Execute(eCtx.Response(), struct {
		Pages     []page
		LogLevel  string
		LogLevels []string
	}{
		Pages:     i.pages,
		LogLevel:  logLevel(),
		LogLevels: logger.Levels,
	})
}

func logLevel() string {
	return zap.L().Level().String()
}
