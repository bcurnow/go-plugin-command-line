package plugin

import (
	"fmt"
	"os"
	"path/filepath"
)

// Traverses the dir and calls the register function for any executable
func Register(pluginType string, dir string, register func(pluginName string, pluginCmd string) error) error {
	logger.Debug(fmt.Sprintf("Loading %s", pluginType), "Dir", dir)
	executables, err := discover(dir)
	if err != nil {
		return err
	}

	for pluginName, pluginCmd := range executables {
		err = register(pluginName, pluginCmd)
		if err != nil {
			return err
		}
	}

	return nil
}

// Until goplugin.Discover is updated to check for the executable bit, this is our own implementation
func discover(dir string) (map[string]string, error) {
	var executables map[string]string = make(map[string]string)

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Don't traverse sub-directories, this is arbitrary but we are keeping it simple
		if d.IsDir() && path != dir {
			logger.Warn("Subdirectories are not supported", "Subdirectory", path)
			return filepath.SkipDir
		}

		// Because we're using WalkDir, we need to get the FileInfo from the DirEntry
		info, err := d.Info()
		if err != nil {
			return err
		}

		// Check if this is a file and if the file is executable
		if info.Mode().IsRegular() {
			// 0111 checks for the execute bit to be set
			if info.Mode()&0111 == 0 {
				logger.Warn("Skipping non-executable file", "File", path)
				return nil
			}

			// Get the absolute path of the file so we can provide the best debugging information
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			executables[filepath.Base(path)] = absPath
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return executables, nil
}
