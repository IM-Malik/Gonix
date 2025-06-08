package webserver

import (
	"fmt"
	"os"
	"text/template"
	"github.com/IM-Malik/Gonix/nginx/sites"
    "github.com/IM-Malik/Gonix/nginx"
)

func AddSite(directoryPath string, domain string, listenPort int, uri string, staticContentPath string, staticContentFileName string) (string, error) {
	file, err := os.OpenFile(directoryPath + domain + ".conf", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create configuration file: %v", err)
	}
	defer file.Close()

    output, err := AddServer(directoryPath, domain,  listenPort)
    if err != nil {
        return "", fmt.Errorf("failed to add a site: %v", err)
    }
	return fmt.Sprintf("adding a site is successful: \n%v", output), nil
}

func RemoveSite(directoryPath string, domain string) (string, error) {
    return nginx.RemoveSite(directoryPath, domain)
//     err := os.Remove(directoryPath + domain + ".conf")
//     if err != nil {
//         return "", fmt.Errorf("failed to remove the config file: %v", err)
//     }
// 	return fmt.Sprintf("removal of config file " + directoryPath + domain + ".conf" + " is successful"), nil
}

// Advanced Feature
func UpdateSite() (string, error) {
    return "", nil
}

//---------------------------------------------------------------------------------------------------
// Sub-Functions

func AddServer(directoryPath string, domain string, listenPort int) (string, error) {
    file, err := os.OpenFile(directoryPath + domain + ".conf", os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return "", fmt.Errorf("failed to open config file: %v", err)
    }
    defer file.Close()

    cfgVars := sites.NewWebConfig()
    cfgVars.ConfigPath = directoryPath
    cfgVars.Domain = domain
    cfgVars.ListenPort = listenPort

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
        return fmt.Sprintf("creating web server config file is successful: %v", directoryPath + domain + ".conf"), nil
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