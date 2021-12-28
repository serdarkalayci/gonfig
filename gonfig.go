package gonfig

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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
				found = fnd
			}
		}
	}
	return value, found
}

func (c Configuration) GetInt(key string) (int, error) {
	val, found := c.findKey(key)
	if !found {
		return 0, errors.New("The key is not found among config sources")
	}
	return convertToInt(val)
}

func convertToInt(val interface{}) (int, error) {
	var i int // your final value
	var err error
	switch t := val.(type) {
	case int:
		i = t
	case int8:
		i = int(t) // standardizes across systems
	case int16:
		i = int(t) // standardizes across systems
	case int32:
		i = int(t) // standardizes across systems
	case int64:
		i = int(t) // standardizes across systems
	case bool:
		if t {
			i = 1
		} else {
			i = 0
		}
	case float32:
		i = int(t) // standardizes across systems
	case float64:
		i = int(t) // standardizes across systems
	case uint:
		i = int(t) // standardizes across systems
	case uint8:
		i = int(t) // standardizes across systems
	case uint16:
		i = int(t) // standardizes across systems
	case uint32:
		i = int(t) // standardizes across systems
	case uint64:
		i = int(t) // standardizes across systems
	case string:
		i, err = strconv.Atoi(t)
	default:
		i = 0
		err = errors.New("Unknown type")
	}
	return i, err
}

func (c Configuration) GetString(key string) (string, error) {
	val, found := c.findKey(key)
	if !found {
		return "", errors.New("The key is not found among config sources")
	}
	return convertToString(val), nil
}

func convertToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}
