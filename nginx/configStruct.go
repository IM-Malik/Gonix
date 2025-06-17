package nginx

type Config struct {
	ConfigPath	string
	Domain      string 
	ListenPort  int    
	URI	        string 
}

type WebConfig struct {
	Config
	StaticContentPath string
	StaticContentFileName string 
}

type RevConfig struct {
	Config
	ProxyPass   string
	EnableSSL   bool
	SSLCertPath string
	SSLKeyPath  string
	HttpOrHttps string
}

type Upstream struct {
	ConfigPath	string
	Name       	string
	ServerIP   	string
	PortNumber 	int
}

func NewRevConfig() *RevConfig {
	return &RevConfig{
		Config: Config{
			ListenPort: 80,
			URI:    "/",
		},
		ProxyPass: "127.0.0.1",
		EnableSSL: false,
		HttpOrHttps: "http",
	}
}

func NewWebConfig() *WebConfig {
	return &WebConfig{
		Config: Config{
			ListenPort: 80,
			URI:    "/",	
		},
		StaticContentPath: "/usr/share/nginx/html",
		StaticContentFileName: "index.html index.htm",
	}
}

func NewUpstream() *Upstream {
	return &Upstream{}
}
// Make a function to change the content of the templates somehow
// at the end give the templates a look
const SERVER_REVERSEPROXY_BLOCK_TMPL = `server {
    {{- if not .EnableSSL}}
	listen				{{.ListenPort}};
	server_name			{{.Domain}};
    {{end}}
    {{- if .EnableSSL}}
    listen				{{.ListenPort}} ssl;
	server_name			{{.Domain}};
	ssl_certificate		{{.SSLCertPath}};
	ssl_certificate_key	{{.SSLKeyPath}};
    {{end}}
`

const LOCATION_REVERSEPROXY_BLOCK_TMPL = `	location {{.URI}} {
	{{- if .ProxyPass}}
		proxy_pass			{{.HttpOrHttps}}://{{.ProxyPass}};
		proxy_set_header	Host $host;
		proxy_set_header	X-Real-IP $remote_addr;
		proxy_set_header	X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header	X-Forwarded-Proto $scheme;
		{{- else}}
		return				404;
		{{end}}
	}
`

const SERVER_WEBSERVER_BLOCK_TMPL = `server {
	listen		{{.ListenPort}};
	server_name	{{.Domain}};
	error_page	500 502 503 504  /50x.html;
	
`

const LOCATION_WEBSERVER_BLOCK_TMPL = `	location {{.URLPath}} {
		root	{{.StaticContentPath}};
		index	{{.StaticContentFileName}};
	}
`

const UPSTREAM_BLOCK_TMPL = `upstream {{.Name}} {
	server {{.ServerIP}}:{{.PortNumber}};
}
`
