package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"

	apiLib "github.com/hornbill/goApiLib"
)

func processCommandLineOptions() {
	flag.StringVar(&gStrInstance, "instance", "", "Instance name")
	flag.StringVar(&gStrAPI, "api", "", "API Key")
	flag.StringVar(&gStrTaskList, "listfile", "", "File name of file containing list of task references - one per line")
	flag.StringVar(&gStrTaskRef, "taskref", "", "Single Task Reference")
	flag.BoolVar(&configAction, "delete", false, "REMOVAL of tasks (default: false - cancellation of tasks")
	flag.BoolVar(&configVersion, "version", false, "Returns the tool verison")
	flag.Parse()
}

func main() {
	processCommandLineOptions()
	//-- If configVersion just output version number and die
	if configVersion {
		fmt.Printf("%v \n", version)
		return
	}
	if configAction {
		taskAction = "taskDelete"
	}
	//Setup instance connection
	if gStrInstance == "" {
		log.Fatal("instance argument is mandatory")
	}
	if gStrAPI == "" {
		log.Fatal("api argument is mandatory")
	}
	espXmlmc = apiLib.NewXmlmcInstance(gStrInstance)
	espXmlmc.SetAPIKey(gStrAPI)

	//Process Tasks
	if gStrTaskRef != "" {
		cancelTask(gStrTaskRef)
	} else if gStrTaskList != "" {
		file, err := os.Open(gStrTaskList)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			cancelTask(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("No input given. Either give a Task Reference (-taskref=TSK###) or a file to read (-listfile=somefile.txt)")
	}
	fmt.Println(countSuccess, " Tasks Cancelled Successfully")
	fmt.Println(countFail, " Task Cancellation Failures")
}

func cancelTask(taskref string) {
	espXmlmc.SetParam("taskId", taskref)
	XMLLogin, xmlmcErr := espXmlmc.Invoke("task", taskAction)
	if xmlmcErr != nil {
		fmt.Println("[" + taskref + "] API FAIL " + taskAction + ": " + xmlmcErr.Error())
		countFail++
		return
	}
	var xmlRespon xmlmcResponse
	err := xml.Unmarshal([]byte(XMLLogin), &xmlRespon)
	if err != nil {
		fmt.Println("[" + taskref + "] UNMARSHAL FAIL " + taskAction + ": " + err.Error())
		countFail++
		return
	}
	if xmlRespon.MethodResult != "ok" {
		fmt.Println("[" + taskref + "] METHOD FAIL " + taskAction + ": " + xmlRespon.State.ErrorRet)
		countFail++
		return
	}
	fmt.Println("[" + taskref + "] SUCCESS " + taskAction)
	countSuccess++
	return
}
