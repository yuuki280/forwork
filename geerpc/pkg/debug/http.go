// pkg/debug/http.go
package debug

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"geerpc/pkg/service"
)

const debugText = `<html>
	<body>
	<title>GeeRPC 服务</title>
	{{range .}}
	<hr>
	服务 {{.Name}}
	<hr>
		<table>
		<th align=center>方法</th><th align=center>调用次数</th>
		{{range $name, $mtype := .Method}}
			<tr>
			<td align=left font=fixed>{{$name}}({{$mtype.ArgType}}, {{$mtype.ReplyType}}) error</td>
			<td align=center>{{$mtype.NumCalls}}</td>
			</tr>
		{{end}}
		</table>
	{{end}}
	</body>
	</html>`

var debug = template.Must(template.New("RPC debug").Parse(debugText))

type debugHTTP struct {
	*service.Service
}

type debugService struct {
	Name   string
	Method map[string]*service.MethodType
}

func (server debugHTTP) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var services []debugService
	server.serviceMap.Range(func(namei, svci interface{}) bool {
		svc := svci.(*service.Service)
		services = append(services, debugService{
			Name:   namei.(string),
			Method: svc.Method,
		})
		return true
	})
	err := debug.Execute(w, services)
	if err != nil {
		_, _ = fmt.Fprintln(w, "rpc: 执行模板错误:", err.Error())
	}
}

func StartDebugHTTP(addr string, path string) {
	http.Handle(path, debugHTTP{})
	log.Printf("rpc: debug HTTP服务启动于 %s%s\n", addr, path)
	go http.ListenAndServe(addr, nil)
}
