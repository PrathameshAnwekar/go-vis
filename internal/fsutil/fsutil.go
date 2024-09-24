package fsutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PrathameshAnwekar/go-vis/internal/log"
)

func GetGoFiles(projectRoot string) ([]string, error) {
	log.D("Getting go files from:", projectRoot)
	fileList := make([]string, 0)
	err := filepath.WalkDir(projectRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".go") {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("error resolving the absolute path for: %s.\n Error: %e", path, err)
			}
			
			log.D("Adding: ", absPath)
			fileList = append(fileList, absPath)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileList, nil
}

func removeCommonPrefix(list []string) ([]string, error) {
	if len(list) <= 1 {
		return list, nil
	}
	prefix := list[0]
	for _, item := range list {
		for strings.Index(item, prefix) != 0 {
			prefix = prefix[:len(prefix)-1]
			if prefix == "" {
				break
			}
		}
	}
	for i, str := range list {
		list[i] = strings.TrimPrefix(str, prefix)
	}
	return list, nil
}
