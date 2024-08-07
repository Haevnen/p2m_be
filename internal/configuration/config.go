package configuration

import (
	"fmt"
	"github.com/Haevnen/p2m_be/pkg/gormdb"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config struct for project config
type Config struct {
	App      app      `yaml:"app"`
	Server   server   `yaml:"server"`
	Database database `yaml:"database"`
}

type database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBname   string `yaml:"dbname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type app struct {
	Name         string `yaml:"name"`
	LogMode      string `yaml:"logmode"`
	LogLevel     string `yaml:"loglevel"`
	LogFile      string `yaml:"logfile"`
	LogFileAudit string `yaml:"logfileAudit"`
	Mode         string `yaml:"mode"`
}

type server struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Monitor string `yaml:"monitor"`
}

var (
	AppConfig *Config
)

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Validate config path
	if err := ValidateConfigPath(configPath); err != nil {
		return nil, err
	}

	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(filepath.Clean(configPath))
	if err != nil {
		return nil, err
	}

	// Solved gosec issue G307 (CWE-703)
	// See: https://github.com/securego/gosec/issues/512#issuecomment-675286833
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	AppConfig = config
	return config, nil
}

// ValidateConfigPath just makes sure,
// that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	// Check path
	s, err := os.Stat(filepath.Clean(path))
	if err != nil {
		return err
	}

	// Check for directory
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory", path)
	}

	return nil
}

// GetGORMConfig build gormdb config from env
func (e *Config) GetGORMConfig() *gormdb.Config {
	return &gormdb.Config{
		DBHost:            e.Database.Host,
		DBPort:            e.Database.Port,
		DBUser:            e.Database.Username,
		DBPass:            e.Database.Password,
		DBName:            e.Database.DBname,
		LogSQL:            true,
		MaxOpenConn:       10,
		MaxLifetimeSecond: 300,
	}
}

// GetURLBase build server config from env
func (e *Config) GetURLBase() string {
	return fmt.Sprintf("%s:%s", e.Server.Host, e.Server.Port)
}

// GetURLProfile build server config from env
func (e *Config) GetURLProfile() string {
	return fmt.Sprintf("%s:%s", e.Server.Host, e.Server.Monitor)
}
