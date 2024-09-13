package settings

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type ChainConfig struct {
    PrivateKey string
    InstanceID string
}

/**
* This function should:
* - Read all the values from the template configuration file
* - Read the values defined in ChainConfig from the chain configuration file
* - Read the environment variables
* - Copy all the values from the template configuration file to the output configuration file
* - Modify the values using the chain configuration file and environment variables (these have more priority)
* - Save the modified configuration file
* - Return an error if any
*/
func AddSettingsToChain(templateFilePath, outputFilePath string) error {

    var chainConfig ChainConfig
    
    chainConfigPath := os.Getenv("SHUTTER_CHAIN_CONFIG_FILE")

    // Check that the chain configuration file exists
    if _, err := os.Stat(chainConfigPath); os.IsNotExist(err) {
        fmt.Println("Chain configuration file does not exist:", chainConfigPath)
        os.Exit(1)
    }


    chainConfigFile, err := os.ReadFile(chainConfigPath)
    if err != nil {
        fmt.Println("Error reading TOML file:", err)
        os.Exit(1)
    }

    err = toml.Unmarshal(chainConfigFile, &chainConfig)
    if err != nil {
        fmt.Println("Error unmarshalling TOML file:", err)
        os.Exit(1)
    }

    // Read the template configuration file
	templateFile, err := os.ReadFile(templateFilePath)
	if err != nil {
		return fmt.Errorf("error reading template TOML file: %v", err)
	}

	// Create a map to hold the template configuration
	var templateConfig map[string]interface{}
	err = toml.Unmarshal(templateFile, &templateConfig)
	if err != nil {
		return fmt.Errorf("error unmarshalling template TOML file: %v", err)
	}

    	// Modify the template configuration based on ChainConfig and environment variables
	// ChainConfig values take priority over the template values
	applyChainConfig(&templateConfig, chainConfig)

	// Apply environment variables, which have even higher priority than the ChainConfig
	applyEnvOverrides(&templateConfig)

	// Marshal the modified configuration to TOML format
	modifiedConfig, err := toml.Marshal(templateConfig)
	if err != nil {
		return fmt.Errorf("error marshalling modified config to TOML: %v", err)
	}

	// Write the modified configuration to the output file
	err = os.WriteFile(outputFilePath, modifiedConfig, 0644)
	if err != nil {
		return fmt.Errorf("error writing modified TOML file: %v", err)
	}

	fmt.Println("TOML file modified successfully and saved to", outputFilePath)
	return nil
}

// applyChainConfig modifies the template configuration based on the values from ChainConfig
func applyChainConfig(templateConfig *map[string]interface{}, chainConfig ChainConfig) {
	if chainConfig.PrivateKey != "" {
		(*templateConfig)["PrivateKey"] = chainConfig.PrivateKey
	}
	if chainConfig.InstanceID != "" {
		(*templateConfig)["InstanceID"] = chainConfig.InstanceID
	}
}

// applyEnvOverrides applies environment variables to the template configuration, giving them the highest priority
func applyEnvOverrides(templateConfig *map[string]interface{}) {
	// Example: Check for environment variables and override values if they exist
	if privateKey := os.Getenv("PRIVATE_KEY"); privateKey != "" {
		(*templateConfig)["PrivateKey"] = privateKey
	}
	if instanceID := os.Getenv("INSTANCE_ID"); instanceID != "" {
		(*templateConfig)["InstanceID"] = instanceID
	}
}