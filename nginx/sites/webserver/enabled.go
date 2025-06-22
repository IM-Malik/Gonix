package webserver

import (
	"os"

	"github.com/IM-Malik/Gonix/nginx"
)

// EnableSite enables the specific available site by domain name
func EnableSite(availableDirectoryPath string, enabledDirectoryPath string, domain string) (string, error) {
	return nginx.EnableSite(availableDirectoryPath, enabledDirectoryPath, domain)
}

// RemoveEnabledSite removes the enabled site, without removing the available site
func RemoveEnabledSite(enabledDirectoryPath string, domain string) (string, error) {
	return nginx.RemoveEnabledSite(enabledDirectoryPath, domain) 
}

// GetEnabledSites return list of all the enabled sites
func GetEnabledSites(enabledDirectoryPath string) ([]os.DirEntry, error) {
	sites, err := nginx.GetSites(enabledDirectoryPath)
	if err != nil {
		return nil, err
	}
    return sites, nil
}