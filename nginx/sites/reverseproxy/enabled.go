package reverseproxy

import (
	"os"
	"github.com/IM-Malik/Gonix/nginx"
)

// Function EnableSite enables the specific available site by domain name
func EnableSite(sourceDirectoryPath string, destDirectoryPath string, domain string) (string, error) {
	return nginx.EnableSite(sourceDirectoryPath, destDirectoryPath, domain)
}

// Function RemoveEnabledSite removes the enabled site, without removing the available site
func RemoveEnabledSite(enabledDirectoryPath string, domainName string) (string, error) {
	return nginx.RemoveEnabledSite(enabledDirectoryPath, domainName)
}

// Function GetEnabledSites return list of all the enabled sites
func GetEnabledSites(enabledDirectoryPath string) ([]os.DirEntry, error) {
	sites, err := nginx.GetSites(enabledDirectoryPath)
	if err != nil {
		return nil, err
	}
	return sites, nil
}