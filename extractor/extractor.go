package extractor

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// NetworkInterfaceInformation
// Currently, we only need to know the name for your interface used
// by the extraction methods.
// For linux: you get this from the command, "ifconfig".
// For windows: you get this from the command, "ipconfig /all".
type NetworkInterfaceInformation struct {
	Name string
}

// ExtractedResults
type ExtractedResults struct {
	BytesSent     int64
	BytesReceived int64
}

// ExtractFromLinux uses the files found in /sys/class/net/.
func ExtractFromLinux(iface NetworkInterfaceInformation) *ExtractedResults {

	path := fmt.Sprintf("/sys/class/net/%s/statistics", iface.Name)

	return &ExtractedResults{
		BytesReceived: readNumberFromFile(path + "/rx_bytes"),
		BytesSent:     readNumberFromFile(path + "/tx_bytes"),
	}
}

// ExtractFromWindows uses wmic.
// This is variably suboptimal. Needs improvement.
func ExtractFromWindows(iface NetworkInterfaceInformation) *ExtractedResults {

	cmd := exec.Command("wmic", "path", "Win32_PerfRawData_Tcpip_NetworkInterface", "get", "BytesReceivedPersec,", "BytesSentPersec,", "Name")
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(stdout), "\n")
	if len(lines) <= 1 {
		log.Fatal("no network interface detected.")
	}
	var bytesSent int64 = 0
	var bytesReceived int64 = 0

	for i := 1; i < len(lines); i++ {

		temp := strings.TrimSuffix(lines[i], "\n")
		temp = strings.TrimSpace(temp)
		if strings.HasSuffix(temp, iface.Name) {
			columns := produceColumnsFromRow(temp)
			bytesReceived = parseInt64FromString(columns[0])
			bytesSent = parseInt64FromString(columns[1])
		}
	}

	return &ExtractedResults{
		BytesSent:     bytesSent,
		BytesReceived: bytesReceived,
	}
}

func readNumberFromFile(path string) int64 {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	output, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	output = strings.TrimSuffix(output, "\n")
	number := parseInt64FromString(output)
	return number
}

func produceColumnsFromRow(row string) []string {

	columns := make([]string, 0)
	var accumulator string = ""
	for _, ch := range row {
		if string(ch) == " " {
			if accumulator == "" {
				continue
			}
			columns = append(columns, accumulator)
			accumulator = ""
			continue
		}
		accumulator += string(ch)
	}
	columns = append(columns, accumulator)
	return columns
}

func parseInt64FromString(input string) int64 {

	number, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return number
}
