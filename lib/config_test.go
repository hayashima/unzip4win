package unzip4win

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestOverrideConfig(t *testing.T) {
	defaultConfig := Config{}
	t.Run("override output config", func(t *testing.T) {
		config := Config{}
		_ = overrideConfig(configPath("override_output_config.toml"), &config)
		if reflect.DeepEqual(defaultConfig.Output, config.Output) {
			t.Errorf("failed to override output config.\nactual: %v", config)
		}
		if !reflect.DeepEqual(defaultConfig.Password, config.Password) ||
			!reflect.DeepEqual(defaultConfig.Spec, config.Spec) {
			t.Errorf("unwanted config is overridden.\nactual: %v", config)
		}
	})

	t.Run("override password config", func(t *testing.T) {
		config := Config{}
		_ = overrideConfig(configPath("override_password_config.toml"), &config)
		if reflect.DeepEqual(defaultConfig.Password, config.Password) {
			t.Errorf("failed to override password config.\nactual: %v", config)
		}
		if !reflect.DeepEqual(defaultConfig.Output, config.Output) ||
			!reflect.DeepEqual(defaultConfig.Spec, config.Spec) {
			t.Errorf("unwanted config is overridden.\nactual: %v", config)
		}
	})

	t.Run("override spec config", func(t *testing.T) {
		config := Config{}
		_ = overrideConfig(configPath("override_specs_config.toml"), &config)
		if reflect.DeepEqual(defaultConfig.Spec, config.Spec) {
			t.Errorf("failed to override spec config.\nactual: %v", config)
		}
		if !reflect.DeepEqual(defaultConfig.Output, config.Output) ||
			!reflect.DeepEqual(defaultConfig.Password, config.Password) {
			t.Errorf("unwanted config is overridden.\nactual: %v", config)
		}
	})

	t.Run("override all config", func(t *testing.T) {
		config := Config{}
		_ = overrideConfig(configPath("override_all_config.toml"), &config)
		if reflect.DeepEqual(defaultConfig.Output, config.Output) ||
			reflect.DeepEqual(defaultConfig.Password, config.Password) ||
			reflect.DeepEqual(defaultConfig.Spec, config.Spec) {
			t.Errorf("failed to override wanted config.\nactual: %v", config)
		}
	})
}

func configPath(name string) string {
	return filepath.Join("..", "_tests", "config", name)
}
