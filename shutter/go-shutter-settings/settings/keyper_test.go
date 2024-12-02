package settings

import (
	"testing"

	"github.com/joho/godotenv"
	"gotest.tools/v3/assert"
)

func TestUnmarshal(t *testing.T) {
	err := godotenv.Load("./test-assets/variables.env")
	assert.NilError(t, err, "assets variables error")
	var generatedConfig map[string]interface{}
	err = UnmarshallFromFile("./test-assets/keyper.toml", &generatedConfig)
	assert.NilError(t, err, "error unmarshalling")
}

func TestAddSettings(t *testing.T) {
	err := godotenv.Load("./test-assets/variables.env")
	assert.NilError(t, err, "assets variables error")
	err = AddSettingsToKeyper("./test-assets/generated.toml", "./test-assets/keyper.toml", "./test-assets/out.toml")
	assert.NilError(t, err, "error")
}
