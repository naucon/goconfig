package examples

type Configuration struct {
	Database Database
	Server   Server
	Debug    bool
	Secret   string
}

type Database struct {
	Dsn string
}

type Server struct {
	Hostname       string
	Port           int
	ReadTimeout    int
	WriteTimeout   int
	TrustedProxies []string `yaml:"trusted_proxies"`
}
