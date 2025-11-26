package flags

import "flag"

const helpFlag = "help"
const showIndirectFlag = "show-indirect"
const verboseFlag = "verbose"
const urlFlag = "url"
const timeoutFlag = "timeout-seconds"
const defaultTimeout = 10
const configFlag = "config-path"

type Flags struct {
	ConfigPath     *string
	Url            *string
	Verbose        *bool
	ShowIndirect   *bool
	Help           *bool
	TimeoutSeconds *int
}

func ParseFlags() *Flags {
	res := &Flags{
		ConfigPath:     flag.String(configFlag, "", "path to config with sensitive data"),
		Help:           flag.Bool(helpFlag, false, "show help"),
		ShowIndirect:   flag.Bool(showIndirectFlag, false, "show indirect requirements"),
		Verbose:        flag.Bool(verboseFlag, false, "show extended information"),
		Url:            flag.String(urlFlag, "", "repository url"),
		TimeoutSeconds: flag.Int(timeoutFlag, defaultTimeout, "timeout for repository scaning"),
	}

	flag.Parse()

	return res
}
