package nginx

import (
	"fmt"
	"html/template"
	"os"

	"github.com/IM-Malik/Gonix/nginx/sites"
)

func RemoveSite(directoryPath string, domainName string) (string, error) {
    err := os.Remove(directoryPath + domainName + ".conf")
    if err != nil {
        return "", fmt.Errorf("failed to remove the config file: %v", err)
    }
	return fmt.Sprintf("removal of config file " + directoryPath + domainName + ".conf" + " is successful"), nil
}

func EnableSite(sourceDirectoryPath string, destDirectoryPath string, domainName string) (string, error) {
	err := os.Symlink(sourceDirectoryPath + domainName + ".conf", destDirectoryPath + domainName + ".conf")
	if err != nil {
		return "", fmt.Errorf("failed to enable the site: %v", err)
	}
	return fmt.Sprintf("enabling the site is successful at: %v", destDirectoryPath + domainName + ".conf"), nil
}

func RemoveEnabledSite(enabledDirectoryPath string, domainName string) (string, error) {
	err := os.Remove(enabledDirectoryPath + domainName + ".conf")
    if err != nil {
        return "", fmt.Errorf("failed to remove the enabled config file: %v", err)
    }
	return fmt.Sprintf("removal of enabled config file " + enabledDirectoryPath + domainName + ".conf" + "is successful"), nil
}

// enabled or available based on the directory path
func GetSites(directoryPath string) (error) {
    sites, err := os.ReadDir(directoryPath)
    if err != nil {
        return fmt.Errorf("failed to read the files inside the 'modules-enabled' directory: %v", err)
    }
    for i := range sites {
        fmt.Println(sites[i])
    }
    return nil
}

func AddUpstream(directoryPath string, domainName string, upstreamName string, serverIP string, portNumber int) (string, error) {
    file, err := os.OpenFile(directoryPath + domainName + ".conf", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

    cfgVars := sites.NewUpstream()
    cfgVars.ConfigPath = directoryPath
    cfgVars.Name = domainName
    cfgVars.ServerIP = serverIP
    cfgVars.PortNumber = portNumber

    status, err := validateConfigUpstream(cfgVars)
    if status {
        tmpl := template.Must(template.New("upstreamBlkTmpl").Parse(sites.UPSTREAM_BLOCK_TMPL))
        if err := tmpl.Execute(file, cfgVars); err != nil {
            return "", fmt.Errorf("upstream template execution failed: %w", err)
        }
        return fmt.Sprintf("upstream block is added succesfully in: %v", directoryPath + domainName + ".conf"), nil
    }
    return "", fmt.Errorf("failed to validate config file: %v", err)
}

func validateConfigUpstream(cfg *sites.Upstream) (bool, error) {
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