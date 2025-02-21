package db

import (
	"fmt"
	"strings"

	"github.com/tecnologer/wheatley/pkg/contants/envvarname"
	"github.com/tecnologer/wheatley/pkg/utils/envvar"
)

const (
	DefaultDB = "wheatley"
)

type Config struct {
	Password string
	DBName   string
}

func NewConfigFromEnvVars() *Config {
	return &Config{
		Password: envvar.ValueWithDefault(envvarname.DBPassword, ""),
		DBName:   envvar.ValueWithDefault(envvarname.DBName, DefaultDB),
	}
}

func (c *Config) fileDBName() (string, error) {
	if err := c.OK(); err != nil {
		return "", fmt.Errorf("invalid db config: %w", err)
	}

	return fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", c.DBName, c.Password), nil
}

func (c *Config) OK() error {
	if c.DBName == "" {
		return fmt.Errorf("missing db name")
	}

	if c.Password == "" {
		return fmt.Errorf("missing db password")
	}

	if !strings.HasSuffix(c.DBName, ".db") {
		c.DBName += ".db"
	}

	return nil
}
