// High-level automation orchestration functions
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/IM-Malik/Gonix/nginx/sites/reverseproxy"
	// "github.com/IM-Malik/Gonix/nginx/modules"
	// "github.com/IM-Malik/Gonix/nginx/sites/webserver"
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
    defaults.ModulesEnabled =  modulesEnabledPath
}

func (defaults Defaults) GetDefaults() *Defaults {
    return &defaults
}

func SetAllDefaults(nginxConfPath string, sitesAvailablePath string, sitesEnabledPath string, modulesEnabledPath string) (*Defaults, error) {
    if(nginxConfPath != "" && sitesAvailablePath != "" && sitesEnabledPath != "" && modulesEnabledPath != ""){
        return &Defaults{
            NginxConf: nginxConfPath,
            SitesAvailable: sitesAvailablePath,
            SitesEnabled: sitesEnabledPath,
            ModulesEnabled: modulesEnabledPath,
        }, nil
    }
    return nil, fmt.Errorf("one of the parameters are not set.") //(you can set any of them with 'd' for default path)
	// Automator.config.BackupConfig = "/etc/nginx/backup-configs/"
}

func GetGlobalConfig(globalConfigFilePath string) (string, error) {
    content, err := ioutil.ReadFile(globalConfigFilePath)
    if err != nil {
        return "", fmt.Errorf("failed to read nginx.conf content: %v", err)
    }
    return "\n"+string(content), nil
}

func GetSiteConfig(siteFilePath string) (string, error) {
    content, err := ioutil.ReadFile(siteFilePath)
    if err != nil {
        return "", fmt.Errorf("failed to read site content: %v", err)
    }
    return "\n"+string(content), nil
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
        return "", fmt.Errorf("failed to test nginx: %v", err)
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

func CreateAndEnableRevProxy(defaults *Defaults, domain string, listenPort int, proxyPass string, uri string, enableSSL bool, certPath string, keyPath string) (string, error) {
    if(!enableSSL) {
        _, err := reverseproxy.AddSite(defaults.SitesAvailable, domain, listenPort, proxyPass, uri, enableSSL, certPath, keyPath)
        if err != nil {
            return "", fmt.Errorf("the site was not created: %v", err)
        }
        _, err = reverseproxy.EnableSite(defaults.SitesAvailable, defaults.SitesEnabled, domain)
        if err != nil {
            return "", fmt.Errorf("the site was not enabled: %v", err)
        }
    }
    return "the site was created and enabled successflly", nil
}

// call this from the UpdateSite function, so before any modifications the old version is saved for easy rollback (RollBackChanges) function
func BackupConfig(FilePath string) (string, error) {
    srcFile, err := os.Open(FilePath)
    if err != nil {
        return "", fmt.Errorf("failed to open file: %v", err)
    }
    defer srcFile.Close()

    dir := filepath.Dir(FilePath)
    base := filepath.Base(FilePath)
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

func RollbackChanges(originalFilePath string) (string, error) {
    dir := filepath.Dir(originalFilePath)
    base := filepath.Base(originalFilePath)
    backupPath := filepath.Join(dir, base+".bak")

    if _, err := os.Stat(backupPath); os.IsNotExist(err) {
        return "", fmt.Errorf("backup file %s is not found", backupPath)
    }
    backupFile, err := os.Open(backupPath)
    if err != nil {
        return "", fmt.Errorf("failed to open backup file: %v", err)
    }
    defer backupFile.Close()

    originalFile, err := os.Create(originalFilePath)
    if err != nil {
        return "", fmt.Errorf("failed to open original file for writing: %v", err)
    }
    defer originalFile.Close()

    if _, err := io.Copy(originalFile, backupFile); err != nil {
        return "", fmt.Errorf("failed to copy backup to original file: %v", err)
    }
    if err := originalFile.Sync(); err != nil {
        return "", fmt.Errorf("failed to sync original file: %v", err)
    }
    if err := os.Remove(backupPath); err != nil {
        return "", fmt.Errorf("failed to remove backup file: %v", err)
    }
    return fmt.Sprintf("rollback is successful at: %v", originalFilePath), nil
}

// Start Here ...
func UpdateSite() (string, error) {
    return "", nil
}

// Composite function that might validate, backup, apply changes, and then reload nginx, making it easier for users to perform all steps with a single call.
//Or CheckAll()
func RunFullCycle() (string, error) {
    return "", nil
}
