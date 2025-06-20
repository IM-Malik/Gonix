// Package reverseproxy is responsible for Adding, deleting, modifying, creating/removing symbolic links in reverse proxy site files
package reverseproxy

import (
	"fmt"
	"os"
	"text/template"
	"github.com/IM-Malik/Gonix/nginx"
)

// Function AddSite adds a complete reverse proxy site, no need for extra function calling.
func AddSite(directoryPath string, domain string, listenPort int, proxyPass string, uri string, enableSSL bool, certPath string, keyPath string, httpOrhttps string) (string, error) {
	file, err := os.OpenFile(directoryPath + domain + ".conf", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		RemoveSite(directoryPath, domain)
		return "", fmt.Errorf("failed to create configuration file: %v", err)
	}
	defer file.Close()

	output, err := AddServer(directoryPath, domain, listenPort, proxyPass, uri, enableSSL, certPath, keyPath, httpOrhttps)
	if err != nil {
		RemoveSite(directoryPath, domain)
		return "", fmt.Errorf("failed to add a site: %v", err)
	}
	return fmt.Sprintf("adding a site is successful: \n%v", output), nil
}

// Function RemoveSite removes any existing site with the specefied domain name
func RemoveSite(directoryPath, domain string) (string, error) {
	return nginx.RemoveSite(directoryPath, domain)
}

// Function AddServer adds a server and location blocks to existing site with the specefied domain name
func AddServer(directoryPath string, domain string, listenPort int, proxyPass string, uri string, enableSSL bool, certPath string, keyPath string, httpOrhttps string) (string, error) {
	file, err := os.OpenFile(directoryPath+domain+".conf", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()
	cfgVars := nginx.NewRevConfig()
	cfgVars.ConfigPath = directoryPath
	cfgVars.Domain = domain
	cfgVars.ListenPort = listenPort
	cfgVars.ProxyPass = proxyPass
	cfgVars.EnableSSL = enableSSL
	cfgVars.SSLCertPath = certPath
	cfgVars.SSLKeyPath = keyPath
	cfgVars.URI = uri
    cfgVars.HttpOrHttps = httpOrhttps

	status, err := validateConfigServer(cfgVars)
	if status {
		tmpl := template.Must(template.New("srvBlkTmpl").Parse(nginx.SERVER_REVERSEPROXY_BLOCK_TMPL))
		if err := tmpl.Execute(file, cfgVars); err != nil {
		    return "", fmt.Errorf("server template execution failed: %w", err)
		}
		tmpl2 := template.Must(template.New("locationBlkTmpl").Parse(nginx.LOCATION_REVERSEPROXY_BLOCK_TMPL))
		if err := tmpl2.Execute(file, cfgVars); err != nil {
		    return "", fmt.Errorf("location template execution failed: %w", err)
		}
		file.WriteString("}\n")
		return fmt.Sprintf("creating config file with SSL is successful: %v", directoryPath + domain + ".conf"), nil
	}
	return "", fmt.Errorf("failed to validate config file: %v", err)
}

// Function GetEnabledSites return list of all the available sites
func GetAvailableSites(availableDirectoryPath string) ([]os.DirEntry, error) {
	sites, err := nginx.GetSites(availableDirectoryPath)
	if err != nil {
		return nil, err
	}
	return sites, nil
}

// Function validateConfigServer checks for all the available information for reverse proxy and returns correct error message when missing
func validateConfigServer(cfg *nginx.RevConfig) (bool, error) {
	if cfg.ConfigPath == "" {
		return false, fmt.Errorf("config file path is not set")
	}
	if cfg.Domain == "" {
		return false, fmt.Errorf("domain name is required")
	}
	if cfg.ListenPort <= 0 || cfg.ListenPort > 65535 {
		return false, fmt.Errorf("port number needs to be between 1-65535")
	}
	if cfg.ProxyPass == "" /*&& cfg.RootDir == ""*/ {
		return false, fmt.Errorf("must specify either ProxyPass or RootDir")
	}
	if cfg.EnableSSL && cfg.ListenPort == 80 {
		return false, fmt.Errorf("cannot use port number 80 with SSL enabled. (use default 443 or other port number)")
	}
	if !cfg.EnableSSL && cfg.ListenPort == 443 {
		return false, fmt.Errorf("cannot use port number 443 without SSL enabled. (change the EnableSSL to true)")
	}
	if cfg.URI != "" && cfg.ProxyPass == "" {
		return false, fmt.Errorf("must specify ProxyPass")
	}
	return true, nil
}