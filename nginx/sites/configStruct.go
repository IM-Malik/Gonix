package sites

type Config struct {
    Domain      string   // `json:"domain"` 	// Domain name (e.g., "example.com")
    ListenPort  int      //`json:"listen_port"` // Port to listen on (e.g., 80)
    RootDir     string   //`json:"root_dir"`	// Root directory (e.g., "/var/www/html")
    IndexFiles  string //`json:"index_files`"	// Index files (e.g., ["index.html"])
    ProxyPass   string   //`json:"proxy_pass`"	// Proxy target (e.g., "http://localhost:3000")
    EnableSSL   bool     // Enable SSL
    SSLCertPath string   // SSL certificate path
    SSLKeyPath  string   // SSL private key path
	URLPath		string	 // Add this field for path customization
}

type Upstream struct {
	Name		string
	ServerIP	string
	PortNumber	int
}

func NewConfig() *Config {
	return &Config{
		ListenPort: 80,
		URLPath: "/",
	}
}

func NewUpstream() *Upstream {
	return &Upstream{}
}

const SERVER_BLOCK_TMPL = `server {
    {{- if not .EnableSSL}}
	listen {{.ListenPort}};
	server_name {{.Domain}};
    {{end}}
    {{- if .EnableSSL}}
    listen {{.ListenPort}} ssl;
	server_name {{.Domain}};
	ssl_certificate {{.SSLCertPath}};
	ssl_certificate_key {{.SSLKeyPath}};
    {{end}}
`

const LOCATION_BLOCK_TMPL = `	location {{if .URLPath}}{{.URLPath}}{{else}}/{{end}} {
			{{- if .ProxyPass}}
			proxy_pass {{.ProxyPass}};
			proxy_set_header Host $host;
			proxy_set_header X-Real-IP $remote_addr;
			proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
			proxy_set_header X-Forwarded-Proto $scheme;
			{{- else}}
			return 404;
			{{end}}
		}
`

const UPSTREAM_BLOCK_TMPL = `upstream {{.Name}} {
	server {{.ServerIP}}:{{.PortNumber}};
}`
