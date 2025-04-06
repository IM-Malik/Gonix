package sites

type Config struct {
	ConfigPath	string
	Domain      string // `json:"domain"` 	// Domain name (e.g., "example.com") --> Both
	ListenPort  int    //`json:"listen_port"` // Port to listen on (e.g., 80) --> Both
	URLPath     string // Add this field for path customization --> Both
	// RootDir     string //`json:"root_dir"`	// Root directory (e.g., "/var/www/html") --> Reverse Proxy
	// IndexFiles  string //`json:"index_files`"	// Index files (e.g., ["index.html"]) --> Reverse Proxy
	// ProxyPass   string //`json:"proxy_pass`"	// Proxy target (e.g., "http://localhost:3000") --> Reverse Proxy
	// EnableSSL   bool   // Enable SSL --> Reverse Proxy
	// SSLCertPath string // SSL certificate path --> Reverse Proxy
	// SSLKeyPath  string // SSL private key path --> Reverse Proxy
	// StaticContentPath string // Use this field to serve static content --> Web Server
	// StaticContentFileName string // Use this field to write the static content file name --> Web Server
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
			URLPath:    "/",
		},
		ProxyPass: "127.0.0.1",
		EnableSSL: false,
	}
}

func NewWebConfig() *WebConfig {
	return &WebConfig{
		Config: Config{
			ListenPort: 80,
			URLPath:    "/",	
		},
		StaticContentPath: "/usr/share/nginx/html",
		StaticContentFileName: "index.html index.htm",
	}
}

func NewUpstream() *Upstream {
	return &Upstream{}
}

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

const LOCATION_REVERSEPROXY_BLOCK_TMPL = `	location {{.URLPath}} {
	{{- if .ProxyPass}}
		proxy_pass			{{.ProxyPass}};
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
