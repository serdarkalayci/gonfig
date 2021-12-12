package gonfig

// SourceType describes the source type of configuration to be read
type SourceType string

const (
	// Env sourcetype is used for reading from environment settings
	Env SourceType = "env"
	// Flag sourcetype is used for reading from command flags
	Flag = "flag"
	// JSON sourcetype is used for reading from JSON formatted files
	JSON = "json"
	// Yaml sourcetype is used for reading from yaml formatted files
	Yaml = "yaml"
)

// ConfigSource is the type that is used to describe various config sources.
type ConfigSource struct {
	// Type is the type SourceType of the ConfigSource
	Type SourceType
	// FilePath is the absolute path to the file if the Type is SourceType.JSON or SourceType.Yaml
	FilePath string
}
