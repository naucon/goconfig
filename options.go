package goconfig

type Options struct {
	DotEnvPath          string
	ConfigPath          string
	ConfigFileExtension string
	TestEnv 			[]string
	Verbose             bool
}
