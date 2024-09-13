package settings

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type KeyperConfig struct {
	InstanceID               int      `env:"_ASSETS_INSTANCE_ID"`
	MaxNumKeysPerMessage      int      `env:"_ASSETS_MAX_NUM_KEYS_PER_MESSAGE"`
	EncryptedGasLimit         int      `env:"_ASSETS_ENCRYPTED_GAS_LIMIT"`
	GenesisSlotTimestamp      int      `env:"_ASSETS_GENESIS_SLOT_TIMESTAMP"`
	SyncStartBlockNumber      int      `env:"_ASSETS_SYNC_START_BLOCK_NUMBER"`
	KeyperSetManager          string   `env:"_ASSETS_KEYPER_SET_MANAGER"`
	KeyBroadcastContract      string   `env:"_ASSETS_KEY_BROADCAST_CONTRACT"`
	Sequencer                 string   `env:"_ASSETS_SEQUENCER"`
	ValidatorRegistry         string   `env:"_ASSETS_VALIDATOR_REGISTRY"`
	DiscoveryNamespace        string   `env:"_ASSETS_DISCOVERY_NAME_PREFIX"`
	CustomBootstrapAddresses  []string `env:"_ASSETS_CUSTOM_BOOTSTRAP_ADDRESSES"`
	DKGPhaseLength            int      `env:"_ASSETS_DKG_PHASE_LENGTH"`
	DKGStartBlockDelta        int      `env:"_ASSETS_DKG_START_BLOCK_DELTA"`

	DatabaseURL               string   `env:"SHUTTER_DATABASE_URL"`
	BeaconAPIURL              string   `env:"SHUTTER_BEACONAPIURL"`
	ContractsURL              string   `env:"SHUTTER_GNOSIS_NODE_CONTRACTSURL"`
	MaxTxPointerAge           int      `env:"SHUTTER_GNOSIS_MAXTXPOINTERAGE"`
	DeploymentDir             string   `env:"SHUTTER_DEPLOYMENT_DIR"`  // Unused, but you can still add an env if needed
	EthereumURL               string   `env:"SHUTTER_GNOSIS_NODE_ETHEREUMURL"`
	ShuttermintURL            string   `env:"SHUTTER_SHUTTERMINT_SHUTTERMINTURL"`
	ListenAddresses           string   `env:"SHUTTER_P2P_LISTENADDRESSES"`
	AdvertiseAddresses        string   `env:"SHUTTER_P2P_ADVERTISEADDRESSES"`
	ValidatorPublicKey        string   `env:"VALIDATOR_PUBLIC_KEY"`
	Enabled                   bool     `env:"SHUTTER_ENABLED"`
}

// AddSettingsToKeyper modifies the keyper settings by combining the template, config, and environment variables.
func AddSettingsToKeyper(templateFilePath, configFilePath, outputFilePath string) error {
	var keyperConfig KeyperConfig

	fmt.Println("Adding user settings to keyper...")

	// Read the keyper config file
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("error reading chain config TOML file: %v", err)
	}

	// Unmarshal the chain config TOML into the chainConfig struct
	err = toml.Unmarshal(configFile, &keyperConfig)
	if err != nil {
		return fmt.Errorf("error unmarshalling chain config TOML file: %v", err)
	}

	// Read the template file
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
	applyKeyperConfig(&templateConfig, keyperConfig)
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

	fmt.Println("Keyper TOML file modified successfully and saved to", outputFilePath)
	return nil
}

// applyKeyperConfig modifies the template based on environment variables, config, and default template values
func applyKeyperConfig(templateConfig *map[string]interface{}, keyperConfig KeyperConfig) {
	v := reflect.ValueOf(keyperConfig)
	t := reflect.TypeOf(keyperConfig)

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		// Get the environment variable tag
		envVar := fieldType.Tag.Get("env")

		// Get the corresponding field value based on the priority (env > config > template)
		fieldName := fieldType.Name

		var finalValue interface{}

		// 1. Check environment variable
		envValue := os.Getenv(envVar)
		if envValue != "" {
			finalValue = parseEnvValue(envValue, fieldType.Type)
		} else if !isZeroValue(fieldValue.Interface()) {
			// 2. Use config value if it exists and is non-zero
			finalValue = fieldValue.Interface()
		} else {
			// 3. Fallback to template value (do nothing as it's already in template)
			finalValue = (*templateConfig)[fieldName]
		}

		setTemplateField(templateConfig, fieldName, finalValue)
	}
}

// Helper function to update template configuration
func setTemplateField(templateConfig *map[string]interface{}, key string, value interface{}) {
	if !isZeroValue(value) {
		(*templateConfig)[key] = value
	}
}

// Parse environment variable value into the correct type
func parseEnvValue(value string, valueType reflect.Type) interface{} {
	switch valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, _ := strconv.ParseInt(value, 10, 64)
		return intValue
	case reflect.Bool:
		boolValue, _ := strconv.ParseBool(value)
		return boolValue
	case reflect.Slice:
		// Assume comma-separated values for slice
		return strings.Split(value, ",")
	default:
		return value
	}
}

// Check if a value is the zero value of its type
func isZeroValue(value interface{}) bool {
	return reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
}
