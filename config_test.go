package goconfig

import (
	"fmt"
	"github.com/naucon/goconfig/examples"
	logMock "github.com/naucon/goconfig/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"os"
	"testing"
)

func TestConfig_NewConfig(t *testing.T) {
	setup()

	t.Run("TestConfig_NewConfig_prod_ShouldPopulateConfigByUsingEnvLocal", func(t *testing.T) {
		os.Clearenv()
		cfg := examples.Configuration{}
		options := Options{
			ConfigPath: "./examples/",
		}
		c := NewConfig(options)
		err := c.Load("prod", &cfg)
		assert.NoError(t, err)

		assert.Equal(t, "secret_dsn_local", cfg.Database.Dsn)
		assert.Equal(t, 3000, cfg.Server.Port)
		assert.Contains(t, cfg.Server.TrustedProxies, "127.0.0.1")
		assert.Equal(t, "secret default", cfg.Secret)
	})

	t.Run("TestConfig_NewConfig_dev_ShouldPopulateConfigByUsingEnvLocal", func(t *testing.T) {
		os.Clearenv()
		cfg := examples.Configuration{}
		options := Options{
			ConfigPath: "./examples/",
		}
		c := NewConfig(options)
		err := c.Load("dev", &cfg)
		assert.NoError(t, err)

		assert.Equal(t, "secret_dsn_local", cfg.Database.Dsn)
		assert.Equal(t, 3000, cfg.Server.Port)
		assert.Contains(t, cfg.Server.TrustedProxies, "127.0.0.1")
		assert.Equal(t, "secret default", cfg.Secret)
	})

	t.Run("TestConfig_NewConfig_Test_ShouldPopulateConfigByUsingEnvTest", func(t *testing.T) {
		os.Clearenv()
		cfg := examples.Configuration{}
		options := Options{
			ConfigPath: "./examples/",
		}
		c := NewConfig(options)
		err := c.Load("test", &cfg)
		assert.NoError(t, err)

		assert.Equal(t, "secret_dsn_test", cfg.Database.Dsn)
		assert.Equal(t, 3000, cfg.Server.Port)
		assert.Contains(t, cfg.Server.TrustedProxies, "127.0.0.1")
		assert.Equal(t, "secret default", cfg.Secret)
	})

	t.Run("TestConfig_NewConfig_EmptyEnv_ShouldReturnError", func(t *testing.T) {
		os.Clearenv()
		cfg := examples.Configuration{}
		options := Options{
			ConfigPath: "./examples/",
		}
		c := NewConfig(options)
		err := c.Load("", &cfg)
		assert.Error(t, err)

		assert.Equal(t, errConfigEmptyEnv, err.Error())
	})

	t.Run("TestConfig_NewConfig_MissingDefaultConfig_ShouldPopulateReturnErr", func(t *testing.T) {
		os.Clearenv()
		cfg := examples.Configuration{}
		options := Options{}
		c := NewConfig(options)
		err := c.Load("prod", &cfg)
		assert.Error(t, err)
	})

	t.Run("TestConfig_NewConfig_test_WithConfigPath_ShouldPopulateConfig", func(t *testing.T) {
		os.Clearenv()
		cfg := examples.Configuration{}
		options := Options{
			ConfigPath:          "./testdata/valid/",
			ConfigFileExtension: "yaml",
		}
		c := NewConfig(options)
		err := c.Load("test", &cfg)
		assert.NoError(t, err)

		assert.Equal(t, "secret_dsn_test", cfg.Database.Dsn)
		assert.Equal(t, 3000, cfg.Server.Port)
		assert.Contains(t, cfg.Server.TrustedProxies, "127.0.0.1")
		assert.Equal(t, "secret default", cfg.Secret)
	})

	t.Run("TestConfig_NewConfig_test_WithConfigAndEnvPath_ShouldPopulateConfigWithoutUsingEnvLocal", func(t *testing.T) {
		os.Clearenv()
		cfg := examples.Configuration{}
		options := Options{
			DotEnvPath:          "./testdata/",
			ConfigPath:          "./testdata/valid/",
			ConfigFileExtension: "yaml",
		}
		c := NewConfig(options)
		err := c.Load("test", &cfg)
		assert.NoError(t, err)

		assert.Equal(t, "secret_dsn_test", cfg.Database.Dsn)
		assert.Equal(t, 3000, cfg.Server.Port)
		assert.Contains(t, cfg.Server.TrustedProxies, "127.0.0.1")
		assert.Equal(t, "secret default", cfg.Secret)
	})

	t.Run("TestConfig_NewConfig_prod_WithConfigAndEnvPath_ShouldPopulateConfigWithUsingEnvProdLocal", func(t *testing.T) {
		os.Clearenv()
		cfg := examples.Configuration{}
		options := Options{
			DotEnvPath:          "./testdata/",
			ConfigPath:          "./testdata/valid/",
			ConfigFileExtension: "yaml",
		}
		c := NewConfig(options)
		err := c.Load("prod", &cfg)
		assert.NoError(t, err)

		assert.Equal(t, "secret_dsn_prod_local", cfg.Database.Dsn)
		assert.Equal(t, 3000, cfg.Server.Port)
		assert.Contains(t, cfg.Server.TrustedProxies, "127.0.0.1")
		assert.Equal(t, "secret default", cfg.Secret)
	})

	t.Run("TestConfig_NewConfig_test_EmptyYaml_ShouldKeepDefaultConfig", func(t *testing.T) {
		os.Clearenv()
		cfg := examples.Configuration{
			Server: examples.Server{
				Port: 3000,
			},
		}
		options := Options{
			DotEnvPath:          "./testdata/",
			ConfigPath:          "./testdata/empty/",
			ConfigFileExtension: "yaml",
		}
		c := NewConfig(options)
		err := c.Load("test", &cfg)
		assert.NoError(t, err)

		assert.Equal(t, "", cfg.Database.Dsn)
		assert.Equal(t, 3000, cfg.Server.Port)
		assert.Empty(t, cfg.Server.TrustedProxies)
		assert.Equal(t, "", cfg.Secret)
	})

	t.Run("TestConfig_NewConfig_test_InvalidYaml_ShouldKeepDefaultConfig", func(t *testing.T) {
		os.Clearenv()
		cfg := examples.Configuration{
			Server: examples.Server{
				Port: 8000,
			},
		}
		options := Options{
			DotEnvPath:          "./testdata/",
			ConfigPath:          "./testdata/invalid/",
			ConfigFileExtension: "yaml",
		}
		c := NewConfig(options)
		err := c.Load("test", &cfg)
		assert.Error(t, err)

		assert.Equal(t, "", cfg.Database.Dsn)
		assert.Equal(t, 8000, cfg.Server.Port)
		assert.Empty(t, cfg.Server.TrustedProxies)
		assert.Equal(t, "", cfg.Secret)
	})

	t.Run("TestConfig_NewConfig_test_Missing_ShouldKeepDefaultConfig", func(t *testing.T) {
		os.Clearenv()
		cfg := examples.Configuration{
			Server: examples.Server{
				Port: 8000,
			},
		}
		options := Options{
			DotEnvPath:          "./testdata/",
			ConfigPath:          "./testdata/missing/",
			ConfigFileExtension: "yaml",
		}
		c := NewConfig(options)
		err := c.Load("test", &cfg)
		assert.Error(t, err)

		assert.Equal(t, "", cfg.Database.Dsn)
		assert.Equal(t, 8000, cfg.Server.Port)
		assert.Empty(t, cfg.Server.TrustedProxies)
		assert.Equal(t, "", cfg.Secret)
	})

	t.Run("TestConfig_NewConfig_ProdVerbose_ShouldLog", func(t *testing.T) {
		os.Clearenv()

		logger := logMock.NewLoggerMock()
		logger.On("Printf", mock.Anything, mock.Anything).Times(7)

		cfg := examples.Configuration{}
		options := Options{
			DotEnvPath:          "./testdata/",
			ConfigPath:          "./testdata/valid/",
			ConfigFileExtension: "yaml",
			Verbose:             true,
			Logger:              logger,
		}
		c := NewConfig(options)
		err := c.Load("prod", &cfg)
		assert.NoError(t, err)

		logger.AssertExpectations(t)
	})

	t.Run("TestConfig_NewConfig_TestVerbose_ShouldLog", func(t *testing.T) {
		os.Clearenv()
		cfg := examples.Configuration{}
		logger := logMock.NewLoggerMock()
		logger.On("Printf", mock.Anything, mock.Anything).Times(5)
		options := Options{
			ConfigPath: "./examples/",
			Verbose:    true,
			Logger:     logger,
		}
		c := NewConfig(options)
		err := c.Load("test", &cfg)
		assert.NoError(t, err)

		logger.AssertExpectations(t)
	})

	t.Run("TestConfig_NewConfig_TestMissingVerbose_ShouldLog", func(t *testing.T) {
		os.Clearenv()

		logger := logMock.NewLoggerMock()
		logger.On("Printf", mock.Anything, mock.Anything).Times(4)

		cfg := examples.Configuration{
			Server: examples.Server{
				Port: 8000,
			},
		}
		options := Options{
			DotEnvPath:          "./testdata/",
			ConfigPath:          "./testdata/missing/",
			ConfigFileExtension: "yaml",
			Verbose:             true,
			Logger:              logger,
		}
		c := NewConfig(options)
		err := c.Load("test", &cfg)
		assert.Error(t, err)

		logger.AssertExpectations(t)
	})

	teardown()
}

func setup() {
	_ = copy("testdata/default.env", ".env")
	_ = copy("testdata/local.env", ".env.local")
	_ = copy("testdata/test.env", ".env.test")
	_ = copy("testdata/default.env", "testdata/.env")
	_ = copy("testdata/local.env", "testdata/.env.local")
	_ = copy("testdata/prod.env", "testdata/.env.prod")
	_ = copy("testdata/test.env", "testdata/.env.test")
	_ = copy("testdata/prod.local.env", "testdata/.env.prod.local")
}

func teardown() {
	_ = remove(".env")
	_ = remove(".env.local")
	_ = remove(".env.test")
	_ = remove("testdata/.env")
	_ = remove("testdata/.env.local")
	_ = remove("testdata/.env.prod")
	_ = remove("testdata/.env.test")
	_ = remove("testdata/.env.prod.local")
}

func copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func remove(src string) error {
	err := os.Remove(src)
	return err
}
