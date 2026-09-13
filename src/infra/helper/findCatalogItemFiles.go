package infraHelper

import (
	"errors"
	"fmt"
	"strings"

	tkInfra "github.com/goinfinite/tk/src/infra"
)

const catalogItemManifestMaxDepth string = "3"

func FindCatalogItemFiles(catalogDirPath string) ([]string, error) {
	searchCommand := fmt.Sprintf(
		"find %s -maxdepth %s -type f \\( -name 'manifest.json' -o "+
			"-name 'manifest.yml' -o -name 'manifest.yaml' \\) "+
			"-not -path '*/.*' -not -name '.*'",
		catalogDirPath, catalogItemManifestMaxDepth,
	)

	rawFilePaths, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command:           searchCommand,
		ShouldUseSubShell: true,
	}).Run()
	if err != nil {
		return nil, errors.New("FindCatalogItemFilesFailed: " + err.Error())
	}
	if rawFilePaths == "" {
		return []string{}, nil
	}

	return strings.Split(rawFilePaths, "\n"), nil
}
