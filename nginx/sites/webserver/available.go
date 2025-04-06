package webserver

import (
	"fmt"
	"os"
	"text/template"
	"github.com/IM-Malik/Gonix/nginx/sites"
    "github.com/IM-Malik/Gonix/nginx"
)

func AddSite(directoryPath string, domainName string, portNumber int, urlPath string, staticContentPath string, staticContentFileName string) (string, error) {
	file, err := os.OpenFile(directoryPath + domainName + ".conf", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create configuration file: %v", err)
	}
	defer file.Close()

    output, err := AddServer(directoryPath, domainName,  portNumber)
    if err != nil {
        return "", fmt.Errorf("failed to add a site: %v", err)
    }
	return fmt.Sprintf("adding a site is successful: \n%v", output), nil
}

func RemoveSite(directoryPath string, domainName string) (string, error) {
    return nginx.RemoveSite(directoryPath, domainName)
//     err := os.Remove(directoryPath + domainName + ".conf")
//     if err != nil {
//         return "", fmt.Errorf("failed to remove the config file: %v", err)
//     }
// 	return fmt.Sprintf("removal of config file " + directoryPath + domainName + ".conf" + " is successful"), nil
}

// Advanced Feature
func UpdateSite() (string, error) {
    return "", nil
}

//---------------------------------------------------------------------------------------------------
// Sub-Functions

func AddServer(directoryPath string, domainName string, portNumber int) (string, error) {
    file, err := os.OpenFile(directoryPath + domainName + ".conf", os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return "", fmt.Errorf("failed to open config file: %v", err)
    }
    defer file.Close()

    cfgVars := sites.NewWebConfig()
    cfgVars.ConfigPath = directoryPath
    cfgVars.Domain = domainName
    cfgVars.ListenPort = portNumber

    status, err := validateConfigServer(cfgVars)
    if status {
        tmpl := template.Must(template.New("srvBlkTmpl").Parse(sites.SERVER_WEBSERVER_BLOCK_TMPL))
        if err := tmpl.Execute(file, cfgVars); err != nil {
            return "", fmt.Errorf("server template execution failed: %w", err)
        }
        tmpl2 := template.Must(template.New("locationBlkTmpl").Parse(sites.LOCATION_WEBSERVER_BLOCK_TMPL))
        if err := tmpl2.Execute(file, cfgVars); err != nil {
            return "", fmt.Errorf("location template execution failed: %w", err)
        }
        file.WriteString("}\n")
        return fmt.Sprintf("creating web server config file is successful: %v", directoryPath + domainName + ".conf"), nil
    }
    return "", fmt.Errorf("failed to validate web server config file: %v", err)
}

func validateConfigServer(cfg *sites.WebConfig) (bool, error) {
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