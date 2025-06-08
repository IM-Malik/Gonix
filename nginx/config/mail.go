package config

import (
	"fmt"
	"os"
)

// just finished with stream and tested it. mail is not tested yet
func GenerateDefaultEmailConfig(globalConfigFilePath string) error {
	defaultEmailConfig := `mail {
	  auth_http 127.0.0.1:9000/cgi-bin/nginxauth.cgi;
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
	file, err := os.OpenFile(globalConfigFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
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
		fmt.Printf("the email default configuration is written correctly in nginx.conf file\n")
		return nil
	}
}
