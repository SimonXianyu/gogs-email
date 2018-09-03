package main

/*
This file is the executable file entry
@Author Simon Xianyu
*/
import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
)

type EmailAccount struct {
	Host     string `yml:"host"`
	Address  string `yml:"address"`
	Secure   bool   `yml:"secure"`
	Account  string `yml:"account"`
	Password string `yml:"Password"`
	From     string `yml:"From"`
}

var account EmailAccount

var ToAddresses = []string{
	"simon_xianyu@163.com",
}

func HandleEvent(outResp http.ResponseWriter, req *http.Request) {
	auth := smtp.PlainAuth("", account.Account, account.Password, account.Host)
	var msg string
	msg = "Gogs Notification on event :" + req.Header.Get("X-Gogs-Event")

	sendErr := smtp.SendMail(account.Address, auth, account.Account, ToAddresses, []byte(msg))
	if sendErr != nil {
		outResp.WriteHeader(http.StatusInternalServerError)
		outResp.Header().Add("Content-Type", "application/json")
		fmt.Fprint(outResp, "{\"success\":false,\"msg\":\"Fail to send notify mail: "+sendErr.Error()+"\" }")
		log.Printf("sendError %s", sendErr)
	} else {
		fmt.Fprint(outResp, "{\"success\":true,\"msg\":\"Succeeded to send notify mail \" }")
	}
}

var serverPort = 9182

func main() {

	// Read and parse config file in Yaml format
	configFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal("Failed to read config file")
		return
	}
	yaml.Unmarshal(configFile, &account)
	//fmt.Println(account)

	http.HandleFunc("/onEvent", HandleEvent)

	fmt.Println("Start to listen at port", serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), nil))
}
