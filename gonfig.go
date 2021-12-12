package gonfig

type source struct {
	items  map[string]interface{}
	source ConfigSource
	err    error
}

type Configuration struct {
	sources  []source
	HasError bool
}

func (c Configuration) AddConfigSource(s ConfigSource) Configuration {
	newSource := source{
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
