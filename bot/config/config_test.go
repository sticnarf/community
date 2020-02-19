package config

import (
	"path"
	"runtime"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func GetTestConfig() (*Config, error) {
	_, localFile, _, _ := runtime.Caller(0)
	pathStr := path.Join(path.Dir(localFile), "../config.example.toml")
	cfg, err := GetConfig(&pathStr)
	if err != nil {
		return nil, errors.Wrap(err, "get test config")
	}
	return cfg, nil
}

func TestT(t *testing.T) {
	cfg, err := GetTestConfig()
	// if err != nil {
	// 	t.Errorf("create bot failed, %+v", err)
	// }
	assert.Equal(t, err, nil, "create bot failed")
	assert.Equal(t, cfg.Database.Address, "127.0.0.1", "wrong address")
	assert.Equal(t, cfg.Repos["owner-repo"].WebhookSecret,
		"secret", "wrong repo secret")
}
