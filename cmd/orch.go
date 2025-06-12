// High-level automation orchestration functions
package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	// "io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/IM-Malik/Gonix/nginx"
	"github.com/IM-Malik/Gonix/nginx/sites/reverseproxy"
	// "github.com/IM-Malik/Gonix/nginx/modules"
	// "github.com/IM-Malik/Gonix/nginx/sites"
	// "github.com/IM-Malik/Gonix/nginx/config"
	// "github.com/IM-Malik/Gonix/nginx"
)

type Defaults struct {
	NginxConf      string // Path to the nginx.conf file
	SitesAvailable string // Path to the sites-available directory
	SitesEnabled   string // Path to the sites-enabled directory
	ModulesEnabled string // Path to the modules-enabled directory
	// BackupConfig   string // Path to the backup of modification on configurations
}

func (defaults *Defaults) SetNgnixConf(nginxConfPath string) {
	defaults.NginxConf = nginxConfPath
}

func (defaults *Defaults) SetSitesAvailable(sitesAvailablePath string) {
	defaults.SitesAvailable = sitesAvailablePath
}

func (defaults *Defaults) SetSitesEnabled(sitesEnabledPath string) {
	defaults.SitesEnabled = sitesEnabledPath
}

func (defaults *Defaults) SetModulesEnabled(modulesEnabledPath string) {
	defaults.ModulesEnabled = modulesEnabledPath
}

func (defaults Defaults) GetDefaults() *Defaults {
	return &defaults
}

func SetAllDefaults(nginxConfPath string, sitesAvailablePath string, sitesEnabledPath string, modulesEnabledPath string) (*Defaults, error) {
	if nginxConfPath != "" && sitesAvailablePath != "" && sitesEnabledPath != "" && modulesEnabledPath != "" {
		return &Defaults{
			NginxConf:      nginxConfPath,
			SitesAvailable: sitesAvailablePath,
			SitesEnabled:   sitesEnabledPath,
			ModulesEnabled: modulesEnabledPath,
		}, nil
	}
	return nil, fmt.Errorf("one or more of the parameters are not set.") //(you can set any of them with 'd' for default path)
	// Automator.config.BackupConfig = "/etc/nginx/backup-configs/"
}

func GetGlobalConfig(globalConfigFilePath string) (string, error) {
	content, err := os.ReadFile(globalConfigFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read nginx.conf content: %v", err)
	}
	return "\n" + string(content), nil
}

func GetSiteConfig(siteFilePath string) (string, error) {
	content, err := os.ReadFile(siteFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read site content: %v", err)
	}
	return "\n" + string(content), nil
}

func ReloadNginx() (string, error) {
	cmd := exec.Command("nginx", "-s", "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to reload nginx: %v", err)
	}
	return fmt.Sprintf("nginx is reloaded successfully: %v", string(output)), nil
}

func TestNginx() (string, error) {
	cmd := exec.Command("nginx", "-t")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to reload nginx: %v\nnginx output:\n%s", err, string(output))
	}
	return fmt.Sprintf("nginx is testing successfully: %v", string(output)), nil
}

func RestartNginx() (string, error) {
	cmd := exec.Command("systemctl", "restart", "nginx")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to restart nginx process: %v", err)
	}
	return "nginx process is restarted successfully", nil
}

func CreateAndEnableRevProxy(defaults *Defaults, domain string, listenPort int, uri string, enableSSL bool, certPath string, keyPath string, upstreamName string, serverIP string, portNum int, httpOrHttps string) (string, error) {
	if httpOrHttps == "http" || httpOrHttps == "https" {
		_, err := reverseproxy.AddSite(defaults.SitesAvailable, domain, listenPort, upstreamName, uri, enableSSL, certPath, keyPath, httpOrHttps)
		_, err1 := reverseproxy.EnableSite(defaults.SitesAvailable, defaults.SitesEnabled, domain)
		_, err2 := nginx.AddUpstream(defaults.SitesAvailable, domain, upstreamName, serverIP, portNum)
		if err != nil || err1 != nil || err2 != nil {
			return "", fmt.Errorf("the reverse proxy site was not created: %v", errors.Join(err, err1, err2))
		}
	}
	return "the site was created and enabled successflly", nil
}

func RemoveSite(defaults *Defaults, domain string) (string, error) {
	_, err := nginx.RemoveSite(defaults.SitesAvailable, domain)
	_, err1 := nginx.RemoveEnabledSite(defaults.SitesEnabled, domain)
	if err != nil || err1 != nil {
		return "", fmt.Errorf("the site was not removed: %v", errors.Join(err, err1))
	}
	return "the site was removed successfully", nil
}

// call this from the UpdateSite function, so before any modifications the old version is saved for easy rollback (RollBackChanges) function
func BackupConfig(defaults *Defaults, domain string) (string, error) {
	srcFile, err := os.Open(defaults.SitesAvailable + domain + ".conf")
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer srcFile.Close()

	dir := filepath.Dir(defaults.SitesAvailable)
	base := filepath.Base(domain + ".conf")
	backupPath := filepath.Join(dir, base+".bak")

	dstFile, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %v", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return "", fmt.Errorf("failed to copy data: %v", err)
	}

	if err := dstFile.Sync(); err != nil {
		return "", fmt.Errorf("failed to sync backup file: %v", err)
	}

	return fmt.Sprintf("backup is created successfully at: %s", backupPath), nil
}

func RollbackChanges(defaults *Defaults, domain string) (string, error) {
	dir := filepath.Dir(defaults.SitesAvailable)
	base := filepath.Base(domain + ".conf")
	backupFile := filepath.Join(dir, base + ".bak")
    oldFile := filepath.Join(dir, base)

	if _, err := os.Stat(backupFile); err != nil {
        return "", fmt.Errorf("backup file %s is not found", backupFile)
	}
    if err := os.Remove(oldFile); err != nil {
        return "", fmt.Errorf("failed to remove backup file: %v", err)
    }
    if err := os.Rename(backupFile, oldFile); err != nil {
        return "", fmt.Errorf("failed to rollback changes: %v", err)
    }

    return fmt.Sprintf("rollback is successful at: %v", oldFile), nil
}

func UpdateSite(defaults *Defaults, domain string, oldText string, newText string) (string, error) {
    _, err := BackupConfig(defaults, domain)
    if err != nil  {
        return "", fmt.Errorf("could not create a backup site file: %v", err)
    }
    data, err := os.ReadFile(defaults.SitesAvailable + domain + ".conf")
    if err != nil {
        return "", fmt.Errorf("the site file does not exist: %v", err)
    }
    if strings.Count(string(data), oldText) == 0 {
        return "", fmt.Errorf("the old text was not found in the site file")
    }
    updated := strings.ReplaceAll(string(data), oldText, newText)
    err = os.WriteFile(defaults.SitesAvailable + domain + ".conf", []byte(updated), 0)    
    if err != nil {
        RollbackChanges(defaults, domain)
        return "", fmt.Errorf("the site file was not written. rolling back changes: %v", err)
    }
    // ReloadNginx()
	return "the site file was updated successfully", nil
}

// Composite function that might validate, backup, apply changes, and then reload nginx, making it easier for users to perform all steps with a single call.
// Or CheckAll()
func RunFullCycle() (string, error) {
	return "", nil
}
