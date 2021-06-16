package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "log"

    "gopkg.in/yaml.v2"
)

type Host struct {
    Ip string `yaml:"ip"`
    Answer bool `yaml:"answer"`
    RttMs string `yaml:"rtt"`
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
    Ping       Host 
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

    c := make(map[string]NetworkClient)
    for i := 0; i<len(network.Hosts); i++ {
        if network.Hosts[i].Answer {
            c[network.Hosts[i].Ip] = NetworkClient{Ping: Host{RttMs: network.Hosts[i].RttMs}}
//            fmt.Printf("%20s %15s\n", " ", network.Hosts[i].Ip)
        }
    }
    for i := 0 ; i<len(clients) ; i++ {
        if client, ok := c[clients[i].IP] ; ok {
            c[clients[i].IP].Hostname = clients[i].Hostname
            c[clients[i].IP].MAC = clients[i].MAC
        } else {
            c[clients[i].IP] = clients[i]
        }
        // fmt.Printf("%-20s %-15s %-17s\n", clients[i].Hostname, clients[i].IP, clients[i].MAC)
    }
    for k, v := range c {
        fmt.Printf("%-20s %-15s %-17s %s\n", v.Hostname, k, v.MAC, v.Ping.RttMs)
    }
}
