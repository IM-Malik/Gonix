// Package webserver is responsible for Adding, deleting, modifying, creating/removing symbolic links in web server site files
package webserver

import (
	"fmt"
	"os"
	"text/template"
    "github.com/IM-Malik/Gonix/nginx"
)

// Function AddSite adds a complete web server site, no need for extra function calling.
func AddSite(directoryPath string, domain string, listenPort int, uri string, staticContentPath string, staticContentFileName string) (string, error) {
	file, err := os.OpenFile(directoryPath + domain + ".conf", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
        RemoveSite(directoryPath, domain)
		return "", fmt.Errorf("failed to create configuration file: %v", err)
	}
	defer file.Close()

    output, err := AddServer(directoryPath, domain,  listenPort)
    if err != nil {
        RemoveSite(directoryPath, domain)
        return "", fmt.Errorf("failed to add a site: %v", err)
    }
	return fmt.Sprintf("adding a site is successful: \n%v", output), nil
}

// Function RemoveSite removes any existing site with the specefied domain name
func RemoveSite(directoryPath string, domain string) (string, error) {
    return nginx.RemoveSite(directoryPath, domain)
}

// Function AddServer adds a server and location blocks to existing site with the specefied domain name
func AddServer(directoryPath string, domain string, listenPort int) (string, error) {
    file, err := os.OpenFile(directoryPath + domain + ".conf", os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return "", fmt.Errorf("failed to open config file: %v", err)
    }
    defer file.Close()
    
    cfgVars := nginx.NewWebConfig()
    cfgVars.ConfigPath = directoryPath
    cfgVars.Domain = domain
    cfgVars.ListenPort = listenPort

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
        return fmt.Sprintf("creating web server config file is successful: %v", directoryPath + domain + ".conf"), nil
    }
    return "", fmt.Errorf("failed to validate web server config file: %v", err)
}

// Function GetEnabledSites return list of all the available sites
func GetAvailableSites(availableDirectoryPath string) ([]os.DirEntry, error) {
    sites, err := nginx.GetSites(availableDirectoryPath)
    if err != nil {
        return nil, err
    } 
    return sites, nil
}

// Function validateConfigServer checks for all the available information for web server and returns correct error message when missing
func validateConfigServer(cfg *nginx.WebConfig) (bool, error) {
    if cfg.ConfigPath == "" {
        return false, fmt.Errorf("config file path is not set")
    }
    if cfg.Domain == "" {
        return false, fmt.Errorf("must set a domain name")
    }
    if cfg.ListenPort == 0 {
        return false, fmt.Errorf("port number needs to be between 1-65535")
    }
    return true, nil
}