package client

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func DockerPS() []string {

	// command to execute
	cmd := exec.Command("docker", "ps")

	// output
	output, err := cmd.Output()
	if err != nil {
		log.Println(err.Error())
		return []string{}
	}

	// Convert the output to string and split it into lines
	return strings.Split(string(output), "\n")

}

func CheckServetActive(monitorServerName string, lines []string) string {

	// Skip the header line (first line) using NR > 1
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line != "" {
			fields := strings.Fields(line)

			lenCount := len(fields)
			name := fields[lenCount-1]
			serverPort := fields[lenCount-2]
			fmt.Println("serverName:", name)
			fmt.Println("serverPort", serverPort)
			if name == monitorServerName {

				return serverPort
			}
		}
	}
	return "stop"
}
