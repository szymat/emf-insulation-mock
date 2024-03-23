package insulationmock

type config struct {
	ScriptsPath  string
	RoutesPath   string
	BeforeScript string
	AfterScript  string
	Port         string
}

// NewConfig creates a new instance of Config.
func NewConfig(scriptsPath, routesPath, beforeScript, afterScript, port string) *config {
	return &config{
		ScriptsPath:  scriptsPath,
		RoutesPath:   routesPath,
		BeforeScript: beforeScript,
		AfterScript:  afterScript,
		Port:         port,
	}
}
