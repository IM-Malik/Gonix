// Adding, deleting, and modifying site files
package reverseproxy

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	// "github.com/IM-Malik/Gonix/nginx/sites"
	"github.com/IM-Malik/Gonix/nginx"
)

// Main Functions
// Or AddConfigFile
func AddSite(directoryPath string, domain string, listenPort int, proxyPass string, uri string, enableSSL bool, certPath string, keyPath string) (string, error) {
    file, err := os.OpenFile(directoryPath + domain + ".conf", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
        RemoveSite(directoryPath, domain)
		return "", fmt.Errorf("failed to create configuration file: %v", err)
	}
	defer file.Close()

    output, err := AddServer(directoryPath, domain,  listenPort, proxyPass, uri, enableSSL, certPath, keyPath)
    if err != nil {
        RemoveSite(directoryPath, domain)
        return "", fmt.Errorf("failed to add a site: %v", err)
    }
	return fmt.Sprintf("adding a site is successful: \n%v", output), nil
}
// Or RemoveConfigFile
func RemoveSite(directoryPath string, domain string) (string, error) {   
    return nginx.RemoveSite(directoryPath, domain)
}

// Advanced Feature
// Or UpdateConfigFile
// Finished with mail and stream. Changed stuff in both 'available'. Start Here Next...
func UpdateSite() (string) {
    return ""
}
// REMOVE AND REPLACE BLOCK OR JUST MODIFY EXISTING BLOCKS?
// Advanced Feature

//---------------------------------------------------------------------------------------------------
// Sub-Functions

func AddServer(directoryPath string, domain string, listenPort int, proxyPass string, uri string, enableSSL bool, certPath string, keyPath string) (string, error) {
    file, err := os.OpenFile(directoryPath + domain + ".conf", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()
    cfgVars := nginx.NewRevConfig()
    cfgVars.ConfigPath = directoryPath
    cfgVars.Domain = domain
    cfgVars.ListenPort = listenPort
    // cfgVars.RootDir = "/var/www/html"
    // cfgVars.IndexFiles = "index.html"
    cfgVars.ProxyPass = proxyPass
    cfgVars.EnableSSL = enableSSL
    cfgVars.SSLCertPath = certPath
    cfgVars.SSLKeyPath = keyPath
    cfgVars.URI = uri

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

func AddLocation(directoryPath string, domain string, proxyPass string, uri string) (string, error) {
    file, err := os.OpenFile(directoryPath + domain + ".conf", os.O_RDWR, 0644)
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

    cfgVars := nginx.NewRevConfig()
    cfgVars.ConfigPath = directoryPath
    cfgVars.ProxyPass = proxyPass
    cfgVars.URI = uri

    status, err := validateConfigLocation(cfgVars)
    if status {
        tmpl := template.Must(template.New("locationBlkTmpl").Parse(nginx.LOCATION_REVERSEPROXY_BLOCK_TMPL))
        if err := tmpl.Execute(file, cfgVars); err != nil {
            return "", fmt.Errorf("server template execution failed: %w", err)
        }
        file.WriteString("}\n")
        return fmt.Sprintf("location block is added successfully in: %v", directoryPath + domain + ".conf"), nil
    }
    return "", fmt.Errorf("failed to validate config file: %v", err)
}

func GetAvailableSites(availableDirectoryPath string) (error) {
    return nginx.GetSites(availableDirectoryPath)
}

// func AddUpstream(directoryPath string, domain string, upstreamName string, serverIP string, listenPort int) (string, error) {
//     file, err := os.OpenFile(directoryPath + domain + ".conf", os.O_APPEND|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to open config file: %v", err)
// 	}
// 	defer file.Close()

//     cfgVars := sites.NewUpstream()
//     cfgVars.ConfigPath = directoryPath
//     cfgVars.Name = domain
//     cfgVars.ServerIP = serverIP
//     cfgVars.listenPort = listenPort

//     status, err := validateConfigUpstream(cfgVars)
//     if status {
//         tmpl := template.Must(template.New("upstreamBlkTmpl").Parse(sites.UPSTREAM_BLOCK_TMPL))
//         if err := tmpl.Execute(file, cfgVars); err != nil {
//             return "", fmt.Errorf("upstream template execution failed: %w", err)
//         }
//         return fmt.Sprintf("upstream block is added succesfully in: %v", directoryPath + domain + ".conf"), nil
//     }
//     return "", fmt.Errorf("failed to validate config file: %v", err)
// }

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

func validateConfigLocation(cfg *nginx.RevConfig) (bool, error) {
    if cfg.ConfigPath == "" {
        return false, fmt.Errorf("config file path is not set")
    }
    // if cfg.ProxyPass == "" && cfg.RootDir == "" {
    //     return false, fmt.Errorf("must specify either ProxyPass or RootDir")
    // }
    if cfg.URI != "" && cfg.ProxyPass == "" {
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
//     if cfg.listenPort == 0 {
//         return false, fmt.Errorf("port number needs to be between 1-65535")
//     }
//     return true, nil
// }