//go:build windows

package main

import (
	"log"

	"golang.getoutline.org/sdk/x/sysproxy"
)

func setSystemProxy(address string, port string) error {
	log.Printf("Setting system proxy to %s:%s\n", address, port)
	return sysproxy.SetWebProxy(address, port)
}

func unsetSystemProxy() error {
	log.Println("Unsetting system proxy")
	return sysproxy.DisableWebProxy()
}
