package main

import (
	"encoding/json"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func convertStatus(status string) int {
	switch {
	case strings.Contains(status, ALLOWED_STRING):
		return ALLOWED
	case strings.Contains(status, BLOCKED_STRING):
		return BLOCKED
	case strings.Contains(status, UNKNOWN_STRING):
		return UNKNOWN
	default:
		return DISCOVERED
	}
}

func deviceResponse(request_id uint32) {
	if DevicesList[request_id].Allowed == BLOCKED {
		PublishDeviceFound(DevicesList[request_id].Device_name,
			DevicesList[request_id].Ip_address,
			DevicesList[request_id].Allowed)
	} else if DevicesList[request_id].Allowed == DISCOVERED {
		log.Error("We did not get a correct status")
	} else if DevicesList[request_id].Allowed == ALLOWED {
		log.Trace("Device is allowed")
	} else if DevicesList[request_id].Allowed == UNKNOWN {
		PublishDeviceFound(DevicesList[request_id].Device_name,
			DevicesList[request_id].Ip_address,
			DevicesList[request_id].Allowed)
	} else {
		log.Error("We shouldn't hit this error")
	}
}

func convertStatusMessage(message MapMessage) bool {
	switch {
	case message.routing_key == STATUSSYP:
		log.Debug(message.message)
		json.Unmarshal([]byte(message.message), &_statusSYP)
		log.Debug("Status for SYP")
		log.Debug("Highest Usage: " + strconv.Itoa(_statusSYP.HighestUsage))
		log.Debug("Temperature CPU: " + strconv.Itoa(_statusSYP.Temperature))
		log.Debug("CPU Memory Left: " + strconv.Itoa(_statusSYP.MemoryLeft))
		postHardware()
	case message.routing_key == STATUSUP:
		json.Unmarshal([]byte(message.message), &_statusUP)
		log.Debug("Status for UP")
		log.Debug("Last access blocked: " + _statusUP.LastAccessBlocked)
		log.Debug("Last access granted: " + _statusUP.LastAccessGranted)
		log.Debug("Last user: " + _statusUP.LastUser)
		postAccess()
	case message.routing_key == STATUSFH:
		json.Unmarshal([]byte(message.message), &_statusFH)
		log.Debug("Status for FH")
		log.Debug("Last Fault: " + _statusFH.LastFault)
		postFault()

	default:
		log.Warn("We received an incorrect status")
		return false
	}
	return true

}

func checkState() {
	for message_id := range SubscribedMessagesMap {
		if SubscribedMessagesMap[message_id].valid {
			switch {
			case strings.Contains(SubscribedMessagesMap[message_id].routing_key, STATUS):
				if convertStatusMessage(*SubscribedMessagesMap[message_id]) {
					SubscribedMessagesMap[message_id].valid = false
				}
			case SubscribedMessagesMap[message_id].routing_key == DEVICERESPONSE:
				log.Warn("Received a device response topic")
				var message DeviceResponse
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				DevicesList[message.Request_id].Allowed = convertStatus(message.Status)
				DevicesList[message.Request_id].Device_name = message.Name
				deviceResponse(message.Request_id)
				SubscribedMessagesMap[message_id].valid = false

			case SubscribedMessagesMap[message_id].routing_key == ALARMEVENT:
				log.Warn("Received a alarm event topic")
				var message AlarmEvent
				json.Unmarshal([]byte(SubscribedMessagesMap[message_id].message), &message)
				postAlarmEvent(message)
				SubscribedMessagesMap[message_id].valid = false

			default:
				log.Warn("We were not expecting this message unvalidating: ",
					SubscribedMessagesMap[message_id].routing_key)
				SubscribedMessagesMap[message_id].valid = false
			}
		}
	}

}
