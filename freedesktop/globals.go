package freedesktop

import (
	"log"
	"os"
	"os/user"
	"path"
	"strings"
)

var xdgData, xdgConfig, xdgIcons []string
var currentLocale string
var AppName string

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err.Error())
	}

	// load the user's xdg config and data dirs

	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = path.Join(usr.HomeDir, ".config")
	}
	xdgConfigDirs := os.Getenv("XDG_CONFIG_DIRS")
	if xdgConfigDirs == "" {
		xdgConfigDirs = "/etc/xdg"
	}
	xdgConfigDirsList := strings.Split(xdgConfigDirs, string(os.PathListSeparator))
	xdgConfig = make([]string, len(xdgConfigDirsList) + 1)
	xdgConfig[0] = xdgConfigHome
	for i, dir := range xdgConfigDirsList {
		xdgConfig[i+1] = dir
	}

	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		xdgDataHome = path.Join(usr.HomeDir, ".local/share")
	}
	xdgDataDirs := os.Getenv("XDG_DATA_DIRS")
	if xdgDataDirs == "" {
		xdgDataDirs = "/usr/local/share/:/usr/share/"
	}
	xdgDataDirsList := strings.Split(xdgDataDirs, string(os.PathListSeparator))
	xdgData = make([]string, len(xdgDataDirsList) + 1)
	xdgIcons = make([]string, len(xdgDataDirsList) + 2)
	xdgData[0] = xdgDataHome
	xdgIcons[0] = path.Join(usr.HomeDir, ".icons")
	for i, dir := range xdgDataDirsList {
		xdgData[i+1] = dir
		xdgIcons[i+1] = path.Join(dir, "icons")
	}
	xdgIcons[len(xdgIcons)-1] = "/usr/share/pixmaps"

	// get the current locale

	currentLocale = os.Getenv("LANG")
}
