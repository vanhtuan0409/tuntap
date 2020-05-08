package main

import (
	"io"
	"log"
	"net"
	"os/exec"

	"github.com/songgao/water"
)

const (
	ifName = "tun98"
)

func runIP(args ...string) error {
	cmd := exec.Command("/usr/bin/ip", args...)
	return cmd.Run()
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

	if err := runIP("addr", "add", "10.2.0.1/24", "dev", ifName); err != nil {
		log.Fatal(err)
	}
	if err := runIP("link", "set", "dev", ifName, "up"); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp", ":6868")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Running server on :6868")
	for {
		conn, err := l.Accept()
		if err != nil {
			break
		}
		log.Println("Received connection")
		go io.Copy(ifce, conn)
	}
}
