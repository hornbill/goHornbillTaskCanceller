package main

import (
	"encoding/base64"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/hornbill/goApiLib"
	"time"
	// read file
	"bufio"
	// download inclusion
	"os"
	//DAV Inclusion
	//	_ "bytes"
	"net"
	"net/http"
)

const (
	version         = "1.0.0"
	regularFontSize = 10
)

var (
	GstrConfigFileName      string
	GstrZone                string
	GstrInstance            string
	GstrAPI                 string
	GboolProcessAttachments bool
	GstrUsername            string
	GstrPassword            string
	GConfigDetails          hbImportConfStruct
	xmlmcInstanceConfig     xmlmcConfigStruct
	espXmlmc                *apiLib.XmlmcInstStruct
	GstrTaskRef             string
	GstrCallList            string
	configDryRun            bool
	configAction            bool
	GCounter                int
	GstrOutputDir           string
	client                  = http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   600 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   1,
			Proxy:                 http.ProxyFromEnvironment,
		},
		//Timeout: time.Duration(120 * time.Second),
	}
	GApplicationStrings map[string]string
)

type stateStruct struct {
	Code     string `xml:"code"`
	ErrorRet string `xml:"error"`
}

type xmlmcResponse struct {
	MethodResult string      `xml:"status,attr"`
	State        stateStruct `xml:"state"`
}

type hbImportConfStruct struct {
	UserName   string
	Password   string
	InstanceID string
	APIKey     string
	DeleteTask bool
	taskAction string
	URL        string
}


func processCommandLineOptions() {
	flag.StringVar(&GstrZone, "zone", "eur", "Override the default Zone the instance sits in")
	flag.BoolVar(&configDryRun, "dryrun", false, "Skipping cancelation of tasks")
	flag.StringVar(&GstrInstance, "instance", "", "Instance name")
	flag.StringVar(&GstrAPI, "api", "", "API Key")
	flag.StringVar(&GstrUsername, "user", "", "Override the API Key to use username and password instead")
	flag.StringVar(&GstrPassword, "pwd", "", "Override the API Key to use username and password instead")
	flag.StringVar(&GstrCallList, "listfile", "", "File name of file containing list of call references - one per line")
	flag.StringVar(&GstrTaskRef, "taskref", "", "Single Call Reference")
	flag.BoolVar(&configAction, "delete", false, "REMOVAL of tasks (default: false - cancelation of tasks")
	flag.Parse()
}

type xmlmcConfigStruct struct {
	instance string
	url      string
	zone     string
	davurl   string
}

func getInstanceURL() string {
	xmlmcInstanceConfig.url = "https://"
	xmlmcInstanceConfig.url += xmlmcInstanceConfig.zone
	xmlmcInstanceConfig.url += "api.hornbill.com/"
	xmlmcInstanceConfig.url += xmlmcInstanceConfig.instance
	xmlmcInstanceConfig.davurl = xmlmcInstanceConfig.url + "/dav/"
	xmlmcInstanceConfig.url += "/xmlmc/"
	return xmlmcInstanceConfig.url
}

func main() {
	processCommandLineOptions()
	//-- Load Configuration File Into Struct

	if configAction {
		GConfigDetails.taskAction = "taskDelete"
	} else {
		GConfigDetails.taskAction = "taskCancel"
	}

	
	if GstrZone != "" {
		xmlmcInstanceConfig.zone = GstrZone
	}
	if GstrInstance != "" {
		xmlmcInstanceConfig.instance = GstrInstance
	} else {
		//xmlmcInstanceConfig.instance = GConfigDetails.InstanceID
		fmt.Println("In order to cancel tasks an Instance Name MUST be given")
		return
	}
	if GstrAPI != "" {
		GConfigDetails.APIKey = GstrAPI
	} else {
		fmt.Println("In order to cancel tasks an API Key Name MUST be given")
		return
	}

	GConfigDetails.URL = getInstanceURL()
	//fmt.Println(GConfigDetails.URL)

	espXmlmc = apiLib.NewXmlmcInstance(GConfigDetails.URL)

	if GstrUsername != "" && GstrPassword != "" {
		var boolLogin = login()
		if boolLogin != true {
			fmt.Println("Unable to Login")
			return
		}
		//-- Defer log out of Hornbill instance until after main() is complete
		defer logout()
	} else {
		espXmlmc.SetAPIKey(GConfigDetails.APIKey)
	}
	
	if GstrTaskRef != "" {
		if !configDryRun {
			cancelTask(GstrTaskRef)
		} else {
			GCounter++
		}
	} else if GstrCallList != "" {
		file, err := os.Open(GstrCallList)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if !configDryRun {
				cancelTask(scanner.Text())
			} else {
				GCounter++
			}
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	} else {
		fmt.Println("No input given. Either give a Task Reference (-taskref=TSK###) or a file to read (-listfile=somefile.txt)")
	}

	if configDryRun {
		fmt.Println(fmt.Sprintf("%d", GCounter) + " items would have been cancelled")
	} else {
		fmt.Println(fmt.Sprintf("%v", GCounter) + " items cancelled")
	}
}

func logout() {
	espXmlmc.Invoke("session", "userLogoff")
}

func login() bool {
	//	espXmlmc = apiLib.NewXmlmcInstance(swImportConf.HBConf.URL)

	espXmlmc.SetParam("userId", GstrUsername)
	espXmlmc.SetParam("password", base64.StdEncoding.EncodeToString([]byte(GstrPassword)))
	XMLLogin, xmlmcErr := espXmlmc.Invoke("session", "userLogon")
	if xmlmcErr != nil {
		panic("Logon Issue : " + fmt.Sprintf("%v", xmlmcErr))
	}

	var xmlRespon xmlmcResponse
	err := xml.Unmarshal([]byte(XMLLogin), &xmlRespon)
	if err != nil {
		panic("Unable to Login: " + fmt.Sprintf("%v", err))
		return false
	}
	if xmlRespon.MethodResult != "ok" {
		panic("Unable to Login: " + xmlRespon.State.ErrorRet)
		return false
	}
	return true
}

func cancelTask(taskref string) bool {
	//espXmlmc.SetParam("appName", "com.hornbill.core")
	espXmlmc.SetParam("taskId", taskref)
	XMLLogin, xmlmcErr := espXmlmc.Invoke("Task", GConfigDetails.taskAction)
	if xmlmcErr != nil {
		fmt.Println("ApplicationString Issue : " + fmt.Sprintf("%v", xmlmcErr))
		return false
	}

	var xmlRespon xmlmcResponse
	err := xml.Unmarshal([]byte(XMLLogin), &xmlRespon)
	if err != nil {
		fmt.Println("Unable to obtain " + GConfigDetails.taskAction + ": " + fmt.Sprintf("%v", err))
		return false
	}
	if xmlRespon.MethodResult != "ok" {
		fmt.Println("Task not " + GConfigDetails.taskAction + ": " + xmlRespon.State.ErrorRet)
		return false
	}
	GCounter++
	return true
}
