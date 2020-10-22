package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type LogConfig struct {
	LogSetLevel int `json:"logSetLevel"`
}

type NodeConfig struct {
	ListenEndPoint string `json:"listen"`
	MaxConn        uint32 `json:"maxConn"`
	Timeout        uint32 `json:"timeout"`
}

type TargetConfig struct {
	ConnEndPoint string `json:"endpoint"`
	MaxConn      uint32 `json:"maxConn"`
	Timeout      uint32 `json:"timeout"`
}

type Config struct {
	ConfigFileName string         `json:"-"`
	Threads        uint32         `json:"threads"`
	Log            LogConfig      `json:"log"`
	Node           NodeConfig     `json:"node"`
	Targets        []TargetConfig `json:"targets"`
}

func loadConfig(cfg *Config) {
	configFileName := "config.json"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}
	configFileName, _ = filepath.Abs(configFileName)
	Info.Printf("Load Config: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		Error.Fatalf("Open File error: ", err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&cfg); err != nil {
		Error.Fatalf("Load Json Config error: ", err.Error())
	}

	cfg.ConfigFileName = configFileName
}

func saveConfig(cfg *Config) {
	Info.Printf("Save Config: %v", cfg.ConfigFileName)

	configFile, err := os.Open(cfg.ConfigFileName)
	if err != nil {
		Error.Fatalf("Open File error: ", err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewEncoder(configFile)
	if err := jsonParser.Encode(&cfg); err != nil {
		Error.Fatalf("Save Json Config error: ", err.Error())
	}
}
