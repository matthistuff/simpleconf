package simpleconf

import (
	"github.com/mitchellh/go-homedir"
	"path"
	"os"
	"github.com/BurntSushi/toml"
	"errors"
)

type Config struct {
	path string
	File string
	Data interface{}
}

func (c *Config) Ensure() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	c.path = path.Join(home, c.File)
	_, err = os.Stat(c.path);
	if os.IsNotExist(err) {
		c.Flush()
	}

	return nil
}

func (c *Config) Load() error {
	if _, err := toml.DecodeFile(c.path, &c.Data); err != nil {
		return err
	}

	return nil
}

func (c *Config) Flush() error {
	if c.Data == nil {
		return errors.New("Cannot flush config without loading/creating it first!")
	}

	file, err := os.Create(c.path)
	if err != nil {
		panic(err)
	}

	encoder := toml.NewEncoder(file)
	err = encoder.Encode(c.Data)
	if err != nil {
		panic(err)
	}

	return nil
}

func New(fileName string, defaultData interface{}) (*Config, error) {
	conf := &Config{
		File: fileName,
		Data: defaultData,
	}

	if err := conf.Ensure(); err != nil {
		return nil, err
	}

	conf.Load()

	return conf, nil
}