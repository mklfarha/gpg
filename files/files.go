package files

import (
	"path/filepath"
	"runtime"
)

func AppDir() string {
	_, callerFile, _, _ := runtime.Caller(1)
	generatorDir := filepath.Dir(callerFile)
	absoluteGeneratorDir, err := filepath.Abs(generatorDir)
	if err != nil {
		panic("could not resolve path")
	}
	return filepath.Dir(absoluteGeneratorDir)
}
