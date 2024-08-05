package configuration

import (
	"fmt"
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
	Port     string `yaml:"port"`
	Dbname   string `yaml:"dbname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type app struct {
	Name         string `yaml:"name"`
	LogMode      string `yaml:"logmode"`
	LogLevel     string `yaml:"loglevel"`
	LogFile      string `yaml:"logfile"`
	LogFileAudit string `yaml:"logfileAudit"`
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
