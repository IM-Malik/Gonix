package config

import (
	"os"
	"fmt"
	"text/template"
	"github.com/IM-Malik/Gonix/nginx"
)

type Stream struct {
	MainDomainName			string
	MainDomainConfigPath	string
}

func NewStream() *Stream {
	return &Stream{
		// MainDomainName: 		"malik.com",
		// MainDomainConfigPath: 	"/etc/nginx/sites-available/",
	}
}

const STREAM_BLOCK_TMPL = `stream {
	ssl_preread on;
	map $ssl_preread_server_name $upstream {
		{{.MainDomainName}}	{{.MainDomainConfigPath}}{{.MainDomainName}}.conf;
	}
	server {
		listen 443;
		proxy_pass $upstream;
	}
	upstream {{.MainDomainName}}.conf {
		server {{.ServerIP}}:{{.PortNumber}};
	}
}
`

// Stopped here. Didn't test yet
func GenerateDefaultStreamConfig(mainDomainConfigPath string, mainDomainName string, upstreamServerIP string,  upstreamPortNumber int) (string, error) {
	file, err := os.OpenFile(mainDomainConfigPath + mainDomainName + ".conf", os.O_APPEND|os.O_WRONLY, 0644)
	
	stream := NewStream()
	stream.MainDomainName = mainDomainName
	stream.MainDomainConfigPath = mainDomainConfigPath
	
	upstream, err := nginx.AddUpstream(mainDomainConfigPath, mainDomainName, mainDomainName + ".conf", upstreamServerIP, upstreamPortNumber)
	if err != nil {
		return "", fmt.Errorf("failed to add upstream block: %v", err)
	}
	fmt.Println(upstream)

	upstreamMap := map[string]interface{}{
		"ServerIP": 	upstreamServerIP,
		"PortNumber":	upstreamPortNumber,
	}
	tmpl := template.Must(template.New("streamBlkTmpl").Parse(STREAM_BLOCK_TMPL))
	err = tmpl.Execute(file, stream)
	if err != nil {
		return "", fmt.Errorf("failed to add stream information to template: %v", err)
	}
	err = tmpl.Execute(file, upstreamMap)
	if err != nil {
		return "", fmt.Errorf("failed to add upstream information to template: %v", err)
	}

	return fmt.Sprintf("default stream is generated successfully at: %v", mainDomainConfigPath + "nginx.conf"), nil
}