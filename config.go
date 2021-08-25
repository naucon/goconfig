package goconfig

import (
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

type Config interface {
	Load(env string, out interface{}) error
}

type config struct {
	options Options
}

func NewConfig(o Options) *config {
	c := config{
		options: o,
	}
	c.optionDefaults()
	return &c
}

func (c *config) Load(env string, out interface{}) error {
	if env == "" {
		return NewConfigError(errConfigEmptyEnv, nil)
	}
	c.loadDotEnv(c.options.DotEnvPath, env)
	err := c.loadYamlConfiguration(c.options.ConfigPath+env+"."+c.options.ConfigFileExtension, out)
	if err != nil {
		return err
	}

	return nil
}

func (c *config) loadDotEnv(path string, env string) {
	var err error
	var filePath string

	if c.options.Verbose {
		log.Printf("loading dotenv from %s ...\n", path)
	}

	// loads eg. .env.prod.local
	filePath = path + ".env." + env + ".local"
	err = godotenv.Load(filePath)
	if c.options.Verbose && err == nil {
		log.Printf("dotenv loaded %s!\n", filePath)
	}

	// loads eg. .env.prod
	filePath = path + ".env." + env
	err = godotenv.Load(filePath)
	if c.options.Verbose && err == nil {
		log.Printf("dotenv loaded %s!\n", filePath)
	}

	// loads eg. .env.local
	if c.isTestEnv(env) == false {
		filePath = path + ".env.local"
		err = godotenv.Load(filePath)
		if c.options.Verbose && err == nil {
			log.Printf("dotenv loaded %s!\n", filePath)
		}
	}

	// loads default .env
	filePath = path + ".env"
	err = godotenv.Load(filePath)
	if c.options.Verbose && err == nil {
		log.Printf("dotenv loaded %s!\n", filePath)
	}
}

func (c *config) isTestEnv(env string) bool {
	for _, optTestEnv := range c.options.TestEnv {
		if optTestEnv == env {
			return true
		}
	}
	return false
}

func (c *config) loadYamlConfiguration(filePath string, out interface{}) error {
	if c.options.Verbose {
		log.Printf("loading yaml config %s ...\n", filePath)
	}
	bytContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return NewConfigError(errConfigMissing, err)
	}

	bytContent = []byte(os.ExpandEnv(string(bytContent)))
	if err := yaml.Unmarshal(bytContent, out); err != nil {
		return NewConfigError(errConfigInvalid, err)
	}
	if c.options.Verbose {
		log.Printf("yaml config loaded %s!\n", filePath)
	}

	return nil
}

func (c *config) optionDefaults() {
	if c.options.DotEnvPath == "" {
		c.options.DotEnvPath = "./"
	}
	if c.options.ConfigPath == "" {
		c.options.ConfigPath = "./config/"
	}
	if c.options.ConfigFileExtension == "" {
		c.options.ConfigFileExtension = "yml"
	}
	if len(c.options.TestEnv) == 0 {
		c.options.TestEnv = append(c.options.TestEnv, "test")
	}
}
