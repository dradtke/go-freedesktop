package freedesktop

import (
	"container/list"
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
)

var entryTypes []string = []string{"Application", "Link", "Directory"}
var entryKeys []string = []string{"Type", "Version", "Name", "GenericName",
	"NoDisplay", "Comment", "Icon", "Hidden", "OnlyShowIn", "NotShowIn",
	"TryExec", "Exec", "Path", "Terminal", "Actions", "MimeType",
	"Categories", "Keywords", "StartupNotify", "StartupWMClass", "URL"}

type DesktopEntry struct {
	/* required fields */
	Name string
	Type string
	Exec string
	URL string  // only required for Link entries

	NoDisplay bool
	Hidden bool
	Terminal bool
	StartupNotify bool

	Version string
	GenericName string
	Comment string
	Icon string
	OnlyShowIn []string ; NotShowIn []string
	TryExec string
	Path string
	Actions []string
	MimeType []string
	Categories []string
	Keywords []string
	StartupWMClass string

	File string // .desktop file
}

// Config files are maps from Group Header -> Key -> Value
func ParseConfigFile(file string) (map[string] map[string] string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	keyRegex, _ := regexp.Compile("[A-Za-z0-9-]+")
	cfg := make(map[string] map[string] string)
	lines := strings.Split(string(data), "\n")
	var header string
	for _, line := range lines {
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			header = strings.TrimRight(strings.TrimLeft(line, "["), "]")
			continue
		}

		if header != "" {
			split := strings.Split(line, "=")
			if len(split) != 2 {
				continue
			}

			var key, value string = strings.TrimSpace(split[0]), strings.TrimSpace(split[1])
			if !keyRegex.MatchString(key) {
				return nil, errors.New("invalid key: " + key)
			}
			if _, ok := cfg[header]; !ok {
				cfg[header] = make(map[string] string)
			}
			cfg[header][key] = value
		}
	}

	return cfg, nil
}

func ParseDesktopEntry(file string) (*DesktopEntry, error) {
	cfg, err := ParseConfigFile(file)
	if err != nil {
		return nil, err
	}

	deHeader := cfg["Desktop Entry"]
	entry := new(DesktopEntry)

	// get the entry type

	entry.Type = deHeader["Type"]
	if entry.Type == "" {
		return nil, errors.New("missing required key: Type")
	}
	if !DesktopEntryTypeIsValid(entry.Type) {
		return nil, missingKeyError("Type")
	}

	// other required fields

	entry.Name = GetLocalizedValue(deHeader, "Name")
	if entry.Name == "" {
		return nil, missingKeyError("Name")
	}

	entry.Exec = deHeader["Exec"]
	if entry.Exec == "" {
		return nil, missingKeyError("Exec")
	}

	entry.URL = deHeader["URL"]
	if entry.URL == "" && entry.Type == "Link" {
		return nil, missingKeyError("URL")
	}

	// non-required fields that need validation

	noDisplay := deHeader["NoDisplay"]
	if noDisplay == "" {
		entry.NoDisplay = false
	} else if val, err := getBoolValue(deHeader["NoDisplay"]); err != nil {
		return nil, err
	} else {
		entry.NoDisplay = val
	}

	hidden := deHeader["Hidden"]
	if hidden == "" {
		entry.Hidden = false
	} else if val, err := getBoolValue(deHeader["Hidden"]); err != nil {
		return nil, err
	} else {
		entry.Hidden = val
	}

	terminal := deHeader["Terminal"]
	if terminal == "" {
		entry.Terminal = false
	} else if val, err := getBoolValue(deHeader["Terminal"]); err != nil {
		return nil, err
	} else {
		entry.Terminal = val
	}

	startupNotify := deHeader["StartupNotify"]
	if startupNotify == "" {
		entry.StartupNotify = false
	} else if val, err := getBoolValue(deHeader["StartupNotify"]); err != nil {
		return nil, err
	} else {
		entry.StartupNotify = val
	}

	if deHeader["OnlyShowIn"] != "" && deHeader["NotShowIn"] != "" {
		return nil, errors.New("only one of either OnlyShowIn or NotShowIn may be specified")
	}

	// non-required fields that don't need validation

	entry.Version = deHeader["Version"]
	entry.GenericName = GetLocalizedValue(deHeader, "GenericName")
	entry.Comment = GetLocalizedValue(deHeader, "Comment")
	entry.Icon = GetLocalizedValue(deHeader, "Icon")
	entry.OnlyShowIn = splitMultiValue(deHeader["OnlyShowIn"])
	entry.NotShowIn = splitMultiValue(deHeader["NotShowIn"])
	entry.TryExec = deHeader["TryExec"]
	entry.Path = deHeader["Path"]
	entry.Actions = splitMultiValue(deHeader["Actions"])
	entry.MimeType = splitMultiValue(deHeader["MimeType"])
	entry.Categories = splitMultiValue(deHeader["Categories"])
	entry.Keywords = splitMultiValue(GetLocalizedValue(deHeader, "Keywords"))
	entry.StartupWMClass = deHeader["StartupWMClass"]

	entry.File = file
	return entry, nil
}

