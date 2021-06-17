package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "log"

    "gopkg.in/yaml.v2"
)

type Host struct {
    Ip     string `yaml:"ip"`
    Name   string `yaml:"name"`
    Answer bool   `yaml:"answer"`
    RttMs  string `yaml:"rtt"`
}

type Network struct {
    Hosts []Host `yaml:"host"`
}

type WiFiParam struct {
    Radio      int    `yaml:"Radio"`
    RSSI       int    `yaml:"RSSI"`
    Rate       string `yaml:"Rate"`
}

type NetworkClient struct {
    Hostname   string `yaml:"hostname"`
    MAC        string `yaml:"MAC"`
    IP         string `yaml:"IP"`
    Down       int    `yaml:"Down"`
    Up         int    `yaml:"Up"`
    ActiveTime string `yaml:"ActiveTime"`
    LinkType   string `yaml:"linktype"`
    Upstream   string `yaml:"Upstream"`
    WiFi       WiFiParam `yaml:"WiFi,omitempty"`
}

type HostState struct {
	Hostname string 
	MAC string
	RttMs string
	UpstreamState []UpstreamState
}

type UpstreamState struct {
	Down       int
	Up         int
	ActiveTime string
	LinkType   string
	Upstream   string
	WiFi       WiFiParam
}

func main() {
    pingerFilePtr := flag.String("ping","ping.json","Output file from pinger")
    wificlientsFilePtr := flag.String("wifi","wifi.yml","Output file from wificlients")
    flag.Parse()
    fmt.Println("ping:",*pingerFilePtr)
    fmt.Println("wifi:",*wificlientsFilePtr)

    b, err := ioutil.ReadFile(*pingerFilePtr)
    if err != nil {
        log.Fatal(err)
    }
    var network Network
    yaml.Unmarshal(b, &network)

    b, err = ioutil.ReadFile(*wificlientsFilePtr)
    if err != nil {
        log.Fatal(err)
    }
    var clients []NetworkClient
    yaml.Unmarshal(b, &clients)

    // Combined database
    hostState := make(map[string]HostState)

    // Fill data from ping
    for _, host := range network.Hosts {
        if host.Answer {
		hostState[host.Ip] = HostState{RttMs: host.RttMs, Hostname: host.Name}
        }
    }
    // Fill data from Access Points
    for _, client := range clients {
	var host HostState
	var ok bool
        if host, ok = hostState[client.IP] ; ok {
            // do nothing
        } else {
            // do nothing
        }
	// Common data
        host.Hostname = client.Hostname
        host.MAC = client.MAC

	// Append data from UpstreamState
	var upstreamState UpstreamState
	upstreamState.Down = client.Down
	upstreamState.Up = client.Up
	upstreamState.ActiveTime = client.ActiveTime
	upstreamState.LinkType = client.LinkType
	upstreamState.Upstream = client.Upstream
	host.UpstreamState = append(host.UpstreamState, upstreamState)

	// store back
	hostState[client.IP] = host

        // fmt.Printf("%-20s %-15s %-17s\n", clients[i].Hostname, clients[i].IP, clients[i].MAC)
    }
    for k, v := range hostState {
        fmt.Printf("%-20s %-15s %s\n", v.Hostname, k, v.RttMs)
	for _, u := range v.UpstreamState {
		fmt.Printf(" @ %s since %s\n", u.Upstream, u.ActiveTime)
	}
    }
}
