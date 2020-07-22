package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"math/rand"
	"time"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var _port string
var _device_status string
var _messages_done bool
var _guid string
var _sent bool

type allData []DataInfo

var data_messages = allData{
}

func init() {
	_sent = false
	_port = "0"
}

func SetGUID() {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, 10)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
	_guid = string(b)
	log.Warn("GUID set to: ", _guid)
}

func SetPort(port string) {
	log.Debug("Port set")
	_port = port
}

func isValidGUID(guid string) bool {
	log.Warn("GUID Check")
	if guid == _guid {
		return true
	}
	log.Warn("InValid GUID")
	return false
}

func isValidRequest(id int) bool {
	log.Debug("Checking whether request id was new and valid")
	if id == current_id {
		log.Warn("Already received this, invalidating")
		return false
	} else {
		current_id = id
		_statusNAC.TimeEscConnected = getTime()
		_messages_done = false
		data_messages = nil
		return true
	}
}

func device_add(w http.ResponseWriter, r *http.Request) {
	var device DeviceAdd
	req_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		json.Unmarshal(req_body, &device)
		if isValidGUID(device.GUID) && isValidRequest(device.Request_id) {
			_statusNAC.TimeEscConnected = getTime()
			log.Debug("Received Device Name: ", device.Name)
			PublishDeviceUpdate(device.Name, device.Mac,
								device.Status, r.Method)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func user_add(w http.ResponseWriter, r *http.Request) {
	var user UserAdd
	req_body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		json.Unmarshal(req_body, &user)
		if isValidGUID(user.GUID) && isValidRequest(user.Request_id) {
			_statusNAC.TimeEscConnected = getTime()
			log.Debug("Received User Name: ", user.User)
			//PublishUserUpdate(user.User, user.Pin, r.Method)
			w.WriteHeader(http.StatusOK)
		} else {
			log.Error("Invalid GUID")
		}
	}
}

func getData(w http.ResponseWriter, r *http.Request) {
	log.Warn("Received data message:", r.URL.Query())

	guid := r.URL.Query().Get("guid")
	if guid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	string_id := r.URL.Query().Get("id")
	id, err := strconv.Atoi(string_id)
	if string_id == "" || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if isValidGUID(guid) && isValidRequest(id) {
		time_from  := r.URL.Query().Get("time_from")
		time_to    := r.URL.Query().Get("time_to")
		event_type := r.URL.Query().Get("event_type")
		PublishRequestDatabase(id, time_from, time_to, event_type)
		loop := 0
		for _messages_done == false && loop < 3 {
			time.Sleep(1 * time.Second)
			log.Warn("Loop: ", loop)
			loop++
		}
		if _messages_done {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(data_messages)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func http_server() {
	if _sent == false {
		PublishGUIDUpdate(_guid)
		_sent = true
	}
	router := mux.NewRouter().StrictSlash(true)
	// Set up of methods
	router.HandleFunc("/device", device_add).Methods("POST")
	router.HandleFunc("/device", device_add).Methods("PATCH")
	router.HandleFunc("/device", device_add).Methods("DELETE")
	//
	router.HandleFunc("/user", user_add).Methods("POST")
	router.HandleFunc("/user", user_add).Methods("PATCH")
	router.HandleFunc("/user", user_add).Methods("DELETE")	
	//
	router.HandleFunc("/data", getData).Methods("GET")
	//	
	log.Fatal(http.ListenAndServe(":" + _port, router))
}