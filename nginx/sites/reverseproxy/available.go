// Adding, deleting, and modifying site files
package reverseproxy

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"github.com/IM-Malik/Gonix/nginx/sites"
	"github.com/IM-Malik/Gonix/nginx"
)

// Main Functions
// Or AddConfigFile
func AddSite(directoryPath string, domainName string, portNumber int, proxyPass string, urlPath string, certPath string, keyPath string, enableSSL bool) (string, error) {
    file, err := os.OpenFile(directoryPath + domainName + ".conf", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create configuration file: %v", err)
	}
	defer file.Close()

    output, err := AddServer(directoryPath, domainName,  portNumber, proxyPass, certPath, keyPath, urlPath, enableSSL)
    if err != nil {
        return "", fmt.Errorf("failed to add a site: %v", err)
    }
	return fmt.Sprintf("adding a site is successful: \n%v", output), nil
}
// Or RemoveConfigFile
func RemoveSite(directoryPath string, domainName string) (string, error) {   
    return nginx.RemoveSite(directoryPath, domainName)
//     err := os.Remove(directoryPath + domainName + ".conf")
//     if err != nil {
//         return "", fmt.Errorf("failed to remove the config file: %v", err)
//     }
// 	return fmt.Sprintf("removal of config file " + directoryPath + domainName + ".conf" + " is successful"), nil
}

// Advanced Feature
// Or UpdateConfigFile
func UpdateSite() (string, error) {
    return "", nil
}
// REMOVE AND REPLACE BLOCK OR JUST MODIFY EXISTING BLOCKS?
// Advanced Feature

//---------------------------------------------------------------------------------------------------
// Sub-Functions

func AddServer(directoryPath string, domainName string, portNumber int, proxyPass string, certPath string, keyPath string, urlPath string, enableSSL bool) (string, error) {
    file, err := os.OpenFile(directoryPath + domainName + ".conf", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()
    cfgVars := sites.NewRevConfig()
    cfgVars.ConfigPath = directoryPath
    cfgVars.Domain = domainName
    cfgVars.ListenPort = portNumber
    // cfgVars.RootDir = "/var/www/html"
    // cfgVars.IndexFiles = "index.html"
    cfgVars.ProxyPass = proxyPass
    cfgVars.EnableSSL = enableSSL
    cfgVars.SSLCertPath = certPath
    cfgVars.SSLKeyPath = keyPath
    cfgVars.URLPath = urlPath

    status, err := validateConfigServer(cfgVars)
    if status {
        tmpl := template.Must(template.New("srvBlkTmpl").Parse(sites.SERVER_REVERSEPROXY_BLOCK_TMPL))
        if err := tmpl.Execute(file, cfgVars); err != nil {
            return "", fmt.Errorf("server template execution failed: %w", err)
        }
        tmpl2 := template.Must(template.New("locationBlkTmpl").Parse(sites.LOCATION_REVERSEPROXY_BLOCK_TMPL))
        if err := tmpl2.Execute(file, cfgVars); err != nil {
            return "", fmt.Errorf("location template execution failed: %w", err)
        }
        file.WriteString("}\n")
        return fmt.Sprintf("creating config file with SSL is successful: %v", directoryPath + domainName + ".conf"), nil
    }
    return "", fmt.Errorf("failed to validate config file: %v", err)
}

func AddLocation(directoryPath string, domainName string, proxyPass string, urlPath string) (string, error) {
    file, err := os.OpenFile(directoryPath + domainName + ".conf", os.O_RDWR, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "false1", err
	}

	size := stat.Size()
	buffer := make([]byte, size)

	_, err = file.Read(buffer)
	if err != nil {
		return "false2", err
	}

	content := string(buffer)
	lastIndex := strings.LastIndex(content, "}")
	if lastIndex == -1 {
		return "false3", err
	}

	insertPos := int64(lastIndex)
	_, err = file.Seek(insertPos, os.SEEK_SET)
	if err != nil {
		return "false4", err
	}

    cfgVars := sites.NewRevConfig()
    cfgVars.ConfigPath = directoryPath
    cfgVars.ProxyPass = proxyPass
    cfgVars.URLPath = urlPath

    status, err := validateConfigLocation(cfgVars)
    if status {
        tmpl := template.Must(template.New("locationBlkTmpl").Parse(sites.LOCATION_REVERSEPROXY_BLOCK_TMPL))
        if err := tmpl.Execute(file, cfgVars); err != nil {
            return "", fmt.Errorf("server template execution failed: %w", err)
        }
        file.WriteString("}\n")
        return fmt.Sprintf("location block is added successfully in: %v", directoryPath + domainName + ".conf"), nil
    }
    return "", fmt.Errorf("failed to validate config file: %v", err)
}

// func AddUpstream(directoryPath string, domainName string, upstreamName string, serverIP string, portNumber int) (string, error) {
//     file, err := os.OpenFile(directoryPath + domainName + ".conf", os.O_APPEND|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to open config file: %v", err)
// 	}
// 	defer file.Close()

//     cfgVars := sites.NewUpstream()
//     cfgVars.ConfigPath = directoryPath
//     cfgVars.Name = domainName
//     cfgVars.ServerIP = serverIP
//     cfgVars.PortNumber = portNumber

//     status, err := validateConfigUpstream(cfgVars)
//     if status {
//         tmpl := template.Must(template.New("upstreamBlkTmpl").Parse(sites.UPSTREAM_BLOCK_TMPL))
//         if err := tmpl.Execute(file, cfgVars); err != nil {
//             return "", fmt.Errorf("upstream template execution failed: %w", err)
//         }
//         return fmt.Sprintf("upstream block is added succesfully in: %v", directoryPath + domainName + ".conf"), nil
//     }
//     return "", fmt.Errorf("failed to validate config file: %v", err)
// }

func validateConfigServer(cfg *sites.RevConfig) (bool, error) {
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
    if cfg.URLPath != "" && cfg.ProxyPass == "" {
        return false, fmt.Errorf("must specify ProxyPass")
    }
    return true, nil
}

func validateConfigLocation(cfg *sites.RevConfig) (bool, error) {
    if cfg.ConfigPath == "" {
        return false, fmt.Errorf("config file path is not set")
    }
    // if cfg.ProxyPass == "" && cfg.RootDir == "" {
    //     return false, fmt.Errorf("must specify either ProxyPass or RootDir")
    // }
    if cfg.URLPath != "" && cfg.ProxyPass == "" {
        return false, fmt.Errorf("must specify ProxyPass")
    }
    return true, nil
}

// func validateConfigUpstream(cfg *sites.Upstream) (bool, error) {
//     if cfg.ConfigPath == "" {
//         return false, fmt.Errorf("config file path is not set")
//     }
//     if cfg.Name == ""{
//         return false, fmt.Errorf("must specify an upstream name")
//     }
//     if cfg.ServerIP == "" {
//         return false, fmt.Errorf("must specify a server IP")
//     }
//     if cfg.PortNumber == 0 {
//         return false, fmt.Errorf("port number needs to be between 1-65535")
//     }
//     return true, nil
// }