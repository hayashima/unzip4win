package unzip4win

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"os"
	"sort"
	"time"
)

type Arguments struct {
	ConfigFile string
	IsDebug    bool
	ZipFile    string
}

func ParseArgs() *Arguments {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, `Usage of Unzip4win:
  %s [OPITIONS] <zip-file-path>
Options
`, os.Args[0])
		flag.PrintDefaults()
	}

	configFile := flag.String("config", "", "Set path to customized config.toml.")
	isDebug := flag.Bool("debug", false, "If this flag is settle, output debug log!")
	flag.Parse()

	if flag.NArg() == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "zip file is not settle.")
		flag.Usage()
		os.Exit(2)
	}

	parsed := Arguments{
		ConfigFile: *configFile,
		IsDebug:    *isDebug,
		ZipFile:    flag.Arg(0),
	}

	return &parsed
}

type Config struct {
	Output   OutputConfig
	Password PasswordConfig
	Spec     []SpecConfig
}

type OutputConfig struct {
	SaveCurrent bool   `toml:"saveCurrent"`
	OutputPath  string `toml:"outputPath"`
}

type PasswordConfig struct {
	TryDays int `toml:"tryDays"`
}

type SpecConfig struct {
	Format    string    `toml:"format"`
	StartDate time.Time `toml:"startDate"`
}

func LoadConfig(configFile string) (*Config, error) {
	var config Config
	err := loadDefaultConfig(&config)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = overrideConfig(configFile, &config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// order by Spec.StartDate DESC
	sort.Sort(config)

	return &config, nil
}

func loadDefaultConfig(config *Config) error {
	f, err := Assets.Open("/config.toml")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = toml.DecodeReader(f, config)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func overrideConfig(configFile string, config *Config) error {
	if len(configFile) == 0 {
		return nil
	}
	var override Config
	meta, err := toml.DecodeFile(configFile, &override)
	if err != nil {
		return errors.WithStack(err)
	}
	if meta.IsDefined("output", "saveCurrent") {
		config.Output.SaveCurrent = override.Output.SaveCurrent
	}
	if meta.IsDefined("output", "outputPath") {
		config.Output.OutputPath = override.Output.OutputPath
	}
	if meta.IsDefined("password", "tryDays") {
		config.Password.TryDays = override.Password.TryDays
	}
	if meta.IsDefined("spec") {
		config.Spec = override.Spec
	}
	return nil
}

func (c Config) Len() int {
	return len(c.Spec)
}

func (c Config) Swap(i, j int) {
	c.Spec[i], c.Spec[j] = c.Spec[j], c.Spec[i]
}

func (c Config) Less(i, j int) bool {
	return c.Spec[i].StartDate.Sub(c.Spec[j].StartDate) > 0
}
