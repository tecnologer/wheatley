package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tecnologer/wheatley/pkg/contants/envvarname"
	"github.com/tecnologer/wheatley/pkg/dao/db"
)

func TestNewConfigFromEnvVars(t *testing.T) {
	t.Run("no_env_var_data", func(t *testing.T) {
		t.Setenv(envvarname.DBPassword, "")
		t.Setenv(envvarname.DBName, "")

		want := &db.Config{
			Password: "",
			DBName:   db.DefaultDB,
		}

		got := db.NewConfigFromEnvVars()
		assert.Equal(t, want, got)
	})

	t.Run("env_var_data", func(t *testing.T) {
		t.Setenv(envvarname.DBPassword, "test_password")
		t.Setenv(envvarname.DBName, "test_db")

		want := &db.Config{
			Password: "test_password",
			DBName:   "test_db",
		}

		got := db.NewConfigFromEnvVars()
		assert.Equal(t, want, got)
	})
}
