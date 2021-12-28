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
	val, err := c.GetFloat("key1")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, val)
	val, err = c.GetFloat("key2")
	assert.EqualError(t, err, "The key is not found among config sources")
	assert.Equal(t, 0.0, val)
}
