package main

import "log"

type Configuration struct {
	Address string
	Static string
	ReadTimeout int64
	WriteTimeout int64
}

var config Configuration
var logger *log.Logger

func init() {

}

func loadConfig() {

}
