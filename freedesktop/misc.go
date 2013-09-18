package freedesktop

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func GetConfig(file string) string {
	return get(xdgConfig, file)
}

// Find a file in the XDG config path, returning the first one found
func GetAppConfig(file string) string {
	if AppName == "" {
		log.Printf("warning: app name not set")
		return ""
	}
	return GetConfig(path.Join(AppName, file))
}

// Find a file in the XDG data path, returning the first one found
func GetData(file string) string {
	return get(xdgData, file)
}

func GetAppData(file string) string {
	if AppName == "" {
		log.Printf("warning: app name not set")
		return ""
	}
	return GetData(path.Join(AppName, file))
}

func GetAllAppData(pattern string) []string {
	if AppName == "" {
		log.Printf("warning: app name not set")
		return make([]string, 0)
	}
	return GetAllData(path.Join(AppName, pattern))
}

// Searching through the XDG data path, return all files matching
// the given pattern
func GetAllData(pattern string) []string {
	return getAll(xdgData, pattern)
}

// Searches a directory list for the given file, then stops and returns
// its absolute path when found
func get(dirs []string, file string) string {
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
func getAll(dirs[]string, pattern string) []string {
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
	userDirs := GetConfig("user-dirs.dirs")
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

func XdgOpen(item string) error {
	return exec.Command("xdg-open", item).Run()
}
