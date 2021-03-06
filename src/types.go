package main

type ConfigTypes struct {
	Settings struct {
		Code string `yaml:"code"`
	} `yaml:"settings"`
}

type FailureNetwork struct {
	Time         string `json:"time"`
	Failure_type string `json:"type_of_failure"`
}

type ApiResponse struct {
	Vendor string `json:"Vendor"`
}

type MapMessage struct {
	message     string
	routing_key string
	time        string
	valid       bool
}

type Device struct {
	Device_name string `json:"name"`
	Mac         string `json:"mac"`
	Ip_address  string `json:"address"`
	Alive       bool   `json:"alive"`
	Allowed     int    `json:"allowed"`
	New         bool   `json:"new"`
}

type DeviceFound struct {
	Device_name string `json:"name"`
	Ip_address  string `json:"address"`
	Status      int    `json:"status"`
}

type DeviceRequest struct {
	Request_id uint32 `json:"id"`
	Name       string `json:"name"`
	Mac        string `json:"mac"`
}

type DeviceResponse struct {
	Request_id uint32 `json:"id"`
	Name       string `json:"name"`
	Mac        string `json:"mac"`
	Status     string `json:"status"`
}

//
const DEVICERESPONSE string = "Device.Response"
const FAILURENETWORK string = "Failure.Network"
const DEVICEFOUND string = "Device.Found"
const DEVICEREQUEST string = "Device.Request"

//
const ACCESSFAIL string = "FAIL"
const ACCESSPASS string = "PASS"
const EXCHANGENAME string = "topics"
const EXCHANGETYPE string = "topic"
const TIMEFORMAT string = "2006/01/02 15:04:05"
const FAILUREPUBLISH string = "Failed to publish"
const UNKNOWN_DEVICE string = "New device connected - "

//
const ALLOWED int = 1
const BLOCKED int = 2
const DISCOVERED int = 3
const UNKNOWN int = 4
const ALLOWED_STRING string = "ALLOWED"
const BLOCKED_STRING string = "BLOCKED"
const DISCOVERED_STRING string = "DISCOVERED"
const UNKNOWN_STRING string = "UNKNOWN"

//
const START_ADDRESS string = "192.168.0."

//
var SubscribedMessagesMap map[uint32]*MapMessage
var DevicesList map[uint32]*Device
var key_id uint32 = 0
var device_id uint32 = 0
