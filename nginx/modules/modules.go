// Package modules is responsible for operations on modules
package modules

import (
	"fmt"
	"os"
    "github.com/IM-Malik/Gonix/orch"
)

// STOPPED HERE IN THE DOCUMENTATION PROCESS. finish and go down in the file structure
func EnableModule(defaults *orch.Defaults, sourceDirectoryPath string, moduleName string) (string, error) {
    err := os.Symlink(sourceDirectoryPath + moduleName[3:], defaults.ModulesEnabled + moduleName)
    if err != nil {
        return "", fmt.Errorf("failed to enable the " + moduleName + " module: %v", err)
    }
    return "the module is enabled successfully", nil
}

func RemoveModule(defaults *orch.Defaults, moduleName string) (string, error) {
    err := os.Remove(defaults.ModulesEnabled + moduleName)
    if err != nil {
        return "", fmt.Errorf("failed to disable the " + moduleName + " module: %v", err)
    }
    return "the module is removed successfully", nil
}

func GetEnabledModules(defaults *orch.Defaults) (error) {
    enabledModules, err := os.ReadDir(defaults.ModulesEnabled)
    if err != nil {
        return fmt.Errorf("failed to read the files inside the 'modules-enabled' directory: %v", err)
    }
    for i := range enabledModules {
        fmt.Println(enabledModules[i])
    }
    return nil
}