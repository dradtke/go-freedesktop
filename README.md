go-freedesktop
==============

`go-freedesktop` is a collection of methods that provide utilities for integrating with freedesktop.org-compliant desktop environments. In particular it simplifies the process of searching for application data and configuration in known locations as defined by freedesktop.org.

Here are a couple examples of what it can do:

Get a List of Installed Applications
------------------------------------

```go
applications := freedestkop.GetInstalledApplications()
for e := applications.Front(); e != nil; e = e.Next() {
	fmt.Println(e.Value.(*freedesktop.DesktopEntry).Name)
}
```

Or if you want to get a list of installed applications that support a certain file type:

```go
applications := freedestkop.GetInstalledApplicationsWhere(func(entry *freedesktop.DesktopEntry) bool {
	for _, mime := range entry.MimeType {
		if strings.HasPrefix(mime, "audio/") {
			return true
		}
	}
	return false
})
for e := applications.Front(); e != nil; e = e.Next() {
	fmt.Println(e.Value.(*freedesktop.DesktopEntry).Name)
}
```

Look up Application Configuration
---------------------------------

You can use this method to locate a configuration file, using the standard freedesktop.org search algorithm:

```go
filepath := freedesktop.FindInConfig("my-app.cfg")
```

Find Application Icon
---------------------

Want the absolute path of an application's icon?

```go
// this will return something like /usr/share/icons/hicolor/48x48/apps/firefox.png
iconpath := freedesktop.AppIcon("firefox")
// use a method like this for a different size
smallicon := freedesktop.AppIconForSize("firefox", "32x32")
```

This package is very much a work in progress and is being developed ad hoc as I find more things to implement.
