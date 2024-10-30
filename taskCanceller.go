package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	apiLib "github.com/hornbill/goApiLib"
)

func processCommandLineOptions() {
	flag.StringVar(&gStrInstance, "instance", "", "Instance name")
	flag.StringVar(&gStrAPI, "api", "", "API Key")
	flag.StringVar(&gStrTaskList, "listfile", "", "File name of file containing list of task references - one per line")
	flag.StringVar(&gStrTaskRef, "taskref", "", "Single Task Reference")
	flag.BoolVar(&configAction, "delete", false, "REMOVAL of tasks (default: false - cancellation of tasks); this can be overriden by -action!")
	flag.BoolVar(&configVersion, "version", false, "Returns the tool verison")
	flag.StringVar(&gStrAction, "action", "", "cancel, delete, complete; overrides -delete option")
	flag.StringVar(&gStrOutcome, "outcome", "completed", "Outcome for completing actions (default: completed)")
	flag.StringVar(&gStrNote, "note", "", "note to be used for completion")
	flag.Parse()
}

func main() {
	processCommandLineOptions()
	//-- Load Configuration File Into Struct
	if configVersion {
		fmt.Printf("%v \n", version)
		return
	}

	if gStrInstance == "" {
		log.Fatal("instance argument is mandatory")
	}
	if gStrAPI == "" {
		log.Fatal("api argument is mandatory")
	}

	// default action is to cancel (can be overriden by -delete and THEN by -action)
	if configAction {
		taskAction = "taskDelete"
	}

	if gStrAction == "delete" {
		taskAction = "taskDelete"
	} else if gStrAction == "cancel" {
		taskAction = "taskCancel"
	} else if gStrAction == "complete" {
		taskAction = "taskComplete"
	} else if gStrAction != "" {
		log.Fatal("please provide a valid action MUST be given")
	}

	if taskAction == "" {
		log.Fatal("in order to process tasks an action MUST be given")
	}

	espXmlmc = apiLib.NewXmlmcInstance(gStrInstance)
	espXmlmc.SetAPIKey(gStrAPI)

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

	fmt.Println(countSuccess, " Tasks processed successfully")
	fmt.Println(countFail, " Task processing failures")
	if (countOutcome > 0){
		fmt.Println(countOutcome, " Task completions need to be reconfigured")
	}
	if (countNotesReq > 0){
		fmt.Println(countNotesReq, " Task completions require notes")
	}
}

// Time Struct
var Time struct {
	timeNow   string
	startTime time.Time
	endTime   time.Duration
}

func cancelTask(taskref string) {
	//espXmlmc.SetParam("appName", "com.hornbill.core")
	espXmlmc.SetParam("taskId", taskref)
	if taskAction == "taskComplete" {
		if gStrOutcome != "" {
			espXmlmc.SetParam("outcome", gStrOutcome)
		}
		if gStrNote != "" {
			espXmlmc.SetParam("details", gStrNote)
		}
	}
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
		error_msg := xmlRespon.State.ErrorRet
		fmt.Println("[" + taskref + "] METHOD FAIL " + taskAction + ": " + error_msg)
		lastX := error_msg[len(error_msg)-32:]
		if (lastX=="requires details to be specified" || lastX=="ied is not defined for this task") {
		
			if (lastX=="requires details to be specified") {
				countNotesReq++
			}
			if (lastX=="ied is not defined for this task") {
				countOutcome++
			}
		
			espXmlmc.SetParam("taskId", taskref)
			TaskDetails, giErr := espXmlmc.Invoke("task", "taskGetInfo")
			if xmlmcErr != nil {
				fmt.Println("[" + taskref + "] API FAIL Get Info: " + giErr.Error())
				return
			}
			var xmlRespon taskResponse
			err := xml.Unmarshal([]byte(TaskDetails), &xmlRespon)
			if err != nil {
				fmt.Println("[" + taskref + "] UNMARSHAL FAIL " + taskAction + ": " + err.Error())
				return
			}
			if xmlRespon.MethodResult != "ok" {
				fmt.Println("[" + taskref + "] METHOD FAIL Get Info: " + xmlRespon.State.ErrorRet)
				return
			}
			
			if (lastX=="requires details to be specified") {
				if (xmlRespon.Params.CompletionDetails != ""){

					fmt.Println("Attempting with stored completion details")

					espXmlmc.SetParam("taskId", taskref)
					if gStrOutcome != "" {
						espXmlmc.SetParam("outcome", gStrOutcome)
					}
					espXmlmc.SetParam("details", xmlRespon.Params.CompletionDetails)
					XMLLogin, xmlmcErr := espXmlmc.Invoke("task", taskAction)
					if xmlmcErr != nil {
						fmt.Println("[" + taskref + "] API FAIL " + taskAction + ": " + xmlmcErr.Error())
						return
					}

					var xmlRespon2 xmlmcResponse
					err := xml.Unmarshal([]byte(XMLLogin), &xmlRespon2)
					if err != nil {
						fmt.Println("[" + taskref + "] UNMARSHAL FAIL " + taskAction + ": " + err.Error())
						return
					}
					if xmlRespon2.MethodResult != "ok" { 
						fmt.Println("[" + taskref + "] METHOD FAIL " + taskAction + ": " + xmlRespon2.State.ErrorRet)
						return
					}




				}
				countNotesReq++
			}
			if (lastX=="ied is not defined for this task") {
				fmt.Println("Please select from: " + xmlRespon.Params.Outcomes)
				return
			}

		}
		countFail++
		return
	}
	fmt.Println("[" + taskref + "] SUCCESS " + taskAction)
	countSuccess++
	GCounter++
}
