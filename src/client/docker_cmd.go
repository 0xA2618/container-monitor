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
			fmt.Println("runServer:", fields)
			lenCount := len(fields)
			name := fields[lenCount-1]
			if name == monitorServerName {
				return fields[lenCount-2]
			}
		}
	}
	return "stop"
}
