package infraHelper

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestFindCatalogItemFiles(t *testing.T) {
	catalogDirPath := t.TempDir()

	expectedFilePaths := []string{
		filepath.Join(catalogDirPath, "database", "redis", "manifest.json"),
		filepath.Join(catalogDirPath, "runtime", "node", "manifest.yml"),
		filepath.Join(catalogDirPath, "runtime", "php-webserver", "manifest.yaml"),
	}
	ignoredFilePaths := []string{
		filepath.Join(catalogDirPath, "runtime", "php-webserver", "assets", "modules.yaml"),
		filepath.Join(catalogDirPath, "runtime", "php-webserver", "assets", "manifest.yaml"),
		filepath.Join(catalogDirPath, "runtime", "php-webserver", "notes.txt"),
		filepath.Join(catalogDirPath, "runtime", "node", "extra.yaml"),
		filepath.Join(catalogDirPath, ".hidden", "manifest.yaml"),
		filepath.Join(catalogDirPath, ".hidden-manifest.yaml"),
	}

	fixtureFilePaths := slices.Concat(expectedFilePaths, ignoredFilePaths)
	for _, fixtureFilePath := range fixtureFilePaths {
		err := os.MkdirAll(filepath.Dir(fixtureFilePath), 0755)
		if err != nil {
			t.Fatalf("CreateCatalogFixtureDirFailed: %v", err)
		}

		err = os.WriteFile(fixtureFilePath, []byte("name: test"), 0644)
		if err != nil {
			t.Fatalf("CreateCatalogFixtureFileFailed: %v", err)
		}
	}

	actualFilePaths, err := FindCatalogItemFiles(catalogDirPath)
	if err != nil {
		t.Fatalf("FindCatalogItemFilesFailed: %v", err)
	}

	slices.Sort(actualFilePaths)
	slices.Sort(expectedFilePaths)
	if !slices.Equal(actualFilePaths, expectedFilePaths) {
		t.Errorf(
			"CatalogItemFilesMismatch: expected %v, got %v",
			expectedFilePaths, actualFilePaths,
		)
	}
}
