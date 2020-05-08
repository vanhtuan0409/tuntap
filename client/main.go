package main

import (
	"log"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

const (
	ifName = "tun99"
)

func runIP(args ...string) error {
	cmd := exec.Command("/usr/bin/ip", args...)
	return cmd.Run()
}

func senderProxy(ifce *water.Interface) {
	buf := make([]byte, 1500)
	for {
		n, err := ifce.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		buf = buf[:n]

		conn, err := net.Dial("tcp", "localhost:6868")
		if err != nil {
			log.Printf("Cannot connect to VPN server. ERR: %v\n", err)
			continue
		}
		conn.Write(buf)
		conn.Close()
		log.Println("Sent data")
	}
}

func main() {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = ifName

	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer ifce.Close()

	if err := runIP("addr", "add", "10.1.0.1/24", "dev", ifName); err != nil {
		log.Fatal(err)
	}
	if err := runIP("link", "set", "dev", ifName, "up"); err != nil {
		log.Fatal(err)
	}
	// if err := runIP("route", "add", "default", "dev", ifName); err != nil {
	// 	log.Fatal(err)
	// }

	senderProxy(ifce)
}
