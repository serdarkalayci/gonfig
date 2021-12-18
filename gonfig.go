package gonfig

import "os"

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
			if val, found := loadedSource.items[key]; found {
				value = val
			}
		case "env":
			if val, found := os.LookupEnv(key); found {
				value = val
			}
		}
	}
	return value, found
}
