package flags

import "flag"

const helpFlag = "help"
const showInderectFlag = "show-inderect"
const verboseFlag = "verbose"
const urlFlag = "url"
const configFlag = "config-path"
const defaultConfigPath = "update-checker-config.yml"

type Flags struct {
	ConfigPath   *string
	Url          *string
	Verbose      *bool
	ShowInderect *bool
	Help         *bool
}

func ParseFlags() *Flags {
	res := &Flags{
		ConfigPath:   flag.String(configFlag, defaultConfigPath, "set config loading path"),
		Help:         flag.Bool(helpFlag, false, "show help"),
		ShowInderect: flag.Bool(showInderectFlag, false, "show indirect requirements"),
		Verbose:      flag.Bool(verboseFlag, false, "show extended information"),
		Url:          flag.String(urlFlag, "", "repository url"),
	}

	flag.Parse()

	return res
}
