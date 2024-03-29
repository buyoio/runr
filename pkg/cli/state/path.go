package state

import (
	"os"
	"path/filepath"
)

const (
	configFileName = ".runr.yaml"
	secretFileName = ".runr.secrets"
)

func resolvePath(path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(
		filepath.Dir(getDefaultConfigPath()),
		path,
	)
}

func getDefaultConfigPath() string {
	if config := os.Getenv("LOK8S_CONFIG"); config != "" {
		return config
	}
	config := configWalkParents()
	if config != "" {
		return config
	}
	// use current working directory
	return configFileName
}

func configWalkParents() string {
	base, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		file := filepath.Join(base, configFileName)
		if _, err := os.Stat(file); err == nil {
			return file
		}
		base = filepath.Dir(base)
		if base == "/" || base == "." || base == "" ||
			base == string(filepath.Separator) ||
			base == filepath.VolumeName(base) {
			break
		}
	}
	return ""
}
