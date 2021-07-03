package main

import (
    "bytes"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "sort"
    "strings"

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

type DHCPLease struct {
  ExpirationTime   int64  `yaml:"ExpirationTime"`
  MAC              string `yaml:"MAC"`
  IP               string `yaml:"IP"`
  Hostname         string `yaml:"Hostname"`
  ClientIdentifier string `yaml:"ClientIdentifier"`
}

func main() {
    pingerFilePtr := flag.String("ping","ping.json","Output file from monitor-ping")
    wificlientsFilePtr := flag.String("wifi","wifi.yml","Output file from monitor-wifi")
    dhcpFilePtr := flag.String("dhcp","dhcp.yml","Output file from monitor-dhcp")
    flag.Parse()

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

    b, err = ioutil.ReadFile(*dhcpFilePtr)
    if err != nil {
        log.Fatal(err)
    }
    var dhcp []DHCPLease
    yaml.Unmarshal(b, &dhcp)

    // Combined database
    hostState := make(map[string]HostState)

    // Fill data from ping
    for _, host := range network.Hosts {
        if host.Answer {
		hostState[host.Ip] = HostState{RttMs: host.RttMs, Hostname: host.Name}
        }
    }
    // Fill data from WiFi Access Points
    for _, client := range clients {
	var host HostState
        if len(client.IP) == 0 {
            for _, dhcpData := range dhcp {
                if client.MAC == dhcpData.MAC {
                    client.IP = dhcpData.IP
                    client.Hostname = dhcpData.Hostname
                }
            }
        }
        // If there is already an entry from ping, get its data out
        host, _ = hostState[client.IP]
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
	upstreamState.WiFi = client.WiFi
	host.UpstreamState = append(host.UpstreamState, upstreamState)

	// store back
	hostState[client.IP] = host
    }
    // Slice with all IP addresses
    IPs := make([]string, 0, len(hostState))
    for k, _ := range hostState {
        IPs = append(IPs, k)
    }
    // sort slice IPs according to the actual IP value
    sort.Slice(
        IPs,
        func(i, j int) bool {
            return bytes.Compare(
              net.ParseIP(IPs[i]), net.ParseIP(IPs[j]))<0
        })
    // loop over the sorted IPs
    for _, k := range IPs {
        v := hostState[k]
	PingLabel := "↺"
        if v.RttMs == "" {
            PingLabel = ""
        }
        fmt.Printf("%-20s %-15s %-17s %s%s\n", v.Hostname, k, v.MAC, PingLabel, v.RttMs)
	for _, u := range v.UpstreamState {
		fmt.Printf(" @ %s from %s ↑%v ↓%v, %vGHz %vdBm %vMbps\n",
		               u.Upstream,
			                strings.Replace(u.ActiveTime, " days", "d", 1),
					            u.Up,
						                   u.Down,
								       map[int]string{0:"2.4", 1:"5"}[u.WiFi.Radio],
								             u.WiFi.RSSI,
                                                                                   u.WiFi.Rate)
	}
    }
}
