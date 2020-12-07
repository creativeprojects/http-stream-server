package cfg

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// Config from the file
type Config struct {
	Servers map[string]Server `yaml:"servers"`
}

// Server configuration
type Server struct {
	Listen      string `yaml:"listen"`
	Certificate string `yaml:"certificate"`
	PrivateKey  string `yaml:"privateKey"`
}

// LoadFileConfig loads the configuration from the file
func LoadFileConfig(fileName string) (Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	return loadConfig(file)
}

func loadConfig(reader io.Reader) (Config, error) {
	config := Config{}
	decoder := yaml.NewDecoder(reader)
	err := decoder.Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
