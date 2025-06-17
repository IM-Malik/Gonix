// this is the CLI entry point
package main

import (
	// "github.com/IM-Malik/Gonix/nginx"
	// "fmt"

	// "github.com/IM-Malik/Gonix/nginx/config"
	// "fmt"
	"log"
	// "github.com/IM-Malik/Gonix/nginx/sites/reverseproxy"
	// "github.com/IM-Malik/Gonix/nginx/modules"
	"github.com/IM-Malik/Gonix/nginx/orch"
	// "github.com/IM-Malik/Gonix/nginx/sites/webserver"
)

func main() {


	
	// ===============================================================
	// ===============================================================
	// ===============================================================
	// 				BEFORE ASKING FOR PEOPLE TO TEST!!!!
	// ===============================================================
	// ===============================================================
	// ===============================================================





	// filePath := os.Getenv("NGINX_CONF_PATH")
	// globalConfigFilePath := "/etc/nginx/nginx.conf"
	// siteConfigFilePath := "/etc/nginx/sites-available/malik.com.conf"
	// sitesPath := "/etc/nginx/sites-available/"
	// log.Println(globalConfigFilePath)

	// err := config.GenerateDefaultGlobalConfig(globalConfigFilePath)
	// if err != nil {
	// 	log.Println(err)
	// }

	// // err = config.GenerateDefaultEmailConfig(filePath)
	// // if err != nil {
	// // 	log.Println(err)
	// // }
	

	// output, err = nginx.ReloadNginx()
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	log.Print(output)
	// }

	// output, err = nginx.RestartNginx()
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	log.Println(output)
	// }

	// output, err = nginx.GetGlobalConfig(globalConfigFilePath)
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	log.Println(output)
	// }

	// output, err = nginx.GetSiteConfig(siteConfigFilePath)
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	log.Println(output)
	// }

	// // output, err = nginx.BackupConfig(siteConfigFilePath)
	// // if err != nil {
	// // 	log.Println(err)
	// // } else {
	// // 	log.Println(output)
	// // }

	// output, err = nginx.RollbackChanges(siteConfigFilePath)
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	log.Println(output)
	// }

	    // serverBlockVariables := sites.Config{
    //     Domain: "malik.com",
    //     ListenPort: 80,
    //     RootDir: "/var/www/html",
    //     IndexFiles: ["index.html"],
    //     ProxyPass: "http://localhost:3000",
    //     EnableSSL: ture,
    //     SSLCertPath: "/etc/letsencrypt/nitaqat/cert",
    //     SSLKeyPath: "/etc/letsencrypt/nitaqat/key",
    //     URLPath: "api",
    // }

	// output1, err1 := reverseproxy.AddSite("/etc/nginx/sites-available/", "malik.com", 80, "http://localhost:3000", "/api", "/etc/letsencrypt/nitaqat/cert", "/etc/letsencrypt/nitaqat/key", false)
	// output2, err2 := reverseproxy.EnableSite("/etc/nginx/sites-available/", "/etc/nginx/sites-enabled/", "malik.com")
	// fmt.Println(output1)
	// fmt.Println(output2)
	// fmt.Println(err1)
	// fmt.Println(err2)
	
	// output1, err1 := reverseproxy.RemoveSite("/etc/nginx/sites-available/", "malik.com")
	// output2, err2 := reverseproxy.RemoveEnabledSite("/etc/nginx/sites-enabled/", "malik.com")
	// fmt.Println(output1)
	// fmt.Println(output2)
	// fmt.Println(err1)
	// fmt.Println(err2)
	
	// output, err := reverseproxy.AddLocation("/etc/nginx/sites-available/", "malik.com", "http://localhost:3000", "/add")
	// fmt.Println(output)
	// fmt.Println(err)
	
	// output, err = reverseproxy.AddUpstream("/etc/nginx/sites-available/", "malik.com", "test", "127.0.0.1", 9086)
	// fmt.Println(output)
	// fmt.Println(err)

	// output, err := webserver.AddSite("/etc/nginx/sites-available/", "malik.com", 80, "/api", "/usr/share/nginx/html", "index.html index.htm")
	// fmt.Println(output)
	// fmt.Println(err)
	
	// output, err = nginx.TestNginx()
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	log.Print(output)
	// }
	// k, err := config.GenerateDefaultStreamConfig("/etc/nginx", "malik.com", "127.0.0.1", 4598)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(k)
	// err := config.GenerateDefaultEmailConfig("/etc/nginx/nginx.conf")
	// if err != nil {
	// 	log.Println(err)
	// }
	// output, err := reverseproxy.AddLocation("/etc/nginx/sites-available/", "malik.com", "http://127.0.0.1:7282", "click")
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	log.Println(output)
	// }
	// output, err := modules.DisableModule("/etc/nginx/modules-enabled/", "50-mod-http-auth-pam.conf")
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(output)
	// output, err := modules.EnableModule("/usr/share/nginx/modules-available/", "/etc/nginx/modules-enabled/", "50-mod-http-auth-pam.conf")
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(output)
	// err := modules.GetEnabledModules("/etc/nginx/modules-enabled/")
	// if err != nil {
	// 	log.Println(err)
	// }
	// err := reverseproxy.GetAvailableSites("/etc/nginx/sites-available/")
	// if err != nil {
	// 	log.Println(err)
	// }
	k, err := orch.SetAllDefaults("/etc/nginx/", "/etc/nginx/sites-available/", "/etc/nginx/sites-enabled/", "/etc/nginx/modules-enabled/")
	if err != nil {
		log.Println(err)
	}
	
	// j, e := RemoveSite(k, "ali.com")
	// if e != nil {
	// 	log.Print(e)
	// }
	// log.Println(j)

	// l, er := CreateAndEnableRevProxy(k, "ali.com", 80, "/docs", false, "", "", "docsUpstream", "127.0.0.1", 4319, "http")
	// if er != nil {
	// 	log.Println(er)
	// }
	// log.Println(l)

	// a, e1 := BackupConfig(k, "ali.com")
	// if e1 != nil {
	// 	log.Println(e1)
	// }
	// log.Println(a)

	// a, e1 := RollbackChanges(k, "ali.com")
	// if err != nil {
	// 	log.Println(e1)
	// }
	// log.Println(a)

	a, err := UpdateSite(k, "ali.com", "8083", "80gd")
	if err != nil {
		log.Println(err)
	}
	log.Println(a)
	res, err := TestNginx()
	if err != nil {
		log.Println(err)
	}
	log.Println(res)

	k.SetModulesEnabled("")
}

