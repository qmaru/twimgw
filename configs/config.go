package configs

import (
	"path/filepath"

	"twimgw/utils"
)


// readFile Read configuration files in configs
func readFile(name string) (map[string]any, error) {
	cfgRoot, err := utils.FileSuite.RootPath("configs")
	if err != nil {
		return nil, err
	}

	cfgPath := filepath.Join(cfgRoot, name)
	err = utils.FileSuite.IsExist(cfgPath)
	if err != nil {
		return nil, err
	}

	raw, err := utils.FileSuite.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	data, err := utils.DataSuite.RawMap2Map(raw)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Secrets load twitter API secret config
func Secrets() (map[string]any, error) {
	return readFile("secrets.json")
}
