// Package nginx is responsible for having functions that are shared between more than one package
package nginx

import (
	"fmt"
	"html/template"
	"os"
)

// RemoveSite removes reverseproxy/webserver available/enabled sites, based on the directory path
func RemoveSite(directoryPath string, domain string) (string, error) {
    err := os.Remove(directoryPath + domain + ".conf")
    if err != nil {
        return "", fmt.Errorf("could not remove the configuration file: %v", err)
    }
	return fmt.Sprintf("removal of config file " + directoryPath + domain + ".conf" + " is successful"), nil
}

// EnableSite enables reverseproxy/webserver available/enabled sites, based on the directory path
func EnableSite(sourceDirectoryPath string, destDirectoryPath string, domain string) (string, error) {
	err := os.Symlink(sourceDirectoryPath + domain + ".conf", destDirectoryPath + domain + ".conf")
	if err != nil {
		return "", fmt.Errorf("failed to enable the site: %v", err)
	}
	return fmt.Sprintf("enabling the site is successful at: %v", destDirectoryPath + domain + ".conf"), nil
}

// RemoveEnabledSite removes the enabled reverseproxy/webserver enabled sites specified by domain name 
func RemoveEnabledSite(enabledDirectoryPath string, domain string) (string, error) {
	err := os.Remove(enabledDirectoryPath + domain + ".conf")
    if err != nil {
        return "", fmt.Errorf("could not remove the enabled configuration file: %v", err)
    }
	return fmt.Sprintf("removal of enabled config file " + enabledDirectoryPath + domain + ".conf" + "is successful"), nil
}

// GetSites returns a slice with all files in a directory in reverseproxy/webserver available/enabled sites, based on the directory path
func GetSites(directoryPath string) ([]os.DirEntry, error) {
    sites, err := os.ReadDir(directoryPath)
    if err != nil {
        return nil, fmt.Errorf("could not read files in directory %s: %v", directoryPath, err)
    }
    return sites, nil
}

// AddUpstream adds an upstream block to the reverseproxy/webserver available/enabled sites, based on the directory path
func AddUpstream(directoryPath string, domain string, upstreamName string, serverIP string, portNumber int) (string, error) {
    file, err := os.OpenFile(directoryPath + domain + ".conf", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

    cfgVars := NewUpstream()
    cfgVars.ConfigPath = directoryPath
    cfgVars.Name = upstreamName
    cfgVars.ServerIP = serverIP
    cfgVars.PortNumber = portNumber

    status, err := validateConfigUpstream(cfgVars)
    if status {
        tmpl := template.Must(template.New("upstreamBlkTmpl").Parse(UPSTREAM_BLOCK_TMPL))
        if err := tmpl.Execute(file, cfgVars); err != nil {
            return "", fmt.Errorf("upstream template execution failed: %w", err)
        }
        return fmt.Sprintf("upstream block is added succesfully in: %v", directoryPath + domain + ".conf"), nil
    }
    return "", fmt.Errorf("configuration validation failed: %v", err)
}

// validateConfigUpstream checks for all the available information for upstream and returns correct error message when missing
func validateConfigUpstream(cfg *Upstream) (bool, error) {
    if cfg.ConfigPath == "" {
        return false, fmt.Errorf("config file path is not set")
    }
    if cfg.Name == ""{
        return false, fmt.Errorf("must specify an upstream name")
    }
    if cfg.ServerIP == "" {
        return false, fmt.Errorf("must specify a server IP")
    }
    if cfg.PortNumber == 0 {
        return false, fmt.Errorf("port number needs to be between 1-65535")
    }
    return true, nil
}