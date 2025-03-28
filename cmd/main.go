// this is the CLI entry point
package main

import (
	"github.com/IM-Malik/Gonix/nginx"
	"github.com/IM-Malik/Gonix/nginx/config"
	"log"
)

func main() {
	// filePath := os.Getenv("NGINX_CONF_PATH")
	globalConfigFilePath := "/etc/nginx/nginx.conf"
	siteConfigFilePath := "/etc/nginx/sites-available/malik.com.conf"
	log.Println(globalConfigFilePath)

	err := config.GenerateDefaultGlobalConfig(globalConfigFilePath)
	if err != nil {
		log.Println(err)
	}

	// err = config.GenerateDefaultEmailConfig(filePath)
	// if err != nil {
	// 	log.Println(err)
	// }
	output, err := nginx.TestNginx()
	if err != nil {
		log.Println(err)
	} else {
		log.Print(output)
	}

	output, err = nginx.ReloadNginx()
	if err != nil {
		log.Println(err)
	} else {
		log.Print(output)
	}

	output, err = nginx.RestartNginx()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(output)
	}

	output, err = nginx.GetGlobalConfig(globalConfigFilePath)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(output)
	}

	output, err = nginx.GetSiteConfig(siteConfigFilePath)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(output)
	}

	// output, err = nginx.BackupConfig(siteConfigFilePath)
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	log.Println(output)
	// }

	output, err = nginx.RollbackChanges(siteConfigFilePath)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(output)
	}
}
