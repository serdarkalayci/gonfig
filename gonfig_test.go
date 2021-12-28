package gonfig

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddConfigSource_JSON_YAML(t *testing.T) {
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
	mockFile("key1: \"value1_2\"\nkey2: \"value2_2\"\nintkey3: 32", nil)
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

func Test_AddConfigSource_JSON_YAML_Array(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type:     JSON,
		FilePath: "testing.json",
	}
	mockFile("{\"key1\":\"value1\", \"strarraykey\":[\"strval1\",\"strval2\"], \"intarraykey\":[123,456,789]}", nil)
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	val, found := c.findKey("strarraykey")
	assert.Equal(t, true, found)
	assert.Equal(t, []interface{}{"strval1", "strval2"}, val)
	val, found = c.findKey("intarraykey")
	assert.Equal(t, true, found)
	assert.Equal(t, []interface{}{123.0, 456.0, 789.0}, val)
	// Let's add a second YAML source with different values
	mockFile("key1: \"value1_2\"\nstrarraykey:\n  - strval1_2\n  - strval2_2\nintarraykey:\n  - 321\n  - 654\n  - 987", nil)
	s = ConfigSource{
		Type:     Yaml,
		FilePath: "testing.yaml",
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 2)
	val, found = c.findKey("strarraykey")
	assert.Equal(t, true, found)
	assert.Equal(t, []interface{}{"strval1_2", "strval2_2"}, val)
	val, found = c.findKey("intarraykey")
	assert.Equal(t, true, found)
	assert.Equal(t, []interface{}{321, 654, 987}, val)
}

func Test_AddConfigSource_Env(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: Env,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	os.Setenv("key1", "value1")
	// Let's try to find a key that's there
	val, found := c.findKey("key1")
	assert.Equal(t, true, found)
	assert.Equal(t, "value1", val)
	// Let's try to find a key that's not there
	val, found = c.findKey("key2")
	assert.Equal(t, false, found)
	assert.Nil(t, val)
	os.Setenv("arrkey1", "[\"val1\",\"val2\",\"val3\"]")
	// Let's try to find an array key that's there
	val, found = c.findKey("arrkey1")
	assert.Equal(t, true, found)
	assert.Equal(t, []string{"\"val1\"", "\"val2\"", "\"val3\""}, val)
}

func Test_GetInt(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: Env,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	os.Setenv("key1", "1")
	// Let's try to find a key that's there
	val, err := c.GetInt("key1")
	assert.Nil(t, err)
	assert.Equal(t, 1, val)
	// Let's try to find an array key that's not there
	val, err = c.GetInt("key2")
	assert.EqualError(t, err, "The key is not found among config sources")
	assert.Equal(t, 0, val)
}

func Test_GetString(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: Env,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	os.Setenv("key1", "1")
	// Let's try to find a key that's there
	val, err := c.GetString("key1")
	assert.Nil(t, err)
	assert.Equal(t, "1", val)
	// Let's try to find an array key that's not there
	val, err = c.GetString("key2")
	assert.EqualError(t, err, "The key is not found among config sources")
	assert.Equal(t, "", val)
}

func Test_GetFloat(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: Env,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	os.Setenv("key1", "1")
	// Let's try to find a key that's there
	val, err := c.GetFloat("key1")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, val)
	// Let's try to find an array key that's not there
	val, err = c.GetFloat("key2")
	assert.EqualError(t, err, "The key is not found among config sources")
	assert.Equal(t, 0.0, val)
}

func Test_GetIntArray(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: Env,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	// Env array variables are interpreted as []string
	os.Setenv("arrkey1", "[123,456,\"abc\",789]")
	// Let's try to find an array key that's there
	val, err := c.GetIntArray("arrkey1")
	assert.Nil(t, err)
	assert.Equal(t, []int{123, 456, 789}, val)
	// JSON and Yaml array variables are interpreted as []interface{}
	mockFile("{\"key1\":\"value1\", \"strarraykey\":[\"strval1\",\"strval2\"], \"intarraykey\":[123,456,789]}", nil)
	s = ConfigSource{
		Type:     JSON,
		FilePath: "testing.json",
	}
	c = c.AddConfigSource(s)
	// Let's try to find an array key that's there
	val, err = c.GetIntArray("intarraykey")
	assert.Nil(t, err)
	assert.Equal(t, []int{123, 456, 789}, val)
	// Let's try to find a value that's not an array
	val, err = c.GetIntArray("key1")
	assert.EqualError(t, err, "The value is not an array or slice")
	assert.Equal(t, []int(nil), val)
	// Let's try to find an array key that's not there
	val, err = c.GetIntArray("key2")
	assert.EqualError(t, err, "The key is not found among config sources")
	assert.Equal(t, []int(nil), val)
}
