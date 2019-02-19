package unzip4win

import "flag"

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
