package config

type Runtime interface {
	Version() string
	Rich() bool
}

// Version TODO: generate that automatically
const (
	Version = "dev-0.0.1"
)

type runtime struct {
	version string
	rich    bool
}

func NewRuntime(version string, cliConfig CliConfig) Runtime {
	return &runtime{
		version: version,
		rich:    cliConfig.Rich,
	}
}

func (d runtime) Version() string {
	return d.version
}

func (d runtime) Rich() bool {
	return d.rich
}
