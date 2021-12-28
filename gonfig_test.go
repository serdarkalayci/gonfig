package gonfig

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddConfigSource_JSON(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type:     JSON,
		FilePath: "testing.json",
	}
	mockFile("{\"key1\":\"value1\", \"key2\":\"value2\", \"intkey3\":3}", nil)
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	val, found := c.findKey("key1")
	assert.Equal(t, true, found)
	assert.Equal(t, "value1", val)
	// Let's add a second YAML source with different values
	mockFile("key1: \"value1_2\"	\nkey2: \"value2_2\"  \nintkey3: 32", nil)
	s = ConfigSource{
		Type:     Yaml,
		FilePath: "testing.yaml",
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 2)
	val, found = c.findKey("key1")
	assert.Equal(t, true, found)
	assert.Equal(t, "value1_2", val)
	// Let's add a second JSON source which is erroneous. This will not change the value of keys because the config source would not be added as a loaded source
	mockFile("{\"key1\":\"value1_2\", \"key2\":\"value2_2\", \"intkey3\":32}", errors.New("File reading error"))
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 3)
	val, found = c.findKey("key1")
	assert.Equal(t, true, found)
	assert.Equal(t, "value1_2", val)
	assert.EqualError(t, c.sources[2].err, "File reading error") // Let's make sure that the relevant config source has the error returned
}

func Test_AddConfigSource_Env(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: Env,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	os.Setenv("key1", "value1")
	val, found := c.findKey("key1")
	assert.Equal(t, true, found)
	assert.Equal(t, "value1", val)
	// Let's try to find
	val, found = c.findKey("key2")
	assert.Equal(t, false, found)
	assert.Nil(t, val)
}

func Test_GetInt(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: Env,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	os.Setenv("key1", "1")
	val, err := c.GetInt("key1")
	assert.Nil(t, err)
	assert.Equal(t, 1, val)
	val, err = c.GetInt("key2")
	assert.EqualError(t, err, "The key is not found among config sources")
	assert.Equal(t, 0, val)
}

func Test_convertToInt(t *testing.T) {
	// int
	var intVal int
	intVal = 1
	val, err := convertToInt(intVal)
	assert.Equal(t, 1, val)
	assert.Nil(t, err)
	// string
	var strVal string
	strVal = "12"
	val, err = convertToInt(strVal)
	assert.Equal(t, 12, val)
	assert.Nil(t, err)
	// float32
	var fl32Val float32
	fl32Val = 21.0
	val, err = convertToInt(fl32Val)
	assert.Equal(t, 21, val)
	assert.Nil(t, err)
	// float64
	var fl64Val float64
	fl64Val = 21.0
	val, err = convertToInt(fl64Val)
	assert.Equal(t, 21, val)
	assert.Nil(t, err)
	// bool
	var boolVal bool
	boolVal = true
	val, err = convertToInt(boolVal)
	assert.Equal(t, 1, val)
	assert.Nil(t, err)
	boolVal = false
	val, err = convertToInt(boolVal)
	assert.Equal(t, 0, val)
	assert.Nil(t, err)
	// int8
	var int8Val int8
	int8Val = 8
	val, err = convertToInt(int8Val)
	assert.Equal(t, 8, val)
	assert.Nil(t, err)
	// int16
	var int16Val int16
	int16Val = 16
	val, err = convertToInt(int16Val)
	assert.Equal(t, 16, val)
	assert.Nil(t, err)
	// int32
	var int32Val int32
	int32Val = 32
	val, err = convertToInt(int32Val)
	assert.Equal(t, 32, val)
	assert.Nil(t, err)
	// int64
	var int64Val int64
	int64Val = 64
	val, err = convertToInt(int64Val)
	assert.Equal(t, 64, val)
	assert.Nil(t, err)
	// uint
	var uintVal uint
	uintVal = 1
	val, err = convertToInt(uintVal)
	assert.Equal(t, 1, val)
	assert.Nil(t, err)
	// uint8
	var uint8Val uint8
	uint8Val = 8
	val, err = convertToInt(uint8Val)
	assert.Equal(t, 8, val)
	assert.Nil(t, err)
	// uint16
	var uint16Val uint16
	uint16Val = 16
	val, err = convertToInt(uint16Val)
	assert.Equal(t, 16, val)
	assert.Nil(t, err)
	// uint32
	var uint32Val uint32
	uint32Val = 32
	val, err = convertToInt(uint32Val)
	assert.Equal(t, 32, val)
	assert.Nil(t, err)
	// uint64
	var uint64Val uint64
	uint64Val = 64
	val, err = convertToInt(uint64Val)
	assert.Equal(t, 64, val)
	assert.Nil(t, err)
	// unknown
	var unknown = make(map[string]string)
	unknown["test"] = "test"
	val, err = convertToInt(unknown)
	assert.Equal(t, 0, val)
	assert.EqualError(t, err, "Unknown type")
}

func Test_GetString(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: Env,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	os.Setenv("key1", "1")
	val, err := c.GetString("key1")
	assert.Nil(t, err)
	assert.Equal(t, "1", val)
	val, err = c.GetString("key2")
	assert.EqualError(t, err, "The key is not found among config sources")
	assert.Equal(t, 0, val)
}

func Test_convertToString(t *testing.T) {
	// int
	var intVal int
	intVal = 1
	val := convertToString(intVal)
	assert.Equal(t, "1", val)
	// string
	var strVal string
	strVal = "12"
	val = convertToString(strVal)
	assert.Equal(t, "12", val)
	// float32
	var fl32Val float32
	fl32Val = 21.0
	val = convertToString(fl32Val)
	assert.Equal(t, "21", val)
	// float64
	var fl64Val float64
	fl64Val = 21.0
	val = convertToString(fl64Val)
	assert.Equal(t, "21", val)
	// bool
	var boolVal bool
	boolVal = true
	val = convertToString(boolVal)
	assert.Equal(t, "true", val)
	boolVal = false
	val = convertToString(boolVal)
	assert.Equal(t, "false", val)
	// int8
	var int8Val int8
	int8Val = 8
	val = convertToString(int8Val)
	assert.Equal(t, "8", val)
	// int16
	var int16Val int16
	int16Val = 16
	val = convertToString(int16Val)
	assert.Equal(t, "16", val)
	// int32
	var int32Val int32
	int32Val = 32
	val = convertToString(int32Val)
	assert.Equal(t, "32", val)
	// int64
	var int64Val int64
	int64Val = 64
	val = convertToString(int64Val)
	assert.Equal(t, "64", val)
	// uint
	var uintVal uint
	uintVal = 1
	val = convertToString(uintVal)
	assert.Equal(t, "1", val)
	// uint8
	var uint8Val uint8
	uint8Val = 8
	val = convertToString(uint8Val)
	assert.Equal(t, "8", val)
	// uint16
	var uint16Val uint16
	uint16Val = 16
	val = convertToString(uint16Val)
	assert.Equal(t, "16", val)
	// uint32
	var uint32Val uint32
	uint32Val = 32
	val = convertToString(uint32Val)
	assert.Equal(t, "32", val)
	// uint64
	var uint64Val uint64
	uint64Val = 64
	val = convertToString(uint64Val)
	assert.Equal(t, "64", val)
	// unknown
	var unknown = make(map[string]string)
	unknown["test"] = "test"
	val = convertToString(unknown)
	assert.Equal(t, "map[test:test]", val)
}
