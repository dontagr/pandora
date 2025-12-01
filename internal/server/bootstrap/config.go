package bootstrap

import (
	"fmt"

	"go.uber.org/fx"

	configInternal "github.com/dontagr/pandora/internal/server/config"
	"github.com/dontagr/pandora/pkg/config"
)

var Config = fx.Options(
	fx.Provide(newConfig),
)

func newConfig() (*configInternal.Config, error) {
	configInt := &configInternal.Config{}

	cnf := &config.Config{
		Data:             configInt,
		DefaultFilePaths: []string{"../../../configs", "./configs"},
		DefaultFileNames: []string{"server.config.json"},
	}

	err := cnf.ReadFromFile()
	if err != nil {
		return nil, fmt.Errorf("ReadFromFile: %v", err)
	}
	err = cnf.ReadFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to read from env: %w", err)
	}

	err = cnf.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return configInt, nil
}
