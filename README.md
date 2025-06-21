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

    msg, err := orch.CreateAndEnableWebServer(defaults, "example.com", 80, "/", "/var/www/example", "index.html")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(msg)
}
```

## Permissions

> **NOTE:**  
> Many operations in this library require elevated permissions to modify system files.  
> When running code that uses this library, you should use `sudo go run ...` or run your binary as root to ensure all file and service operations succeed.  
> If you build your program first, remember to run the built binary with `sudo ./yourbinary`—using `sudo` only during the build step is not sufficient.

## Documentation

See GoDoc for full API documentation:  
https://pkg.go.dev/github.com/IM-Malik/Gonix

## License

See [LICENSE.md](LICENSE.md) for details.