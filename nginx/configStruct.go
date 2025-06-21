package nginx

// Config struct holds shared configuration parameters for both web and reverse proxy servers
type Config struct {
	ConfigPath	string
	Domain      string 
	ListenPort  int    
	URI	        string 
}

// WebConfig struct holds configuration parameters specific to a web server
type WebConfig struct {
	Config
	StaticContentPath string
	StaticContentFileName string 
}

// RevConfig struct holds configuration parameters specific to a reverse proxy server
type RevConfig struct {
	Config
	ProxyPass   string
	EnableSSL   bool
	SSLCertPath string
	SSLKeyPath  string
	HttpOrHttps string
}

// Upstream struct holds configuration parameters for upstream servers in a reverse proxy setup
type Upstream struct {
	ConfigPath	string
	Name       	string
	ServerIP   	string
	PortNumber 	int
}

// NewRevConfig creates a new instance of RevConfig with default values
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

// NewWebConfig creates a new instance of WebConfig with default values
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

// NewUpstream creates a new instance of Upstream with default values
func NewUpstream() *Upstream {
	return &Upstream{}
}

// SERVER_REVERSEPROXY_BLOCK_TMPL is the template for the server block in a reverse proxy configuration
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

// LOCATION_REVERSEPROXY_BLOCK_TMPL is the template for the location block in a reverse proxy configuration
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

// SERVER_WEBSERVER_BLOCK_TMPL is the template for the server block in a web server configuration
const SERVER_WEBSERVER_BLOCK_TMPL = `server {
	listen		{{.ListenPort}};
	server_name	{{.Domain}};
	error_page	500 502 503 504  /50x.html;
	
`

// LOCATION_WEBSERVER_BLOCK_TMPL is the template for the location block in a web server configuration
const LOCATION_WEBSERVER_BLOCK_TMPL = `	location {{.URI}} {
		root	{{.StaticContentPath}};
		index	{{.StaticContentFileName}};
	}
`

// UPSTREAM_BLOCK_TMPL is the template for the upstream block in a reverse proxy configuration
const UPSTREAM_BLOCK_TMPL = `upstream {{.Name}} {
	server {{.ServerIP}}:{{.PortNumber}};
}
`
