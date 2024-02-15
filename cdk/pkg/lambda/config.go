package lambda

type Config struct {
	BuildDir        string
	BuildOutputPath string
	GolangPath      string
	Handler         string
	FunctionName    string
}

func NewConfig(buildDir string, golangPath string, functionName string) Config {
	return Config{
		BuildDir:        buildDir,
		BuildOutputPath: buildDir + "/" + "bootstrap",
		GolangPath:      golangPath,
		Handler:         "main",
		FunctionName:    functionName,
	}
}
