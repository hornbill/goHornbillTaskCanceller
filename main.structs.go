package main

import apiLib "github.com/hornbill/goApiLib"

const (
	version = "1.1.0"
)

var (
	configVersion bool
	gStrInstance  string
	gStrAPI       string
	taskAction    = "taskCancel"
	espXmlmc      *apiLib.XmlmcInstStruct
	gStrTaskRef   string
	gStrTaskList  string
	configDryRun  bool
	configAction  bool
	countSuccess  int
	countFail     int
)

type stateStruct struct {
	Code     string `xml:"code"`
	ErrorRet string `xml:"error"`
}

type xmlmcResponse struct {
	MethodResult string      `xml:"status,attr"`
	State        stateStruct `xml:"state"`
}
