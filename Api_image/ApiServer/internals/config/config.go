package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

func LoadGatewayConfig(filename string) *GatewayConfig {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil
	}

	config := &GatewayConfig{
		GatewayHost: "localhost",
		GatewayPort: 8000,
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Printf("Error unmarshalling config: %v", err)
		return nil
	}

	return config
}

func LoadAnalyticsConfig(filename string) *AnalyticsConfig {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil
	}

	config := &AnalyticsConfig{
		AnalyticsHost: "localhost",
		AnalyticsPort: 8000,
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Printf("Error unmarshalling config: %v", err)
		return nil
	}

	return config
}

func LoadResourceConfig(filename string) *ResourceConfig {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil
	}

	config := &ResourceConfig{
		ResourceHost: "localhost",
		ResourcePort: 8000,
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Printf("Error unmarshalling config: %v", err)
		return nil
	}

	return config
}

func LoadConnectorConfig(filename string) *ConnectorConfig {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil
	}

	config := &ConnectorConfig{
		ConnectorHost: "localhost",
		ConnectorPort: 8000,
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Printf("Error unmarshalling config: %v", err)
		return nil
	}

	return config
}

func LoadDBConfig(filename string) *DBConfig {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading config file %s", filename)
		return nil
	}

	config := &DBConfig{}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Printf("Error reading config file %s", filename)
		return nil
	}

	return config
}
