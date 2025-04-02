// Creating/removing symbolic links
package reverseproxy

import (
	"fmt"
	"os"
)

func EnableSite(sourceDirectoryPath string, destDirectoryPath string, domainName string) (string, error) {
	err := os.Symlink(sourceDirectoryPath + domainName + ".conf", destDirectoryPath + domainName + ".conf")
	if err != nil {
		return "", fmt.Errorf("failed to enable the site: %v", err)
	}
	return fmt.Sprintf("enabling the site is successful at: %v", destDirectoryPath + domainName + ".conf"), nil
}

func RemoveEnabledSite(enabledDirectoryPath string, domainName string) (string, error) {
	err := os.Remove(enabledDirectoryPath + domainName + ".conf")
    if err != nil {
        return "", fmt.Errorf("failed to remove the enabled config file: %v", err)
    }
	return fmt.Sprintf("removal of enabled config file " + enabledDirectoryPath + domainName + ".conf" + "is successful"), nil
}