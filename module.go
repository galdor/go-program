package program

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"runtime"
)

func ModuleDirectoryPath() (string, error) {
	const refFilePath = "go.mod"

	// Skip one call frame, i.e. the one associated with this very function,
	// since we want to locate the module of the calling function.
	_, filePath, _, _ := runtime.Caller(1)

	dirPath := path.Dir(filePath)

	for dirPath != "/" {
		dirPath = path.Join(dirPath, "..")

		filePath := path.Join(dirPath, refFilePath)

		_, err := os.Stat(filePath)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		} else if err != nil {
			return "", fmt.Errorf("cannot stat %q: %v", filePath, err)
		}

		break
	}

	if dirPath == "/" {
		return "", fmt.Errorf("%q not found in parent directories", refFilePath)
	}

	return dirPath, nil
}
