// High-level automation orchestration functions
package nginx

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type AutomatorConfig struct {
	NginxConf      string // Path to the nginx.conf file
	SitesAvailable string // Path to the sites-available directory
	SitesEnabled   string // Path to the sites-enabled directory
	ModulesEnabled string // Path to the modules-enabled directory
	BackupConfig   string // Path to the backup of modification on configurations
}

type Automator struct {
	config AutomatorConfig
	//logger *logger.mylogger
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

// Start Here...
func CreateAndEnableNewConfig() (string, error) {
    return "", nil
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

// // Advanced Feature
// func UpdateSite() (string, error) {
//     return "", nil
// }
// // Advanced Feature

// Composite function that might validate, backup, apply changes, and then reload nginx, making it easier for users to perform all steps with a single call.
//Or CheckAll()
func RunFullCycle() (string, error) {
    return "", nil
}

func NewAutomator() (string, error) {
    return "", nil
}

func DefaultAutomator() (string, error) {
	// Automator.config.NginxConf = "/etc/nginx/nginx.conf"
	// Automator.config.SitesAvailable = "/etc/nginx/sites-available/"
	// Automator.config.SitesEnabled = "/etc/nginx/sites-enabled/"
	// Automator.config.ModulesEnabled = "/etc/nginx/modules-enabled/"
	// Automator.config.BackupConfig = "/etc/nginx/backup-configs/"
    return "", nil
}



/* An example of how to use the structs
// NewAutomator initializes an Automator with the provided configuration.
func NewAutomator(cfg AutomatorConfig, logger *log.Logger) *Automator {
    // Set defaults if necessary
    if cfg.NginxConf == "" {
        cfg.NginxConf = "/etc/nginx/nginx.conf"
    }
    if cfg.SitesAvailable == "" {
        cfg.SitesAvailable = "/etc/nginx/sites-available"
    }
    if cfg.SitesEnabled == "" {
        cfg.SitesEnabled = "/etc/nginx/sites-enabled"
    }
    if cfg.ModulesEnabled == "" {
        cfg.ModulesEnabled = "/etc/nginx/modules-enabled"
    }
    return &Automator{config: cfg, logger: logger}
}
*/
