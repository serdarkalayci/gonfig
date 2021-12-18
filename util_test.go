package gonfig

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ReadFile func(filename string) ([]byte, error)

func mockFile(payload string, err error) {
	myReadFile = func(filename string) ([]byte, error) {
		return []byte(payload), err
	}
}

func Test_readJSON_Success(t *testing.T) {
	mockFile("{\"key1\":\"value1\", \"key2\":\"value2\", \"intkey3\":3}", nil)
	result, err := readJSON("nothing")
	expected := make(map[string]interface{})
	expected["key1"] = "value1"
	expected["key2"] = "value2"
	expected["intkey3"] = (float64)(3) //json.Unmarshal return float64 for all JSON Numbers. https://pkg.go.dev/encoding/json#Unmarshal
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func Test_readJSON_FileReadError(t *testing.T) {
	mockFile("{\"key1\":\"value1\", \"key2\":\"value2\", \"intkey3\":3}", errors.New("File reading error"))
	result, err := readJSON("nothing")
	assert.EqualError(t, err, "File reading error")
	assert.Nil(t, result)
}

func Test_readJSON_UnmarshalError(t *testing.T) {
	mockFile("{\"key1\":\"value1\", \"key2\":, \"intkey3\":3}", nil) // key2 should have a value
	result, err := readJSON("nothing")
	expected := make(map[string]interface{})
	expected["key1"] = "value1"
	expected["key2"] = "value2"
	expected["intkey3"] = (float64)(3) //json.Unmarshal return float64 for all JSON Numbers. https://pkg.go.dev/encoding/json#Unmarshal
	assert.EqualError(t, err, "invalid character ',' looking for beginning of value")
	assert.Nil(t, result)
}

func Test_readYaml_Success(t *testing.T) {
	mockFile("key1: \"value1\"	\nkey2: \"value2\"  \nintkey3: 3", nil)
	result, err := readYaml("nothing")
	expected := make(map[string]interface{})
	expected["key1"] = "value1"
	expected["key2"] = "value2"
	expected["intkey3"] = 3
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func Test_readYaml_FileReadError(t *testing.T) {
	mockFile("key1: \"value1\"	\nkey2: \"value2\"  \nintkey3: 3", errors.New("File reading error"))
	result, err := readYaml("nothing")
	assert.EqualError(t, err, "File reading error")
	assert.Nil(t, result)
}

func Test_readYaml_UnmarshalError(t *testing.T) {
	mockFile("key1: \"value1\"	key2: \"value2\"  \nintkey3: 3", nil) // key2 should be on a new line
	result, err := readYaml("nothing")
	expected := make(map[string]interface{})
	expected["key1"] = "value1"
	expected["key2"] = "value2"
	expected["intkey3"] = (float64)(3)
	assert.EqualError(t, err, "yaml: did not find expected key")
	assert.Nil(t, result)
}
