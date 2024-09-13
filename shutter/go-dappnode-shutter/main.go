package main

import (
	"flag"
	"fmt"
	"go-dappnode-shutter/settings"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Define flags for the template, config, and output paths
	templateFilePath := flag.String("template", "", "Path to the template file where settings will be included")
	configFilePath := flag.String("config", "", "Path to the config file where the settings will be read")
	outputFilePath := flag.String("output", "", "Path where the modified settings will be saved")

	// Parse the flags
	flag.Parse()

	// Load environment variables from the .env file
	err := godotenv.Load(os.Getenv("ASSETS_DIR") + "/variables.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Check for additional arguments, e.g., keyper or chain
	if len(flag.Args()) < 1 {
		fmt.Println("Error: missing argument. Use 'include-keyper-settings' or 'include-chain-settings'.")
		os.Exit(1)
	}

	// Read the argument passed to the program
	argument := flag.Arg(0)

	// Call appropriate function based on the command
	switch argument {
	case "include-keyper-settings":
		// Ensure template, config, and output paths are provided
		if *templateFilePath == "" || *configFilePath == "" || *outputFilePath == "" {
			fmt.Println("Error: --template, --config, and --output flags must be provided for keyper settings.")
			flag.Usage()
			os.Exit(1)
		}

		// Call the function to configure keyper
		err := settings.AddSettingsToKeyper(*templateFilePath, *configFilePath, *outputFilePath)
		if err != nil {
			log.Fatalf("Failed to configure keyper: %v", err)
		}

	case "include-chain-settings":
		// Ensure config and output paths are provided
		if *configFilePath == "" || *outputFilePath == "" {
			fmt.Println("Error: --config and --output flags must be provided for chain settings.")
			flag.Usage()
			os.Exit(1)
		}

		// Call the function to configure chain
		err := settings.AddSettingsToChain(*configFilePath, *outputFilePath)
		if err != nil {
			log.Fatalf("Failed to configure chain: %v", err)
		}

	default:
		fmt.Println("Invalid argument. Use 'include-keyper-settings' or 'include-chain-settings'.")
		os.Exit(1)
	}

	fmt.Println("Configuration completed successfully!")
}
