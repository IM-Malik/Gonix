// Package config is responsible for nginx.conf file
package config

import (
	"fmt"
	"github.com/IM-Malik/Gonix/orch"
	"os"
	"text/template"
)

// Stream holds information to be inserted in the stream block template [DEFAULT_STREAM_BLOCK_TMPL]
type Stream struct {
	Domain     string
	ServerIP   string
	PortNumber int
}

// NewStream creates a new instance of the [Stream] struct
func NewStream() *Stream {
	return &Stream{}
}

// DEFAULT_GLOBAL_CONFIGURATION_TMPL holds the immutable value of the default global configuration of nginx.conf
const DEFAULT_GLOBAL_CONFIGURATION = `user  www-data;
worker_processes  auto;
pid        /run/nginx.pid;
include    /etc/nginx/modules-enabled/*.conf;

events {
    worker_connections  768;
    # multi_accept on;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    # Logging settings
    access_log  /var/log/nginx/access.log;
    error_log   /var/log/nginx/error.log warn;

    # Performance optimizations
    sendfile        on;
    tcp_nopush      on;
    tcp_nodelay     on;
    keepalive_timeout  65;
    types_hash_max_size 2048;

    # Gzip compression
    gzip on;
    # gzip_disable "msie6";

    # Include additional configurations (site-specific and extra settings)
    include /etc/nginx/conf.d/*.conf;
    include /etc/nginx/sites-enabled/*;
}
`

// DEFAULT_EMAIL_BLOCK_TMPL holds the immutable value of the default email block of nginx.conf
const DEFAULT_EMAIL_BLOCK = `mail {
	  auth_http 127.0.0.1:9000/cgi-bin/nginxauth.cgi;
      # See sample authentication script at:
      # http://wiki.nginx.org/ImapAuthenticateWithApachePhpScript
      
      # auth_http localhost/auth.php;
      # pop3_capabilities "TOP" "USER";
      # imap_capabilities "IMAP4rev1" "UIDPLUS";
      
      server {
        listen     localhost:110;
              protocol   pop3;
              proxy      on;
      }

      server {
              listen     localhost:143;
              protocol   imap;
              proxy      on;
      }
}
`

// DEFAULT_STREAM_BLOCK_TMPL holds the immutable dynamic template of the default stream block of nginx.conf
const DEFAULT_STREAM_BLOCK_TMPL = `
stream {
    ssl_preread on;
    map $ssl_preread_server_name $upstream {
        {{.Domain}}	{{.Domain}}.conf;
    }
    server {
        listen 443;
        proxy_pass $upstream;
    }
    upstream {{.Domain}}.conf {
        server {{.ServerIP}}:{{.PortNumber}};
    }
}
`

// GenerateDefaultGlobalConfig generate a default nginx.conf configuration based on the template
func GenerateDefaultGlobalConfig(defaults *orch.Defaults) (string, error) {
	file, err := os.OpenFile(defaults.NginxConf+"nginx.conf", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open the nginx.conf file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(DEFAULT_GLOBAL_CONFIGURATION)
	if err != nil {
		return "", fmt.Errorf("failed to write in the nginx.conf file: %v", err)
	}

	return "the default configuration is written correctly in nginx.conf file\n", nil
}

// GenerateDefaultEmailConfig generates a default mail block with mail servers using pop3 and imap based on the template (works but the template is not finalized)
func GenerateDefaultEmailConfig(defaults *orch.Defaults) (string, error) {
	file, err := os.OpenFile(defaults.NginxConf+"nginx.conf", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open the nginx.conf file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(DEFAULT_EMAIL_BLOCK)
	if err != nil {
		return "", fmt.Errorf("failed to write in the nginx.conf file: %v", err)
    }
	return "the email default configuration is written correctly in nginx.conf file\n", nil
}

// GenerateDefaultStreamConfig generates a default stream block
func GenerateDefaultStreamConfig(defaults *orch.Defaults, domain string, upstreamServerIP string, upstreamPortNumber int) (string, error) {
	file, err := os.OpenFile(defaults.NginxConf + "nginx.conf", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open the nginx.conf file: %v", err)
	}
	defer file.Close()

	stream := NewStream()
	stream.Domain = domain
	stream.ServerIP = upstreamServerIP
	stream.PortNumber = upstreamPortNumber

	tmpl := template.Must(template.New("streamBlkTmpl").Parse(DEFAULT_STREAM_BLOCK_TMPL))
	err = tmpl.Execute(file, stream)
	if err != nil {
		return "", fmt.Errorf("failed to add stream information to template: %v", err)
	}
	return fmt.Sprintf("default stream is generated successfully at: %v", defaults.NginxConf + "nginx.conf"), nil
}
