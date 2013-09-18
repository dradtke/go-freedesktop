package freedesktop

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

var xdgIconExtensions []string = []string{".png", ".svg", ".xpm"}

// Returns the current icon theme, if it can be found
func GetIconTheme() string {
	session := os.Getenv("DESKTOP_SESSION")
	switch session {
		case "gnome":
			cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "icon-theme")
			out, err := cmd.Output()
			if err == nil {
				return strings.Trim(string(out), "' \n")
			}
		case "kde":
			// TODO: look up KDE icon theme
			// ???: where is this setting stored?
			// try looking more in ~/.kde4/share/config/kdeglobals
	}
	return ""
}

func AppIcon(icon string) string {
	return AppIconForSize(icon, "48x48")
}

func AppIconForSize(icon, size string) (filename string) {
	theme := GetIconTheme()
	if theme != "" {
		filename = FindIconHelper(icon, size, theme)
		if filename != "" {
			return
		}
	}

	filename = FindIconHelper(icon, size, "hicolor")
	if filename != "" {
		return
	}

	return LookupFallbackIcon(icon)
}

func FindIconHelper(icon, size, theme string) (filename string) {
	filename = LookupIcon(icon, size, theme)
	if filename != "" {
		return
	}

	// TODO: look in this theme's parents
	// ???: how to get the icon theme's parents?
	return ""
}

func LookupIcon(icon, size, theme string) string {
	themeDirs := make([]string, 0)
	for _, dir := range xdgIcons {
		themeDir := path.Join(dir, theme)
		if _, err := os.Stat(themeDir); err != nil {
			// theme not found
			continue
		}
		themeDirs = append(themeDirs, themeDir)
	}

	for _, themeDir := range themeDirs {
		sizeDir := path.Join(themeDir, size)
		if _, err := os.Stat(sizeDir); err != nil {
			// size not found
			continue
		}

		categories, err := ioutil.ReadDir(sizeDir)
		if err != nil {
			// couldn't read directory
			continue
		}

		for _, category := range categories {
			categoryDir := path.Join(sizeDir, category.Name())

			for _, ext := range xdgIconExtensions {
				file := path.Join(categoryDir, icon + ext)
				if _, err := os.Stat(file); err == nil {
					return file
				}
			}
		}
	}

	return ""
}

func LookupFallbackIcon(icon string) string {
	for _, dir := range xdgIcons {
		for _, ext := range xdgIconExtensions {
			file := path.Join(dir, icon + ext)
			if _, err := os.Stat(file); err == nil {
				return file
			}
		}
	}

	return ""
}
