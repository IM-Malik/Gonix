package config

import (
	"log"
	"os"
    "fmt"
)

func GenerateDefaultGlobalConfig(globalConfigFilePath string) (error) {
	defaultConfig := `user  www-data;
worker_processes  auto;
pid        /run/nginx.pid;
include    /etc/nginx/modules-enabled/*.conf;

events {
    worker_connections  768;
    # multi_accept on;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    # Logging settings
    access_log  /var/log/nginx/access.log;
    error_log   /var/log/nginx/error.log warn;

    # Performance optimizations
    sendfile        on;
    tcp_nopush      on;
    tcp_nodelay     on;
    keepalive_timeout  65;
    types_hash_max_size 2048;

    # Gzip compression
    gzip on;
    # gzip_disable "msie6";

    # Include additional configurations (site-specific and extra settings)
    include /etc/nginx/conf.d/*.conf;
    include /etc/nginx/sites-enabled/*;
}
`
	// filePath := os.Getenv("NGINX_CONF_PATH")
	file, err := os.OpenFile(globalConfigFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open the nginx.conf file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to write in the nginx.conf file: %v", err)
	} else {
		log.Printf("the default configuration is written correctly in nginx.conf file\n")
		return nil
	}
}



func GlobalRateLimiting() () {

}

func GlobalServer() () {

}