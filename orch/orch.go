// Package orch is a high-level package provide automation functions
//
// NOTE: Many operations in this library require elevated permissions to modify system files.
// When running code that uses this library, you should use `sudo go run ...` or run your binary as root
// to ensure all file and service operations succeed. If you build your program first, remember to run the
// built binary with `sudo ./yourbinary`—using `sudo` only during the build step is not sufficient.
package orch

import (
	"errors"
	"fmt"
	"github.com/IM-Malik/Gonix/nginx"
	"github.com/IM-Malik/Gonix/nginx/sites/reverseproxy"
	"github.com/IM-Malik/Gonix/nginx/sites/webserver"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Defaults holds default information about the Nginx setup in your machine
type Defaults struct {
	// NginxConf is the path to nginx.conf file
	NginxConf string
	// SitesAvailable is the path to sites-available directory
	SitesAvailable string
	// SitesEnabled is the path to sites-enabled directory
	SitesEnabled string
	// ModulesEnabled is the path to modules-enabled directory
	ModulesEnabled string
}

// SetNginxConf set/reset only the NginxConf value in the defaults instance
func (defaults *Defaults) SetNginxConf(nginxConfPath string) {
	defaults.NginxConf = nginxConfPath
}

// SetSitesAvailable set/reset only the SitesAvailable value in the defaults instance
func (defaults *Defaults) SetSitesAvailable(sitesAvailablePath string) {
	defaults.SitesAvailable = sitesAvailablePath
}

// SetSitesEnabled set/reset only the SitesEnabled value in the defaults instance
func (defaults *Defaults) SetSitesEnabled(sitesEnabledPath string) {
	defaults.SitesEnabled = sitesEnabledPath
}

// SetModulesEnabled set/reset only the ModulesEnabled value in the defaults instance
func (defaults *Defaults) SetModulesEnabled(modulesEnabledPath string) {
	defaults.ModulesEnabled = modulesEnabledPath
}

// GetDefaults returns the values inside the defaults instance
func (defaults Defaults) GetDefaults() *Defaults {
	return &defaults
}

// SetAllDefaults set/reset all the values inside the default instance
func (defaults *Defaults) SetAllDefaults(nginxConfPath string, sitesAvailablePath string, sitesEnabledPath string, modulesEnabledPath string) (*Defaults, error) {
	if nginxConfPath != "" && sitesAvailablePath != "" && sitesEnabledPath != "" && modulesEnabledPath != "" {
		return &Defaults{
			NginxConf:      nginxConfPath,
			SitesAvailable: sitesAvailablePath,
			SitesEnabled:   sitesEnabledPath,
			ModulesEnabled: modulesEnabledPath,
		}, nil
	}
	return nil, fmt.Errorf("one or more of the parameters are not set")
}

// GetGlobalConfig returns the configuration inside the nginx.conf file
func GetGlobalConfig(defaults *Defaults) (string, error) {
	content, err := os.ReadFile(defaults.NginxConf)
	if err != nil {
		return "", fmt.Errorf("unable to read nginx.conf: %v", err)
	}
	return "\n" + string(content), nil
}

// GetSiteConfig returns the configuration inside the siteName.conf file
func GetSiteConfig(defaults *Defaults, siteName string) (string, error) {
	content, err := os.ReadFile(defaults.SitesAvailable + siteName + ".conf")
	if err != nil {
		return "", fmt.Errorf("failed to read site content: %v", err)
	}
	return "\n" + string(content), nil
}

// ReloadNginx execute the linux command "nginx -s reload" to reload Nginx without fully shutting down the Nginx service
// ReloadNginx is used when there is a new configuration file and nginx does not see it or if there is changes in existing configuration files
func ReloadNginx() (string, error) {
	cmd := exec.Command("nginx", "-s", "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("could not reload Nginx: %v", err)
	}
	return fmt.Sprintf("nginx is reloaded successfully: %v", string(output)), nil
}

// RestartNginx execute the linux command "systemctl restart nginx" to fully shutdown the Nginx service and starting it again
// RestartNginx is used when Nginx is not working properly or if there is any binary upgrades in Nginx
func RestartNginx() (string, error) {
	cmd := exec.Command("systemctl", "restart", "nginx")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("could not restart the Nginx service: %v", err)
	}
	return "nginx process is restarted successfully", nil
}

// TestNginx execute the linux command "nginx -t" to test if there is syntax errors in any Nginx configuratioin file
func TestNginx() (string, error) {
	cmd := exec.Command("nginx", "-t")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("could not reload Nginx: %v\nOutput: %s", err, string(output))
	}
	return fmt.Sprintf("nginx is testing successfully: %v", string(output)), nil
}

