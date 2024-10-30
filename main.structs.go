package main

import apiLib "github.com/hornbill/goApiLib"

const (
	version = "2.0.0"
)

var (
	configVersion bool
	configAction  bool
	gStrInstance  string
	gStrAPI       string
	gStrTaskRef   string
	gStrTaskList  string
	espXmlmc      *apiLib.XmlmcInstStruct
	gStrAction    string
	gStrNote      string
	gStrOutcome   string
	taskAction    = "taskCancel"
	GCounter      int
	GstrOutputDir string
	countSuccess  int
	countFail     int
	countOutcome int
	countNotesReq int
)

type stateStruct struct {
	Code     string `xml:"code"`
	ErrorRet string `xml:"error"`
}

type xmlmcResponse struct {
	MethodResult string      `xml:"status,attr"`
	State        stateStruct `xml:"state"`
}

type taskResponse struct {
	MethodResult string      `xml:"status,attr"`
	State        stateStruct `xml:"state"`
	Params        taskInfoStruct `xml:"params"`
}

type taskInfoStruct struct {
	Outcomes string      `xml:"outcomes"`
	CompletionDetails string      `xml:"completionDetails"`
}