package freedesktop

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Find a file in the XDG config path, returning the first one found
func FindInConfig(file string) string {
	return find(xdgConfig, file)
}

// Find a file in the XDG data path, returning the first one found
func FindInData(file string) string {
	return find(xdgData, file)
}

// Searching through the XDG data path, return all files matching
// the given pattern
func CollectFromData(pattern string) []string {
	return collect(xdgData, pattern)
}

// Searches a directory list for the given file, then stops and returns
// its absolute path when found
func find(dirs []string, file string) string {
	for _, dir := range dirs {
		p := path.Join(dir, file)
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// Iterates through a directory list and collects all files matching
// the given pattern
func collect(dirs[]string, pattern string) []string {
	collection := make(map[string] bool)
	for _, dir := range dirs {
		matches, _ := filepath.Glob(path.Join(dir, pattern))
		for _, match := range matches {
			if _, ok := collection[match]; !ok {
				collection[match] = true
			}
		}
	}

	results := make([]string, len(collection)); i := 0
	for key, _ := range collection {
		results[i] = key
		i++
	}
	return results
}

/* --- Miscellaneous --- */

// Takes the name of a user directory, e.g. "music" or "pictures"
// and returns its absolute path
func GetUserDir(d string) (string, error) {
	userDirs := FindInConfig("user-dirs.dirs")
	data, err := ioutil.ReadFile(userDirs)
	if err != nil {
		return "", err
	}

	name := "XDG_" + strings.ToUpper(d) + "_DIR"
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}

		split := strings.Split(line, "=")
		if len(split) != 2 {
			continue
		}

		var key, value string = split[0], split[1]

		if key == name {
			return os.ExpandEnv(strings.Trim(value, "\"")), nil
		}
	}

	return "", nil
}
