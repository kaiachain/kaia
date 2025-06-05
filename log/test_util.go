package log

import (
	"flag"
	"io"
	"os"
	"testing"

	"github.com/kaiachain/kaia/v2/log/term"
	"github.com/mattn/go-colorable"
)

// Enable logging to STDERR
// Exmaple use
//
//	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
//
// `normalLvl` is used in most cases
// `verboseLvl` is used if `go test -v` flag is given
func EnableLogForTest(normalLvl, verboseLvl Lvl) {
	lvl := Lvl(normalLvl)
	// testing.Verbose() panics in the init() function. flag.Parsed() prevents that.
	if flag.Parsed() && testing.Verbose() {
		lvl = Lvl(verboseLvl)
	}

	usecolor := term.IsTty(os.Stderr.Fd()) && os.Getenv("TERM") != "dumb"
	output := io.Writer(os.Stderr)
	if usecolor {
		output = colorable.NewColorableStderr()
	}

	glogger := NewGlogHandler(StreamHandler(output, TerminalFormat(usecolor)))
	PrintOrigins(true)
	ChangeGlobalLogLevel(glogger, lvl)
	glogger.Vmodule("")
	glogger.BacktraceAt("")
	Root().SetHandler(glogger)
}
