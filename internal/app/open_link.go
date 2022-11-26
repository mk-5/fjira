package app

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func OpenLink(url string) {
	err := openLinkForSystem(runtime.GOOS, url)
	if err != nil {
		log.Fatal(err)
	}
}

func openLinkForSystem(system string, url string) error {
	switch system {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
