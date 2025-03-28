package config

import (
	"log"
	"os"
    "fmt"
)

func GenerateDefaultGlobalConfig(env_filePath string) (error) {
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
	file, err := os.OpenFile(env_filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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

func GenerateDefaultEmailConfig(env_filePath string) (error) {
	defaultEmailConfig := `mail {
      # See sample authentication script at:
      # http://wiki.nginx.org/ImapAuthenticateWithApachePhpScript

      # auth_http localhost/auth.php;
      # pop3_capabilities "TOP" "USER";
      # imap_capabilities "IMAP4rev1" "UIDPLUS";

      server {
              listen     localhost:110;
              protocol   pop3;
              proxy      on;
      }

      server {
              listen     localhost:143;
              protocol   imap;
              proxy      on;
      }
}
`
	file, err := os.OpenFile(env_filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		// log.Fatalf("failed to open the nginx.conf file: %v\n", err)
		return fmt.Errorf("failed to open the nginx.conf file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(defaultEmailConfig)
	if err != nil {
		// log.Fatalf("failed to write in the nginx.conf file: %v\n", err)
		return fmt.Errorf("failed to write in the nginx.conf file: %v", err)
	} else {
		log.Printf("the email default configuration is written correctly in nginx.conf file\n")
		return nil
	}
}
