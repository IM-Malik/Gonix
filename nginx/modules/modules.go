// Package modules is responsible for operations on modules
package modules

import (
	"fmt"
	"os"
    "github.com/IM-Malik/Gonix/orch"
)

// Function EnableModule enables the installed modules for Nginx
func EnableModule(defaults *orch.Defaults, sourceDirectoryPath string, moduleName string) (string, error) {
    err := os.Symlink(sourceDirectoryPath + moduleName[3:], defaults.ModulesEnabled + moduleName)
    if err != nil {
        return "", fmt.Errorf("failed to enable the " + moduleName + " module: %v", err)
    }
    return "the module is enabled successfully", nil
}

// Function RemoveModule removes the enabled modules from Nginx
func RemoveModule(defaults *orch.Defaults, moduleName string) (string, error) {
    err := os.Remove(defaults.ModulesEnabled + moduleName)
    if err != nil {
        return "", fmt.Errorf("failed to disable the " + moduleName + " module: %v", err)
    }
    return "the module is removed successfully", nil
}

// Function GetEnabledModules returns all the modules in Nginx
func GetEnabledModules(defaults *orch.Defaults) ([]os.DirEntry, error) {
    enabledModules, err := os.ReadDir(defaults.ModulesEnabled)
    if err != nil {
        return nil, fmt.Errorf("failed to read the files inside the 'modules-enabled' directory: %v", err)
    }
    return enabledModules, nil
}