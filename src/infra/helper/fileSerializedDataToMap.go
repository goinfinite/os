package infraHelper

import (
	"encoding/json"
	"os"

	"github.com/goinfinite/os/src/domain/valueObject"
	"gopkg.in/yaml.v3"
)

func FileSerializedDataToMap(
	filePath valueObject.UnixFilePath,
) (outputMap map[string]interface{}, err error) {
	fileHandler, err := os.Open(filePath.String())
	if err != nil {
		return outputMap, err
	}
	defer fileHandler.Close()

	itemFileExt, err := filePath.GetFileExtension()
	if err != nil {
		return outputMap, err
	}

	isYamlFile := itemFileExt == "yml" || itemFileExt == "yaml"
	if isYamlFile {
		itemYamlDecoder := yaml.NewDecoder(fileHandler)
		err = itemYamlDecoder.Decode(&outputMap)
		if err != nil {
			return outputMap, err
		}

		return outputMap, nil
	}

	itemJsonDecoder := json.NewDecoder(fileHandler)
	err = itemJsonDecoder.Decode(&outputMap)
	if err != nil {
		return outputMap, err
	}

	return outputMap, nil
}
