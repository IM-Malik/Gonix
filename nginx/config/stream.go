package config

import (
	"os"
	"fmt"
	"text/template"
	// "github.com/IM-Malik/Gonix/nginx"
)

type Stream struct {
	MainDomainName			string
	// MainDomainConfigName	string
	ServerIP				string
	PortNumber				int
}

func NewStream() *Stream {
	return &Stream{
		// MainDomainName: 		"malik.com",
		// MainDomainConfigPath: 	"/etc/nginx/sites-available/",
	}
}



func GenerateDefaultStreamConfig(globalConfigFilePath string, mainDomainName string, upstreamServerIP string,  upstreamPortNumber int) (string, error) {
	STREAM_BLOCK_TMPL := `
stream {
	ssl_preread on;
	map $ssl_preread_server_name $upstream {
		{{.MainDomainName}}	{{.MainDomainName}}.conf;
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
	file, err := os.OpenFile(globalConfigFilePath + "/nginx.conf", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open the nginx.conf file: %v", err)
	}
	defer file.Close()
	
	stream := NewStream()
	stream.MainDomainName = mainDomainName
	stream.ServerIP = upstreamServerIP
	stream.PortNumber = upstreamPortNumber

	tmpl := template.Must(template.New("streamBlkTmpl").Parse(STREAM_BLOCK_TMPL))
	err = tmpl.Execute(file, stream)
	if err != nil {
		return "", fmt.Errorf("failed to add stream information to template: %v", err)
	}

	return fmt.Sprintf("default stream is generated successfully at: %v", globalConfigFilePath), nil
}