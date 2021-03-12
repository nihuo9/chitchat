package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Address string
	Static string
	ReadTimeout int64
	WriteTimeout int64
}

var config Configuration
var logger *log.Logger

func init() {
	loadConfig()
	file, err := os.OpenFile("chitchat.log", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Cannot open log file", err)
	}
	logger = log.New(file, "INFO", log.Ldate | log.Ltime | log.Lshortfile)
}

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("Cannot get configuration from file", err)
	}
}
