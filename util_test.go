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
	tacl := testingACL{t: t}
	mockFile("{\"key1\":\"value1\", \"key2\":\"value2\", \"intkey3\":3}", nil)
	result, err := readJSON("nothing")
	expected := make(map[string]interface{})
	expected["key1"] = "value1"
	expected["key2"] = "value2"
	expected["intkey3"] = (float64)(3) //json.Unmarshal return float64 for all JSON Numbers. https://pkg.go.dev/encoding/json#Unmarshal
	assert.Nil(tacl, err)
	assert.Equal(tacl, expected, result)
}

func Test_readJSON_FileReadError(t *testing.T) {
	tacl := testingACL{t: t}
	mockFile("{\"key1\":\"value1\", \"key2\":\"value2\", \"intkey3\":3}", errors.New("File reading error"))
	result, err := readJSON("nothing")
	assert.EqualError(tacl, err, "File reading error")
	assert.Nil(tacl, result)
}

func Test_readJSON_UnmarshalError(t *testing.T) {
	tacl := testingACL{t: t}
	mockFile("{\"key1\":\"value1\", \"key2\":, \"intkey3\":3}", nil) // key2 should have a value
	result, err := readJSON("nothing")
	expected := make(map[string]interface{})
	expected["key1"] = "value1"
	expected["key2"] = "value2"
	expected["intkey3"] = (float64)(3) //json.Unmarshal return float64 for all JSON Numbers. https://pkg.go.dev/encoding/json#Unmarshal
	assert.EqualError(tacl, err, "invalid character ',' looking for beginning of value")
	assert.Nil(tacl, result)
}

func Test_readYaml_Success(t *testing.T) {
	tacl := testingACL{t: t}
	mockFile("key1: \"value1\"	\nkey2: \"value2\"  \nintkey3: 3", nil)
	result, err := readYaml("nothing")
	expected := make(map[string]interface{})
	expected["key1"] = "value1"
	expected["key2"] = "value2"
	expected["intkey3"] = 3
	assert.Nil(tacl, err)
	assert.Equal(tacl, expected, result)
}

func Test_readYaml_FileReadError(t *testing.T) {
	tacl := testingACL{t: t}
	mockFile("key1: \"value1\"	\nkey2: \"value2\"  \nintkey3: 3", errors.New("File reading error"))
	result, err := readYaml("nothing")
	assert.EqualError(tacl, err, "File reading error")
	assert.Nil(tacl, result)
}

func Test_readYaml_UnmarshalError(t *testing.T) {
	tacl := testingACL{t: t}
	mockFile("key1: \"value1\"	key2: \"value2\"  \nintkey3: 3", nil) // key2 should be on a new line
	result, err := readYaml("nothing")
	expected := make(map[string]interface{})
	expected["key1"] = "value1"
	expected["key2"] = "value2"
	expected["intkey3"] = (float64)(3)
	assert.EqualError(tacl, err, "yaml: did not find expected key")
	assert.Nil(tacl, result)
}
