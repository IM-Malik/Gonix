// Operations for module management (/modules-enabled)
package modules

import (
	"fmt"
	"os"
)

func EnableModule(sourceDirectoryPath string, destDirectoryPath string, moduleName string) (string, error) {
    err := os.Symlink(sourceDirectoryPath + moduleName[3:], destDirectoryPath + moduleName)
    if err != nil {
        return "", fmt.Errorf("failed to enable the " + moduleName + " module: %v", err)
    }
    return "the module is enabled successfully", nil
}

// Make defualt option for directoryPath in all functions
func DisableModule(directoryPath string, moduleName string) (string, error) {
    err := os.Remove(directoryPath + moduleName)
    if err != nil {
        return "", fmt.Errorf("failed to disable the " + moduleName + " module: %v", err)
    }
    return "the module is disabled successfully", nil
}

func GetEnabledModules(directoryPath string) (error) {
    enabledModules, err := os.ReadDir(directoryPath)
    if err != nil {
        return fmt.Errorf("failed to read the files inside the 'modules-enabled' directory: %v", err)
    }
    for i := range enabledModules {
        fmt.Println(enabledModules[i])
    }
    return nil
}