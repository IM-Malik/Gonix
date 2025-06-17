// Package webserver is used to create and manipulate web servers in Nginx
package webserver

import (
	"github.com/IM-Malik/Gonix/nginx"
)

// EnableSite is a function that gives the user the ability to enable an available site in Nginx
func EnableSite(sourceDirectoryPath string, destDirectoryPath string, domainName string) (string, error) {
	return nginx.EnableSite(sourceDirectoryPath, destDirectoryPath, domainName)
}

func RemoveEnabledSite(enabledDirectoryPath string, domainName string) (string, error) {
	return nginx.RemoveEnabledSite(enabledDirectoryPath, domainName) 
}

func GetEnabledSites(enabledDirectoryPath string) (error) {
    return nginx.GetSites(enabledDirectoryPath)
}