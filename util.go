package gitcaddy

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands a leading ~/ to the users home directory
func ExpandPath(path string) string {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return filepath.Join(home, path[1:])
}

func FileExists(path string) bool {
	st, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	if st.IsDir() == true {
		return false
	}
	return true
}
