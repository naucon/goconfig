# Go Config

[![Build](https://github.com/naucon/goconfig/actions/workflows/go-ci.yml/badge.svg)](https://github.com/naucon/goconfig/actions/workflows/go-ci.yml)
[![Coverage](https://codecov.io/gh/naucon/goconfig/branch/master/graph/badge.svg?token=3R985BKFKB)](https://codecov.io/gh/naucon/goconfig)

This package is a lightweight configuration management for go projects.
It manages yaml configurations for different environments like prod, stage, test.
Credentials and secrets should be defined with placeholders like `${APP_DATABASE_DSN}` in the yaml file.
The placeholders will be replaced with OS ENV variables.

Additionally, the package supports populating OS ENV variables with dot ENV files by using the `joho/godotenv` package.

The following dot ENV files are supported

`.env` defines the default values of the env variables that are required. Typically, this file will be committed to git (like a `.env.dist` file).
`.env.local` overrides the default values for all environments, if not testing environment. This file MUST NOT be committed to git and should be added to a `.gitignore` file.
`.env.<environment>` (e.g. `.env.test` or `.env.prod`) overrides again the values with environment specific values. Typically, these files are committed to git.
`.env.<environment>.local` (e.g. `.env.prod.local`): overwrites again the environment specific values. This file MUST NOT be committed to git and should also be added to a `.gitignore` file.


**NOTICE:** Never commit credentials and secrets or `.env` files that contain credentials and secrets to your git repository.

## Requires

* Go 1.13 or newer
* gopkg.in/yaml.v
* github.com/joho/godotenv

## Installation

install the latest version via go get

```
go get -u github.com/naucon/goconfig
```

## Import package

```
import (
  "github.com/naucon/goconfig"
)
```

## Usage

### Configuration struct

First create a struct that defines the desired configuration entries.
Also, a nested struct can be used to group configuration entries.
Some config entries may require `yaml tags` to define a map to the yaml properties.
More information how to define `yaml tags` in your struct can be found on the [go-yaml docs](https://pkg.go.dev/gopkg.in/yaml.v3#Unmarshal).

```go
type MyConfiguration struct {
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
```

### create Yaml files

Next create yaml files that contain the values of the configuration entries. The file name must match the app environment eg. `prod.yml`.
The directory with your config yaml files should look like this:

```
dev.yml
prod.yml
test.yml
```

Credentials and secrets in your yaml file should be defined with placeholders. The placeholders will be replaced with os ENV variables.
Placeholder have the following structure `${var}`. `var` must be the os ENV variable name.
Internally we're using [os.ExpandEnv](https://pkg.go.dev/os#ExpandEnv) function to replace placeholders with os ENV variables.

```yaml
database:
  dsn: "${APP_DATABASE_DSN}"
server:
  port: 3000
  trusted_proxies:
    - "127.0.0.1"
debug: false
secret: "${APP_SECRET}"
```

### create dot env files (optional)

optionally you can define a dot files. They will be used by the package to define os ENV variables.

```
APP_DATABASE_DSN="postgres://user:pwd@localhost:5432/database_name"
APP_SECRET="..."
```

Dot env files do not overwrite the OS ENV variables.

`.env` defines the default values of the env variables that are required. Typically, this file will be committed to git (like a `.env.dist` file).
`.env.local` overrides the default values for all environments, if not testing environment. This file MUST NOT be committed to git and should be added to a `.gitignore` file.
`.env.<environment>` (e.g. `.env.test` or `.env.prod`) overrides again the values with environment specific values. Typically, these files are committed to git.
`.env.<environment>.local` (e.g. `.env.prod.local`): overwrites again the environment specific values. This file MUST NOT be committed to git and should also be added to a `.gitignore` file.

The following pattern should be added to a `.gitignore` file.

```
.env.local
.env.*.local
```

### Loading configuration

Finally, to load the configuration we create an instance with `goconfig.NewConfig()`. Pass in `goconfig.Options`.
Then populate your configuration struct by calling `Load()` from the config instance and pass in app environment eg. `prod` and your struct (with a pointer).

```go
		cfg := MyConfiguration{}
		c := goconfig.NewConfig(goconfig.Options{})
		err := c.Load("dev", &cfg)
```

#### config options

The behavior of the config loader can be changed with the `Options` struct that is passed on creation.
The structs contains the following options:

```go
type Options struct {
	DotEnvPath          string // path where the dot env files are located, by default "./"
	ConfigPath          string // path where the yaml config files are located, by default "./config/"
	ConfigFileExtension string // file extension of the yaml config files, by default "yml"
	TestEnv             []string // defines the testing environments, by default "test"
	Verbose             bool   // verbose mode logs debug messages, by default false
}
```

```go
c := goconfig.NewConfig(goconfig.Options{
  DotEnvPath: "./",
  ConfigPath: "./config/"
  ConfigFileExtension: "yml",
  Verbose: false,
})
err := c.Load("dev", &cfg)
```

## Example

```go
package main

import (
	"fmt"
  "github.com/naucon/goconfig"
)

type MyConfiguration struct {
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

func main() {
  cfg := MyConfiguration{}
  c := goconfig.NewConfig(goconfig.Options{})
  if err := c.Load("dev", &cfg); err != nil {
    panic(err)
  }

  initDatabase(cfg.Database.Dsn)
  initServer(cfg.Server)
  if cfg.Debug {
    fmt.Println("Debug mode on")
  }
}
```
