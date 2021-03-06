package main

// Mostly based on https://github.com/golang/net/blob/master/icmp/ping_test.go
// All ye beware, there be dragons below...

import (
	"context"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/Ullaakut/nmap"
	log "github.com/sirupsen/logrus"
)

const (
	ProtocolICMP = 1
)

// Default to listen on all IPv4 interfaces
var ListenAddr = "0.0.0.0"

func runARP() {
	log.Debug("### Running ARP ###")
	data, err := exec.Command("arp", "-a").Output()
	if err != nil {
		PublishFailureNetwork(getTime(), "Arp failed")
		log.Error(err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// strip brackets around IP
		ip := strings.Replace(fields[1], "(", "", -1)
		ip = strings.Replace(ip, ")", "", -1)
		new_device := true
		mac := fields[3]
		for id := range DevicesList {
			if DevicesList[id].Ip_address == ip {
				new_device = false
				log.Trace("Device found in Arp table")
				DevicesList[id].Alive = true
				if DevicesList[id].Allowed != DISCOVERED {
					DevicesList[id].New = false
				}
			}
		}
		if new_device {
			if mac != "<incomplete>" {
				log.Warn("Adding device ip: ", ip)
				response, err := http.Get("https://api.macvendors.com/" + mac)

				defer response.Body.Close()

				data, _ := ioutil.ReadAll(response.Body)

				if err != nil {
					log.Error("The HTTP request failed with error \n", err)
					PublishFailureNetwork(getTime(), "Api call failed")
				} else {
					log.Trace(response)
					log.Debug("Vendor Name: ", string(data))
					device := Device{string(data), mac, ip, true, DISCOVERED, true}
					DevicesList[device_id] = &device
					device_id++
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
	for id := range DevicesList {
		if DevicesList[id].New && DevicesList[id].Allowed != DISCOVERED {
			log.Trace("Device found in Arp table")
			DevicesList[id].Alive = false
			DevicesList[id].New = false
		}
	}
}

func nmap_scan() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Equivalent to `/usr/local/bin/nmap -p 80,443,843 google.com facebook.com youtube.com`,
	// with a 2 minute timeout.
	scanner, err := nmap.NewScanner(
		nmap.WithTargets("192.168.0.0-255"),
		nmap.WithPorts("80,443,843"),
		nmap.WithContext(ctx),
	)
	if err != nil {
		log.Error("Unable to create nmap scanner: ", err)
	}

	result, warnings, err := scanner.Run()
	if err != nil {
		log.Error("Unable to run nmap scan: ", err)
	}

	if warnings != nil {
		log.Error("Warnings: ", warnings)
	}

	log.Debug("Nmap done: ", len(result.Hosts), " hosts up scanned in seconds ", result.Stats.Finished.Elapsed)
}

func stateDevices(id uint32) {
	if DevicesList[id].Allowed == DISCOVERED {
		PublishDeviceRequest(id,
			DevicesList[id].Device_name,
			DevicesList[id].Mac)
	}
}

func checkDevices() {
	for {
		nmap_scan()
		runARP()
		log.Warn("### Devices ###")
		for id := range DevicesList {
			log.Warn("Device - ", DevicesList[id].Device_name, " : ",
				DevicesList[id].Ip_address, " : ",
				DevicesList[id].Mac, " : ",
				DevicesList[id].Alive, " : ",
				DevicesList[id].Allowed, " : ",
				DevicesList[id].New)
			if DevicesList[id].Alive {
				stateDevices(id)
			}
		}
		log.Debug("### End of ARP ###")
		time.Sleep(4 * time.Minute)
	}
}
