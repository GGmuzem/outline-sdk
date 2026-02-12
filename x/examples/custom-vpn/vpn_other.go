//go:build !windows

package main

import "errors"

func setSystemProxy(address string, port string) error {
	return errors.New("system proxy not supported on this platform yet")
}

func unsetSystemProxy() error {
	return nil
}
