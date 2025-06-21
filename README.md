# Gonix

Gonix is a Go library for automating Nginx configuration and management. It provides high-level functions for creating, enabling, updating, and removing Nginx site configurations, as well as managing modules and global settings.

## Features

- Create and enable Nginx web server and reverse proxy sites
- Update and rollback site configurations safely
- Manage Nginx modules (enable/disable)
- Generate default Nginx configuration files
- Backup and restore site configurations
- Utility functions for reloading, restarting, and testing Nginx

## Requirements

- Go 1.18 or later
- Linux system with Nginx installed
- **Root privileges are required** for most operations (see below)

## Installation

```sh
go get github.com/IM-Malik/Gonix
```

## Usage

```go
import "github.com/IM-Malik/Gonix/orch"
```

### Example

```go
package main

import (
    "fmt"
    "github.com/IM-Malik/Gonix/orch"
)

func main() {
    defaults := &orch.Defaults{
        NginxConf:      "/etc/nginx/",
        SitesAvailable: "/etc/nginx/sites-available/",
        SitesEnabled:   "/etc/nginx/sites-enabled/",
        ModulesEnabled: "/etc/nginx/modules-enabled/",
    }

    msg, err := orch.CreateAndEnableRevProxy(
        defaults,
        "example.com",     // domain
        80,                // listen port
        "/",               // URI path
        false,             // EnableSSL
        "",                // SSLCertPath
        "",                // SSLKeyPath
        "backend",         // upstreamName
        "127.0.0.1",       // serverIP
        8080,              // portNum
        "http",            // httpOrHttps
    )
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(msg)
}
```

### Output

If the operation is successful, you will see output similar to:

```
the site was created and enabled successfully: example.com
```

### What happens to Nginx

A file `/etc/nginx/sites-available/example.com.conf` will be created with content similar to:

```nginx
server {
    listen 80;
    server_name example.com;
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

A symbolic link to this file will also be created in `/etc/nginx/sites-enabled/`, enabling the site in Nginx.


## Advanced Usage

While most users will only need the high-level `orch` package, Gonix also exposes lower-level packages for more granular control over Nginx configuration and management. You can import and use these packages directly if you need to customize or extend functionality beyond what `orch` provides.

### Example: Using the `nginx/sites/reverseproxy` Package Directly

```go
package main

import (
    "fmt"
    "github.com/IM-Malik/Gonix/nginx/sites/reverseproxy"
)

func main() {
    // Add a new reverse proxy site configuration directly
    msg, err := reverseproxy.AddSite(
        "/etc/nginx/sites-available/", // Or you could use a Defaults instance
        "example.com",
        80,
        "127.0.0.1:8080",
        "/",
        false, // EnableSSL
        "",    // SSLCertPath
        "",    // SSLKeyPath
        "http",
    )
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(msg)
}
```

### Output

```
adding a site is successful: 
creating config file is successful: /etc/nginx/sites-available/example.com.conf
```

### What happens to Nginx

A file `/etc/nginx/sites-available/example.com.conf` will be created with content similar to:

```nginx
server {
    listen 80;
    server_name example.com;
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

In this scenarion you will need to manually enable the site by creating a symbolic link in `/etc/nginx/sites-enabled/`.

## Permissions

> **NOTE:**  
> Many operations in this library require elevated permissions to modify system files.  
> When running code that uses this library, you should use `sudo go run ...` or run your binary as root to ensure all file and service operations succeed.  
> If you build your program `sudo go build ...` first, remember to run the built binary with `sudo ./yourbinary`—using `sudo` only during the build step is not sufficient.



### Other Available Packages

- `github.com/IM-Malik/Gonix/nginx` – Core Nginx configuration structures and helpers
- `github.com/IM-Malik/Gonix/nginx/sites/webserver` – Web server site configuration management
- `github.com/IM-Malik/Gonix/nginx/sites/reverseproxy` – Reverse proxy site configuration management
- `github.com/IM-Malik/Gonix/nginx/modules` – Nginx module management

See the [Documentation](#documentation) section below for more details.

## Documentation

See GoDoc for full API documentation:  
https://pkg.go.dev/github.com/IM-Malik/Gonix

## License

See [LICENSE.md](LICENSE.md) for details.