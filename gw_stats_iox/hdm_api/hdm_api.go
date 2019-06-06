package hdm_api
//package main

import (
	"net/http"
	"io/ioutil"
	"crypto/tls"
	"encoding/json"
	"fmt"
	//"os"
	//"os/exec"
	"strings"
	//"regexp"
	"golang.org/x/crypto/ssh"
	"net"
	"bytes"
	"../iox_conf"

)

type dec struct {
	Output     string      `json:"output"`
	RetCode     string    `json:"return-code"`
}

// Below functionality for source_env and oauth_creds has been removed.  Can Delete functions

/*func source_env() {
	cmd := exec.Command("sh", "-c", "source /data/.env")
	cmd.Run()
	fmt.Println(os.Getenv("OAUTH_CLIENT_ID"))
}

func oauth_creds() (url string, ID string, secret string){
	dt_creds, err := ioutil.ReadFile(".env")
	if err != nil {
		panic(err)
	}

	//fmt.Println(string(dt))
	re0 := regexp.MustCompile(`OAUTH_TOKEN_URL=(\S+)`)
	url = re0.FindAllStringSubmatch(string(dt_creds), -1)[0][1]
	//fmt.Println(url)

	re1 := regexp.MustCompile(`OAUTH_CLIENT_ID=(\S+)`)
	ID = re1.FindAllStringSubmatch(string(dt_creds), -1)[0][1]
	//fmt.Println(ID)

	re2 := regexp.MustCompile(`OAUTH_CLIENT_SECRET=(\S+)`)
	secret = re2.FindAllStringSubmatch(string(dt_creds), -1)[0][1]
	//fmt.Println(secret)

	return url, ID, secret
}*/

// This function executes terminal commands on the SSH console session to the main router.

func Show_cmd_ssh(command string) (string){
	 cfg, err := iox_conf.Conf_file()

         hostname := cfg.Section("ir_router_info").Key("IP").String()
         port := cfg.Section("ir_router_info").Key("port").String()

         username := cfg.Section("ir_router_info").Key("user").String()
	 password := cfg.Section("ir_router_info").Key("pass").String()
         config := &ssh.ClientConfig{
                 User: username,
                 Auth: []ssh.AuthMethod{ssh.Password(password)},
                 HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
                         return nil
                 },
         }

         hostaddress := strings.Join([]string{hostname, port}, ":")
         client, err := ssh.Dial("tcp", hostaddress, config)
         if err != nil {
                 fmt.Println(err.Error())
		 return "Error"
         }else {

		 session, err := client.NewSession()
		 if err != nil {
			 fmt.Println(err.Error())
			 return "Error"
		 }else {
			 //defer session.Close()

			 // fmt.Println("To exit this program, hit Control-C")
			 // fmt.Printf("Enter command to execute on %s : ", hostname)

			 // fmt.Scanf is unable to accept command with parameters
			 // see solution at
			 // https://www.socketloop.com/tutorials/golang-accept-input-from-user-with-fmt-scanf-skipped-white-spaces-and-how-to-fix-it
			 //fmt.Scanf("%s", &cmd)


			 cmd := command
			 //log.Printf(cmd)
			 //fmt.Println("Executing command ", cmd)

			 // capture standard output
			 // will NOT be able to handle refreshing output such as TOP command
			 // executing top command will result in panic
			 var buff bytes.Buffer
			 session.Stdout = &buff
			 if err := session.Run(cmd); err != nil {
				 session.Close()
				 fmt.Println(err.Error())
				 return "Error"
			 }else {

				 ssh_return := buff.String()
				 //fmt.Println(ssh_return)
				 session.Close()

				 return ssh_return
			 }
		 }
	 }
}

// This function executes terminal commands on the WiFi AP terminal session.

func Show_cmd_ssh_ap(ip string, command string) (string){
	 cfg, err := iox_conf.Conf_file()
	 var hostname string

	 if ip == "unassigned" {
		 hostname = cfg.Section("ap_info").Key("ap_IP").String()
	 }else {
		 hostname = ip
	 }
     port := cfg.Section("ap_info").Key("ap_port").String()

     username := cfg.Section("ap_info").Key("ap_user").String()
	 password := cfg.Section("ap_info").Key("ap_pass").String()
     config := &ssh.ClientConfig{
                 User: username,
                 Auth: []ssh.AuthMethod{ssh.Password(password)},
                 HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
                         return nil
                 },
         }
         //fmt.Println("\nConnecting to ", hostname, port)

         hostaddress := strings.Join([]string{hostname, port}, ":")
         client, err := ssh.Dial("tcp", hostaddress, config)
         if err != nil {
                 fmt.Println(err.Error())
		 return "Error"
         }else {

		 session, err := client.NewSession()
		 if err != nil {
			 fmt.Println(err.Error())
			 return "Error"
		 }else {
			 //defer session.Close()

			 // fmt.Println("To exit this program, hit Control-C")
			 // fmt.Printf("Enter command to execute on %s : ", hostname)

			 // fmt.Scanf is unable to accept command with parameters
			 // see solution at
			 // https://www.socketloop.com/tutorials/golang-accept-input-from-user-with-fmt-scanf-skipped-white-spaces-and-how-to-fix-it
			 //fmt.Scanf("%s", &cmd)


			 cmd := command
			 //log.Printf(cmd)
			 //fmt.Println("Executing command ", cmd)

			 // capture standard output
			 // will NOT be able to handle refreshing output such as TOP command
			 // executing top command will result in panic
			 var buff bytes.Buffer
			 session.Stdout = &buff
			 if err := session.Run(cmd); err != nil {
				 session.Close()
				 fmt.Println(err.Error())
				 return "Error"
			 }else {

				 ssh_return := buff.String()
				 //fmt.Println(ssh_return)
				 session.Close()

				 return ssh_return
			 }
		 }
	 }
}


// This function executes terminal commands through the HDM API interface that the IOx team released.
func Show_cmd(cmd string) (string) {
	cfg, err := iox_conf.Conf_file()

	url := cfg.Section("hdm_api_info").Key("url").String()

	payload := strings.NewReader(cmd)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Authorization", cfg.Section("hdm_api_info").Key("auth_token").String())

	tr := &http.Transport{
	    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return "Error"
	}else {
		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		var decode dec

		strbody := string(body)

		json.Unmarshal([]byte(strbody), &decode)

		output := decode.Output

		return output
	}

}
