package flags

import "flag"

// Debug is a flag to turn debug mode.
var Debug = flag.Bool("debug", false, "Print debug information?")

func init() {
	flag.Parse()
}
