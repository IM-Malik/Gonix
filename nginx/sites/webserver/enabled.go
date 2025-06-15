package webserver

import (
	"github.com/IM-Malik/Gonix/nginx"
)

func EnableSite(sourceDirectoryPath string, destDirectoryPath string, domainName string) (string, error) {
	return nginx.EnableSite(sourceDirectoryPath, destDirectoryPath, domainName)
}

func RemoveEnabledSite(enabledDirectoryPath string, domainName string) (string, error) {
	return nginx.RemoveEnabledSite(enabledDirectoryPath, domainName) 
}

func GetEnabledSites(enabledDirectoryPath string) (error) {
    return nginx.GetSites(enabledDirectoryPath)
}