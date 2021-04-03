package main

import (
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {
	var packetCounter uint

	ServerAddr, err := net.ResolveUDPAddr("udp", "[::1]:6343")
	checkError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "[::1]:0")
	checkError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	checkError(err)

	defer Conn.Close()
	fileName := "./dump.pcap"

	if handle, err := pcap.OpenOffline(fileName); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			appLayer := packet.Layer(layers.LayerTypeUDP)
			//fmt.Println(appLayer.LayerPayload())
			//fmt.Printf(packet.String())
			_, err := Conn.Write(appLayer.LayerPayload())
			if err != nil {
				fmt.Println(err)
				continue
			}
			packetCounter++
		}
	}
	fmt.Println("Packet prosessed:", packetCounter)
}
