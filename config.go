package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type LogConfig struct {
	LogSetLevel int `json:"logSetLevel"`
}

type NodeConfig struct {
	ListenEndPoint string `json:"listen"`
	MaxConn        uint32 `json:"maxConn"`
	TimeoutRead    uint32 `json:"timeoutRead"`
}

type TargetConfig struct {
	ConnEndPoint string `json:"endpoint"`
	MaxConn      uint32 `json:"maxConn"`
	TimeoutConn  uint32 `json:"timeoutConn"`
	TimeoutRead  uint32 `json:"timeoutRead"`
}

type Config struct {
	ConfigFileName string         `json:"-"`
	Threads        uint32         `json:"threads"`
	Log            LogConfig      `json:"log"`
	Node           NodeConfig     `json:"node"`
	Targets        []TargetConfig `json:"targets"`
}

func LoadConfig(cfg *Config) {
	configFileName := "config.json"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}
	configFileName, _ = filepath.Abs(configFileName)
	log.Printf("Load Config: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("Open File error: ", err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&cfg); err != nil {
		log.Fatal("Load Json Config error: ", err.Error())
	}

	cfg.ConfigFileName = configFileName
}

func SaveConfig(cfg *Config) {
	log.Printf("Save Config: %v", cfg.ConfigFileName)

	configFile, err := os.Open(cfg.ConfigFileName)
	if err != nil {
		log.Fatal("Open File error: ", err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewEncoder(configFile)
	if err := jsonParser.Encode(&cfg); err != nil {
		log.Fatal("Save Json Config error: ", err.Error())
	}
}
