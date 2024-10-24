package gonfig

// SourceType describes the source type of configuration to be read
type SourceType string

const (
	// SourceTypeEnv is used for reading from environment settings
	SourceTypeEnv SourceType = "env"
	// SourceTypeJSON is used for reading from JSON formatted files
	SourceTypeJSON = "json"
	// SourceTypeYaml is used for reading from yaml formatted files
	SourceTypeYaml = "yaml"
)

// ConfigSource is the type that is used to describe various config sources.
type ConfigSource struct {
	// Type is the type SourceType of the ConfigSource
	Type SourceType
	// FilePath is the absolute path to the file if the Type is SourceType.JSON or SourceType.Yaml
	FilePath string
}
