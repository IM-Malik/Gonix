package webserver

import (
	// 	"fmt"
	// 	"os"
	"github.com/IM-Malik/Gonix/nginx"
)

func EnableSite(sourceDirectoryPath string, destDirectoryPath string, domainName string) (string, error) {
	return nginx.EnableSite(sourceDirectoryPath, destDirectoryPath, domainName)
// 	err := os.Symlink(sourceDirectoryPath + domainName + ".conf", destDirectoryPath + domainName + ".conf")
// 	if err != nil {
// 		return "", fmt.Errorf("failed to enable the site: %v", err)
// 	}
// 	return fmt.Sprintf("enabling the site is successful at: %v", destDirectoryPath + domainName + ".conf"), nil
}

func RemoveEnabledSite(enabledDirectoryPath string, domainName string) (string, error) {
	return nginx.RemoveEnabledSite(enabledDirectoryPath, domainName) 
	// 	err := os.Remove(enabledDirectoryPath + domainName + ".conf")
//     if err != nil {
//         return "", fmt.Errorf("failed to remove the enabled config file: %v", err)
//     }
// 	return fmt.Sprintf("removal of enabled config file " + enabledDirectoryPath + domainName + ".conf" + "is successful"), nil
}

func GetEnabledSites(enabledDirectoryPath string) (error) {
    return nginx.GetSites(enabledDirectoryPath)
}