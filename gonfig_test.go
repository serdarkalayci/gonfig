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
		Type: SourceTypeJSON,
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
		Type: SourceTypeYaml,
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
		Type: SourceTypeJSON,
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
		Type: SourceTypeYaml,
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
		Type: SourceTypeEnv,
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
		Type: SourceTypeEnv,
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

func Test_GetIntOrDefault(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: SourceTypeEnv,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	os.Setenv("key1", "1")
	// Let's try to find a key that's there
	val := c.GetIntOrDefault("key1", 25)
	assert.Equal(t, 1, val)
	os.Setenv("key3", "{-}")
	// Let's try to find a key that's there but not convertible
	val = c.GetIntOrDefault("key3", 27)
	assert.Equal(t, 27, val)
}

func Test_GetString(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: SourceTypeEnv,
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

func Test_GetStringOrDefault(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: SourceTypeEnv,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	os.Setenv("key1", "1")
	// Let's try to find a key that's there
	val := c.GetStringOrDefault("key1", "default value 1")
	assert.Equal(t, "1", val)
	// Let's try to find an array key that's not there
	val = c.GetStringOrDefault("key2", "default value 2")
	assert.Equal(t, "default value 2", val)
}

func Test_GetFloat(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: SourceTypeEnv,
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

func Test_GetFloatOrDefault(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: SourceTypeEnv,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	os.Setenv("key1", "1")
	// Let's try to find a key that's there
	val := c.GetFloatOrDefault("key1", 12.7)
	assert.Equal(t, 1.0, val)
	os.Setenv("key3", "{-}")
	// Let's try to find a key that's there but now convertible
	val = c.GetFloatOrDefault("key3", 12.9)
	assert.Equal(t, 12.9, val)
}

func Test_GetBool(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: SourceTypeEnv,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	os.Setenv("key1", "1")
	// Let's try to find a key that's there
	val, err := c.GetBool("key1")
	assert.Nil(t, err)
	assert.Equal(t, true, val)
	// Let's try to find an array key that's not there
	val, err = c.GetBool("key2")
	assert.EqualError(t, err, "The key is not found among config sources")
	assert.Equal(t, false, val)
}

func Test_GetBoolOrDefault(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: SourceTypeEnv,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	os.Setenv("key1", "1")
	// Let's try to find a key that's there
	val := c.GetBoolOrDefault("key1", false)
	assert.Equal(t, true, val)
	// Let's try to find an array key that's not there
	val = c.GetBoolOrDefault("key2", true)
	assert.Equal(t, true, val)
	// Let's try to find a key that's there
	os.Setenv("key3", "{-}")
	val = c.GetBoolOrDefault("key3", true)
	assert.Equal(t, true, val)
}

func Test_GetIntArray(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: SourceTypeEnv,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	// Env array variables are interpreted as []string
	os.Setenv("intarraykey", "[123,456,\"abc\",789]")
	// Let's try to find an array key that's there
	val, err := c.GetIntArray("intarraykey")
	assert.Nil(t, err)
	assert.Equal(t, []int{123, 456, 789}, val)
	// JSON and Yaml array variables are interpreted as []interface{}
	mockFile("{\"key1\":\"value1\", \"strarraykey\":[\"strval1\",\"strval2\"], \"intarraykey\":[321,654,987]}", nil)
	s = ConfigSource{
		Type: SourceTypeJSON,
		FilePath: "testing.json",
	}
	c = c.AddConfigSource(s)
	// Let's try to find an array key that's there
	val, err = c.GetIntArray("intarraykey")
	assert.Nil(t, err)
	assert.Equal(t, []int{321, 654, 987}, val)
	// Let's try to find a value that's not an array
	val, err = c.GetIntArray("key1")
	assert.EqualError(t, err, "The value is not an array or slice")
	assert.Equal(t, []int(nil), val)
	// Let's try to find an array key that's not there
	val, err = c.GetIntArray("key2")
	assert.EqualError(t, err, "The key is not found among config sources")
	assert.Equal(t, []int(nil), val)
}

func Test_GetIntArrayOrDefault(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: SourceTypeEnv,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	defArray := make([]int, 2)
	defArray[0] = 999
	defArray[1] = 888
	// Env array variables are interpreted as []string
	os.Setenv("intarraykey", "[123,456,\"abc\",789]")
	// Let's try to find an array key that's there
	val := c.GetIntArrayOrDefault("intarraykey", defArray)
	assert.Equal(t, []int{123, 456, 789}, val)
	// JSON and Yaml array variables are interpreted as []interface{}
	mockFile("{\"key1\":\"value1\", \"strarraykey\":[\"strval1\",\"strval2\"], \"intarraykey\":[321,654,987]}", nil)
	s = ConfigSource{
		Type: SourceTypeJSON,
		FilePath: "testing.json",
	}
	c = c.AddConfigSource(s)
	// Let's try to find an array key that's there
	val = c.GetIntArrayOrDefault("intarraykey", defArray)
	assert.Equal(t, []int{321, 654, 987}, val)
	// Let's try to find a value that's not an array
	val = c.GetIntArrayOrDefault("key1", defArray)
	assert.Equal(t, defArray, val)
	// Let's try to find an array key that's not there
	val = c.GetIntArrayOrDefault("key2", defArray)
	assert.Equal(t, defArray, val)
}

func Test_GetStringArray(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: SourceTypeEnv,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	// Env array variables are interpreted as []string
	os.Setenv("strarraykey", "[strval1,strval2,123]")
	// Let's try to find an array key that's there
	val, err := c.GetStringArray("strarraykey")
	assert.Nil(t, err)
	assert.Equal(t, []string{"strval1", "strval2", "123"}, val)
	// JSON and Yaml array variables are interpreted as []interface{}
	mockFile("{\"key1\":\"value1\", \"strarraykey\":[\"strval1_2\",\"strval2_2\"], \"intarraykey\":[123,456,789]}", nil)
	s = ConfigSource{
		Type: SourceTypeJSON,
		FilePath: "testing.json",
	}
	c = c.AddConfigSource(s)
	// Let's try to find an array key that's there
	val, err = c.GetStringArray("strarraykey")
	assert.Nil(t, err)
	assert.Equal(t, []string{"strval1_2", "strval2_2"}, val)
	// Let's try to find a value that's not an array
	val, err = c.GetStringArray("key1")
	assert.EqualError(t, err, "The value is not an array or slice")
	assert.Equal(t, []string(nil), val)
	// Let's try to find an array key that's not there
	val, err = c.GetStringArray("key2")
	assert.EqualError(t, err, "The key is not found among config sources")
	assert.Equal(t, []string(nil), val)
}

func Test_GetStringArrayOrDefault(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: SourceTypeEnv,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	defArray := make([]string, 2)
	defArray[0] = "999"
	defArray[1] = "888"
	// Env array variables are interpreted as []string
	os.Setenv("strarraykey", "[strval1,strval2,123]")
	// Let's try to find an array key that's there
	val := c.GetStringArrayOrDefault("strarraykey", defArray)
	assert.Equal(t, []string{"strval1", "strval2", "123"}, val)
	// JSON and Yaml array variables are interpreted as []interface{}
	mockFile("{\"key1\":\"value1\", \"strarraykey\":[\"strval1_2\",\"strval2_2\"], \"intarraykey\":[123,456,789]}", nil)
	s = ConfigSource{
		Type: SourceTypeJSON,
		FilePath: "testing.json",
	}
	c = c.AddConfigSource(s)
	// Let's try to find an array key that's there
	val = c.GetStringArrayOrDefault("strarraykey", defArray)
	assert.Equal(t, []string{"strval1_2", "strval2_2"}, val)
	// Let's try to find a value that's not an array
	val = c.GetStringArrayOrDefault("key1", defArray)
	assert.Equal(t, defArray, val)
	// Let's try to find an array key that's not there
	val = c.GetStringArrayOrDefault("key2", defArray)
	assert.Equal(t, defArray, val)
}

func Test_GetFloatArray(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: SourceTypeEnv,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	// Env array variables are interpreted as []string
	os.Setenv("fltarraykey", "[123.45,456.78,\"abc\",789.01]")
	// Let's try to find an array key that's there
	val, err := c.GetFloatArray("fltarraykey")
	assert.Nil(t, err)
	assert.Equal(t, []float64{123.45, 456.78, 789.01}, val)
	// JSON and Yaml array variables are interpreted as []interface{}
	mockFile("{\"key1\":\"value1\", \"strarraykey\":[\"strval1\",\"strval2\"], \"fltarraykey\":[123.12,456.45,789.78]}", nil)
	s = ConfigSource{
		Type: SourceTypeJSON,
		FilePath: "testing.json",
	}
	c = c.AddConfigSource(s)
	// Let's try to find an array key that's there
	val, err = c.GetFloatArray("fltarraykey")
	assert.Nil(t, err)
	assert.Equal(t, []float64{123.12, 456.45, 789.78}, val)
	// Let's try to find a value that's not an array
	val, err = c.GetFloatArray("key1")
	assert.EqualError(t, err, "The value is not an array or slice")
	assert.Equal(t, []float64(nil), val)
	// Let's try to find an array key that's not there
	val, err = c.GetFloatArray("key2")
	assert.EqualError(t, err, "The key is not found among config sources")
	assert.Equal(t, []float64(nil), val)
}

func Test_GetFloatArrayOrDefault(t *testing.T) {
	var c Configuration
	s := ConfigSource{
		Type: SourceTypeEnv,
	}
	c = c.AddConfigSource(s)
	assert.Equal(t, len(c.sources), 1)
	defArray := make([]float64, 2)
	defArray[0] = 999.9
	defArray[1] = 888.8
	// Env array variables are interpreted as []string
	os.Setenv("fltarraykey", "[123.45,456.78,\"abc\",789.01]")
	// Let's try to find an array key that's there
	val := c.GetFloatArrayOrDefault("fltarraykey", defArray)
	assert.Equal(t, []float64{123.45, 456.78, 789.01}, val)
	// JSON and Yaml array variables are interpreted as []interface{}
	mockFile("{\"key1\":\"value1\", \"strarraykey\":[\"strval1\",\"strval2\"], \"fltarraykey\":[123.12,456.45,789.78]}", nil)
	s = ConfigSource{
		Type: SourceTypeJSON,
		FilePath: "testing.json",
	}
	c = c.AddConfigSource(s)
	// Let's try to find an array key that's there
	val = c.GetFloatArrayOrDefault("fltarraykey", defArray)
	assert.Equal(t, []float64{123.12, 456.45, 789.78}, val)
	// Let's try to find a value that's not an array
	val = c.GetFloatArrayOrDefault("key1", defArray)
	assert.Equal(t, defArray, val)
	// Let's try to find an array key that's not there
	val = c.GetFloatArrayOrDefault("key2", defArray)
	assert.Equal(t, defArray, val)
}
