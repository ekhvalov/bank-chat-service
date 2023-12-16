package main

var shutdownFunctions = make([]func() error, 0)

func registerShutdownFunc(f func() error) {
	shutdownFunctions = append(shutdownFunctions, f)
}