func DesktopEntryTypeIsValid(t string) bool {
	var isValid bool
	for _, typ := range entryTypes {
		if t == typ {
			isValid = true
			break
		}
	}
	return isValid
}

func GetLocalizedValue(cfg map[string] string, property string) string {
	localeRegex, _ := regexp.Compile("^(?P<lang>.+?)(?P<country>_.+?)?(?P<encoding>\\..+?)?(?P<modifier>@.+?)?$")
	locale := localeRegex.FindStringSubmatch(currentLocale)
	if locale == nil {
		println("locale didn't match!")
		return ""
	}

	var lang, country, modifier string = locale[1], locale[2], locale[4]
	if lang == "" {
		return ""
	}

	// test lang_COUNTRY@MODIFIER
	if country != "" && modifier != "" {
		if key, ok := cfg[property + "[" + lang + country + modifier + "]"]; ok {
			return key
		}
	}

	// test lang_COUNTRY
	if country != "" {
		if key, ok := cfg[property + "[" + lang + country + "]"]; ok {
			return key
		}
	}

	// test lang@MODIFIER
	if modifier != "" {
		if key, ok := cfg[property + "[" + lang + modifier + "]"]; ok {
			return key
		}
	}

	// test lang
	if key, ok := cfg[property + "[" + lang + "]"]; ok {
		return key
	}

	// nothing found, use the default
	return cfg[property]
}

// Given a config file group (e.g. "Desktop Entry"), returns a map
// of its values keyed by localization. The default value uses the
// empty string as a key
func GetLocalizedKeys(cfg map[string] string, property string) map[string] string {
	results := make(map[string] string)
	regex, _ := regexp.Compile(property + "\\[(.+)\\]")
	for key, value := range cfg {
		if key == property {
			results[""] = value
		} else if matches := regex.FindStringSubmatch(key); matches != nil {
			results[matches[1]] = value
		}
	}
	return results
}

func GetInstalledApplications() *list.List {
	return GetInstalledApplicationsWhere(func(entry *DesktopEntry) bool {
		return true
	})
}

func GetInstalledApplicationsWhere(f func (*DesktopEntry) bool) *list.List {
	return GetInstalledDesktopEntriesWhere(func (entry *DesktopEntry) bool {
		return entry.Type == "Application" && f(entry)
	})
}

func GetInstalledDesktopEntries() *list.List {
	return GetInstalledDesktopEntriesWhere(func(entry *DesktopEntry) bool {
		return true
	})
}

func GetInstalledDesktopEntriesWhere(f func (*DesktopEntry) bool) *list.List {
	files := GetAllData("applications/*.desktop")
	entries := list.New()
	for _, file := range files {
		entry, err := ParseDesktopEntry(file)
		if err != nil || !f(entry) {
			continue
		}
		entries.PushBack(entry)
	}
	return entries
}

func missingKeyError(key string) error {
	return errors.New("missing required key: " + key)
}

func splitMultiValue(value string) []string {
	regex, _ := regexp.Compile(`[^\\];`)
	values := regex.ReplaceAllString(value, "\n")
	return strings.Split(values, "\n")
}

func getBoolValue(value string) (bool, error) {
	switch value {
		case "true":  return true, nil
		case "false": return false, nil
	}
	return false, errors.New("invalid boolean value: " + value)
}