// CreateAndEnableRevProxy combine 3 functions (AddSite, EnableSite, AddUpstream) from the reverseproxy package to automate the new creation and enable of the Nginx configuration file
func CreateAndEnableRevProxy(defaults *Defaults, domain string, listenPort int, uri string, enableSSL bool, certPath string, keyPath string, upstreamName string, serverIP string, portNum int, httpOrHttps string) (string, error) {
	if httpOrHttps == "http" || httpOrHttps == "https" {
		_, err := reverseproxy.AddSite(defaults.SitesAvailable, domain, listenPort, upstreamName, uri, enableSSL, certPath, keyPath, httpOrHttps)
		if err != nil {
			return "", fmt.Errorf("the reverse proxy site was not created due to: %v", err)
		}
		_, err = reverseproxy.EnableSite(defaults.SitesAvailable, defaults.SitesEnabled, domain)
		if err != nil {
			return "", fmt.Errorf("the reverse proxy site was not enabled due to: %v", err)
		}
		_, err = nginx.AddUpstream(defaults.SitesAvailable, domain, upstreamName, serverIP, portNum)
		if err != nil {
			return "", fmt.Errorf("the reverse proxy site was not enabled due to: %v", err)
		}
	}
	return "the site was created and enabled successflly", nil
}

// CreateAndEnableWebServer combine 2 functions (AddSite, EnableSite) from the webserver package to automate the new creation and enable of the Nginx configuration file
func CreateAndEnableWebServer(defaults *Defaults, domain string, listenPort int, uri string, staticContentPath string, staticContentFileName string) (string, error) {
	_, err := webserver.AddSite(defaults.SitesAvailable, domain, listenPort, uri, staticContentPath, staticContentFileName)
	_, err1 := webserver.EnableSite(defaults.SitesAvailable, defaults.SitesEnabled, domain)
	if err != nil || err1 != nil {
		return "", fmt.Errorf("the web server was not created: %v", errors.Join(err, err1))
	}
	return "the site was created and enabled successflly", nil
}

// RemoveSite combine 2 functions (RemoveSite, RemoveEnabledSite) to automate the complete removal of an existing Nginx configuration file
func RemoveSite(defaults *Defaults, domain string) (string, error) {
	_, err := nginx.RemoveSite(defaults.SitesAvailable, domain)
	_, err1 := nginx.RemoveEnabledSite(defaults.SitesEnabled, domain)
	if err != nil || err1 != nil {
		return "", fmt.Errorf("failed to remove the site: %v", errors.Join(err, err1))
	}
	return "the site was removed successfully", nil
}

// BackupConfig creates a backup configuration file from an existing configuration file
// BackupConfig is used before modifying a configuration if anything went wrong
func BackupConfig(defaults *Defaults, domain string) (string, error) {
	srcFile, err := os.Open(defaults.SitesAvailable + domain + ".conf")
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer srcFile.Close()

	backupPath := filepath.Join(defaults.SitesAvailable, domain+".conf.bak")

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

// RollBackChanges revert to the backup configuration file and removes the old one
// RollBackChanges is used if anything went wrong when modifying the existing configuration file, with condition to have an existing backup file created by BackupConfig function
func RollBackChanges(defaults *Defaults, domain string) (string, error) {
	dir := filepath.Dir(defaults.SitesAvailable)
	base := filepath.Base(domain + ".conf")
	backupFile := filepath.Join(dir, base+".bak")
	oldFile := filepath.Join(dir, base)

	if _, err := os.Stat(backupFile); err != nil {
		return "", fmt.Errorf("backup file %s not found", backupFile)
	}
	if err := os.Remove(oldFile); err != nil {
		return "", fmt.Errorf("failed to remove the modified configuration file: %v", err)
	}
	if err := os.Rename(backupFile, oldFile); err != nil {
		return "", fmt.Errorf("could not restore the backup configuration: %v", err)
	}

	return fmt.Sprintf("rollback is successful at: %v", oldFile), nil
}

// UpdateSite update the existing configuration file
func UpdateSite(defaults *Defaults, domain string, oldText string, newText string) (string, error) {
	_, err := BackupConfig(defaults, domain)
	if err != nil {
		return "", fmt.Errorf("could not create a backup of the site configuration: %v", err)
	}
	data, err := os.ReadFile(defaults.SitesAvailable + domain + ".conf")
	if err != nil {
		return "", fmt.Errorf("site configuration file not found: %v", err)
	}
	if strings.Count(string(data), oldText) == 0 {
		return "", fmt.Errorf("the specified text to replace was not found in the site configuration file")
	}
	updated := strings.ReplaceAll(string(data), oldText, newText)
	err = os.WriteFile(defaults.SitesAvailable+domain+".conf", []byte(updated), 0644)
	if err != nil {
		RollBackChanges(defaults, domain)
		return "", fmt.Errorf("failed to update the site configuration file. Rolling back changes: %v", err)
	}
	return "the site file was updated successfully", nil
}
