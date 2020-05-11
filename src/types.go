package main

type ConfigTypes struct {
	Settings struct {
		Code    string `yaml:"code"`
		Default_Pin int `yaml:"Pin"`
	} `yaml:"settings"`
	Primary struct {
		Name string `yaml:"name"`
		Mac  string `yaml:"mac"`
		Ip   string `yaml:"ip"`
	} `yaml:"primary"`
	Secondary struct {
		Name string `yaml:"name"`
		Mac  string `yaml:"mac"`
		Ip   string `yaml:"ip"`
	} `yaml:"secondary"`
	Tertiary struct {
		Name string `yaml:"name"`
		Mac  string `yaml:"mac"`
		Ip   string `yaml:"ip"`
	} `yaml:"tertiary"`
}

type FailureNetwork struct {
	Time         string `json:"time"`
	Failure_type string `json:"type_of_failure"`
}

type RequestImage struct {
	Request_id string `json:"request_id"`
	Time_from string `json:"time_from"`
	Time_to string `json:"time_to"`
}

type RequestData struct {
	Request_id string `json:"request_id"`
	Time_from string `json:"time_from"`
	Time_to string `json:"time_to"`
	Type string `json:"type"`
}

type RequestDatabase struct {
	Request_id int `json:"request_id"`
	Time_from string `json:"time_from"`
	Time_to string `json:"time_to"`
	Type string `json:"type"`
}

type DataInfo struct {
	Id int `json:"id"`
    Message string `json:"message"`
	Time string `json:"Time"`
}

type DeviceFound struct {
	Device_name string `json:"name"`
	Mac string `json:"mac"`
	Ip_address string `json:"address"`
	Alive bool `json:"alive"`
}

type AccessResponse struct {
	Id int `json:"id"`
	Result string `json:"result"`
}

type RequestAccess struct {
	Id int `json:"id"`
	Pin int `json:"pin"`
}

type UnauthorisedConnection struct {
	Mac string `json:"mac"`
	Time string `json:"time"`
	Alive bool `json:"alive"`
}

type EventNAC struct {
	Component    string `json:"Component"`
	Message      string `json:"Message"`
	Time         string `json:"Time"`
}

type MapMessage struct {
	message     string
	routing_key string
	time        string
	valid       bool
}

//Topics
const REQUESTDATA string = "Request.Data"
const AUTHENTICATIONREQUEST string = "Authentication.Request"
const DATAINFO string = "Data.Info"
const REQUESTACCESS string = "Request.Access"
//
const FAILURENETWORK string = "Failure.Network"
const EVENTNAC string = "Event.NAC"
const REQUESTDATABASE string = "Request.Database"
const DATARESPONSE string = "Data.Response"
const DEVICEFOUND string = "Device.Found"
const AUTHENTICATIONRESPONSE string = "Authentication.Response"
const ACCESSRESPONSE string = "Access.Response"
const UNAUTHORISEDCONNECTION string = "Unauthorised.Connection"
//
const ACCESSFAIL string = "FAIL"
const ACCESSPASS string = "PASS"
const EXCHANGENAME string = "topics"
const EXCHANGETYPE string = "topic"
const TIMEFORMAT string = "20060102150405"
const COMPONENT string = "NAC"
const FAILUREPUBLISH string = "Failed to publish"
const SERVERERROR string = "Server is failing to send"

var SubscribedMessagesMap map[uint32]*MapMessage
var DevicesList map[uint32]*DeviceFound
var key_id uint32 = 0
var device_id uint32 = 0