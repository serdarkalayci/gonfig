package gonfig

import (
	"errors"
	"os"
	"strings"
)

type loadedSource struct {
	items  map[string]interface{}
	source ConfigSource
	err    error
}

// Configuration is the collection of loaded configuration sources
type Configuration struct {
	sources  []loadedSource
	HasError bool
}

// AddConfigSource adds multiple configuration sources to the collection.
// Config sources will be evaluated in the order they are added.
func (c Configuration) AddConfigSource(s ConfigSource) Configuration {
	newSource := loadedSource{
		source: s,
	}
	switch s.Type {
	case "json":
		newSource.items, newSource.err = readJSON(s.FilePath)
	case "yaml":
		newSource.items, newSource.err = readYaml(s.FilePath)
	}
	if newSource.err != nil {
		c.HasError = true
	}
	c.sources = append(c.sources, newSource)
	return c
}

func (c Configuration) findKey(key string) (interface{}, bool) {
	var value interface{}
	var found bool
	for _, loadedSource := range c.sources {
		switch loadedSource.source.Type {
		case "json", "yaml":
			if val, fnd := loadedSource.items[key]; fnd {
				value = val
				found = fnd
			}
		case "env":
			if val, fnd := os.LookupEnv(key); fnd {
				value = val
				if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") { // We will assume the returned val is an array if it starts with "[" and ends with "]"
					val = strings.TrimPrefix(val, "[")
					val = strings.TrimSuffix(val, "]")
					varr := strings.Split(val, ",")
					value = varr
				}
				found = fnd
			}
		}
	}
	return value, found
}

// GetInt returns the int value if the key is amongst the config sources and if the value is convertable to int
// Returns an error otherwise
func (c Configuration) GetInt(key string) (int, error) {
	val, found := c.findKey(key)
	if !found {
		return 0, errors.New("The key is not found among config sources")
	}
	return convertToInt(val)
}

// GetString returns the string value if the key is amongst the config sources
// Returns an error otherwise
func (c Configuration) GetString(key string) (string, error) {
	val, found := c.findKey(key)
	if !found {
		return "", errors.New("The key is not found among config sources")
	}
	return convertToString(val), nil
}

// GetFloat returns the float value if the key is amongst the config sources and if the value is convertable to float
// Returns an error otherwise
func (c Configuration) GetFloat(key string) (float64, error) {
	val, found := c.findKey(key)
	if !found {
		return 0, errors.New("The key is not found among config sources")
	}
	return convertToFloat(val)
}

// GetIntArray returns the []int value if the key is amongst the config sources.
// It'll ignore if the items cannot be converted to int by skipping them
// Returns an error otherwise
func (c Configuration) GetIntArray(key string) ([]int, error) {
	val, found := c.findKey(key)
	if !found {
		return nil, errors.New("The key is not found among config sources")
	}
	arr := make([]int, 0)
	switch val := val.(type) {
	case []string:
		for _, value := range val {
			if newval, err := convertToInt(value); err == nil {
				arr = append(arr, newval)
			}
		}
	case []interface{}:
		for _, value := range val {
			if newval, err := convertToInt(value); err == nil {
				arr = append(arr, newval)
			}
		}
	default:
		return nil, errors.New("The value is not an array or slice")
	}
	return arr, nil
}

// GetStringArray returns the []string value if the key is amongst the config sources.
// Returns an error otherwise
func (c Configuration) GetStringArray(key string) ([]string, error) {
	val, found := c.findKey(key)
	if !found {
		return nil, errors.New("The key is not found among config sources")
	}
	arr := make([]string, 0)
	switch val := val.(type) {
	case []string:
		arr = val
	case []interface{}:
		for _, value := range val {
			arr = append(arr, convertToString(value))
		}
	default:
		return nil, errors.New("The value is not an array or slice")
	}
	return arr, nil
}
