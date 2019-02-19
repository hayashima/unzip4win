package unzip4win

import (
	"flag"
	"github.com/BurntSushi/toml"
	"sort"
	"time"
)

type Arguments struct {
	ConfigFile string
	IsDebug    bool
	ZipFile    string
}

func ParseArgs() *Arguments {
	configFile := flag.String("config", "config.toml", "Set config path.")
	isDebug := flag.Bool("debug", false, "If this flag is settle, output debug log!")
	flag.Parse()

	parsed := Arguments{
		ConfigFile: *configFile,
		IsDebug:    *isDebug,
		ZipFile:    "",
	}

	if flag.NArg() > 0 {
		parsed.ZipFile = flag.Arg(0)
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
	_, err := toml.DecodeFile(configFile, &config)
	// order by Spec.StartDate DESC
	sort.Sort(config)

	return &config, err
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
