// Package nginx is responsible for having functions that are shared between more than one package
package nginx

import (
	"fmt"
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
