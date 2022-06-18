package extractor

import (
	"fmt"
	"runtime"
	"testing"
)

func TestExtractFromLinux(t *testing.T) {

	if runtime.GOOS == "linux" {

		networkInterfaceInfo := NetworkInterfaceInformation{
			Name: "eth0",
		}

		results := ExtractFromLinux(networkInterfaceInfo)
		fmt.Printf("bytes received: %d, bytes sent: %d", results.BytesReceived, results.BytesSent)
	}
}

func TestExtractFromWindows(t *testing.T) {

	if runtime.GOOS == "windows" {

		networkInterfaceInfo := NetworkInterfaceInformation{
			Name: "Realtek RTL8822BE 802.11ac PCIe Adapter",
		}

		results := ExtractFromWindows(networkInterfaceInfo)
		fmt.Printf("bytes received: %d, bytes sent: %d", results.BytesReceived, results.BytesSent)
	}
}
