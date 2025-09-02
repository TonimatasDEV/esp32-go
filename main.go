package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func findEsp32Port() (string, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return "", err
	}

	for _, port := range ports {
		if port.IsUSB {
			if strings.Contains(strings.ToLower(port.VID), "1a86") || strings.Contains(strings.ToLower(port.VID), "0403") {
				fmt.Printf("Detected ESP32 in %s (%s - %s)\n", port.Name, port.VID, port.PID)
				return port.Name, nil
			}
		}
	}
	return "", fmt.Errorf("ESP32 not found in USB")
}

func main() {
	portName, err := findEsp32Port()
	if err != nil {
		log.Fatal(err)
	}

	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	fmt.Printf("Connected to the port: %s\n", portName)

	for {
		cpuPercent, _ := cpu.Percent(0, false)
		vmStat, _ := mem.VirtualMemory()
		diskStat, _ := disk.Usage("C:/")

		line := fmt.Sprintf("CPU:%d%% MEM:%d%% DISK:%d%%\n",
			int(cpuPercent[0]),
			int(vmStat.UsedPercent),
			int(diskStat.UsedPercent),
		)

		_, err := port.Write([]byte(line))
		if err != nil {
			log.Printf("Error sent: %v", err)
		}

		time.Sleep(1 * time.Second)
	}
}
