package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
)

func main() {
	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		log.Fatalf("Failed to discover gateway: %v", err)
	}

	client := natpmp.NewClient(gatewayIP)
	response, err := client.GetExternalAddress()
	if err != nil {
		return
	}
	fmt.Printf("External IP address: %v\n", response.ExternalIPAddress)

	pm, err := client.AddPortMapping("tcp", 8088, 80, 10*60)
	if err != nil {
		log.Fatalf("Failed to add port mapping: %v", err)
	}
	fmt.Printf("Mapped external port %v to internal port %v for %v seconds\n", pm.MappedExternalPort, pm.InternalPort, pm.PortMappingLifetimeInSeconds)

	fmt.Printf("To test, run listener:\n")
	fmt.Printf("nc -l %v\n", pm.InternalPort)

	ipv4 := net.IPv4(response.ExternalIPAddress[0], response.ExternalIPAddress[1], response.ExternalIPAddress[2], response.ExternalIPAddress[3])

	fmt.Printf("Connect to the listener from outside NAT:\n")
	fmt.Printf("nc %v %v\n", ipv4, pm.MappedExternalPort)

}
