package bootstrap

import (
	"runtime"

	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

// ConfigureMaxProc configures the program's GOMAXPROCS to the given
// value if not 0. If 0, the default GOMAXPROCS will be used.
func ConfigureMaxProc(overrideMax int) {

	if overrideMax == 0 {
		if _, err := maxprocs.Set(maxprocs.Logger(func(msg string, args ...any) {})); err != nil {
			zap.L().Fatal("Unable to set automaxprocs", zap.Error(err))
		}
	} else {
		runtime.GOMAXPROCS(overrideMax)
	}

	zap.L().Info("GOMAXPROCS configured", zap.Int("max", runtime.GOMAXPROCS(0)))
}
