// Package webserver is responsible for Adding, deleting, modifying, creating/removing symbolic links in web server site files
package webserver

import (
	"fmt"
	"os"
	"text/template"
    "github.com/IM-Malik/Gonix/nginx"
)

// AddSite adds a complete web server site, no need for extra function calling.
func AddSite(availableDirectoryPath string, domain string, listenPort int, uri string, staticContentPath string, staticContentFileName string) (string, error) {
	file, err := os.OpenFile(availableDirectoryPath + domain + ".conf", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
        RemoveSite(availableDirectoryPath, domain)
		return "", fmt.Errorf("failed to create configuration file: %v", err)
	}
	defer file.Close()

    output, err := AddServer(availableDirectoryPath, domain,  listenPort, staticContentPath, staticContentFileName)
    if err != nil {
        RemoveSite(availableDirectoryPath, domain)
        return "", fmt.Errorf("failed to add a site: %v", err)
    }
	return fmt.Sprintf("adding a site is successful: \n%v", output), nil
}

// RemoveSite removes any existing site with the specefied domain name
func RemoveSite(availableDirectoryPath string, domain string) (string, error) {
    return nginx.RemoveSite(availableDirectoryPath, domain)
}

// AddServer adds a server and location blocks to existing site with the specefied domain name
func AddServer(availableDirectoryPath string, domain string, listenPort int, staticContentPath string, staticContentFileName string) (string, error) {
    file, err := os.OpenFile(availableDirectoryPath + domain + ".conf", os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return "", fmt.Errorf("failed to open config file: %v", err)
    }
    defer file.Close()
    
    cfgVars := nginx.NewWebConfig()
    cfgVars.ConfigPath = availableDirectoryPath
    cfgVars.Domain = domain
    cfgVars.ListenPort = listenPort
    cfgVars.StaticContentPath = staticContentPath
    cfgVars.StaticContentFileName = staticContentFileName

    status, err := validateConfigServer(cfgVars)
    if status {
        tmpl := template.Must(template.New("srvBlkTmpl").Parse(nginx.SERVER_WEBSERVER_BLOCK_TMPL))
        if err := tmpl.Execute(file, cfgVars); err != nil {
            return "", fmt.Errorf("server template execution failed: %w", err)
        }
        tmpl2 := template.Must(template.New("locationBlkTmpl").Parse(nginx.LOCATION_WEBSERVER_BLOCK_TMPL))
        if err := tmpl2.Execute(file, cfgVars); err != nil {
            return "", fmt.Errorf("location template execution failed: %w", err)
        }
        file.WriteString("}\n")
        return fmt.Sprintf("creating web server config file is successful: %v", availableDirectoryPath + domain + ".conf"), nil
    }
    return "", fmt.Errorf("web server configuration validation failed: %v", err)
}

// GetEnabledSites return list of all the available sites
func GetAvailableSites(availableDirectoryPath string) ([]os.DirEntry, error) {
    sites, err := nginx.GetSites(availableDirectoryPath)
    if err != nil {
        return nil, err
    } 
    return sites, nil
}

// validateConfigServer checks for all the available information for web server and returns correct error message when missing
func validateConfigServer(cfg *nginx.WebConfig) (bool, error) {
    if cfg.ConfigPath == "" {
        return false, fmt.Errorf("config file path is not set")
    }
    if cfg.Domain == "" {
        return false, fmt.Errorf("domain name is required")
    }
    if cfg.ListenPort == 0 {
        return false, fmt.Errorf("port number needs to be between 1-65535")
    }
    return true, nil
}