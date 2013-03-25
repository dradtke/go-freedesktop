go-freedesktop
==============

`go-freedesktop` is a collection of methods that provide utilities for integrating with freedesktop.org-compliant desktop environments. In particular it simplifies the process of searching for application data and configuration in known locations as defined by freedesktop.org.

Here are a couple examples of what it can do:

Look up Application Data
------------------------

Say your application is called `my-app` and you install an interface file at `/usr/share/my-app/interface.ui`. You can look up its location using the following:

```go
freedesktop.AppName = "my-app"
file := freedesktop.GetAppData("interface.ui")
```

Or say you have a folder `/usr/share/my-app/plugins/` filled with plugins that have an extension of `.so`:

```go
freedesktop.AppName = "my-app"
plugins := freedesktop.GetAllAppData("plugins/*.so")
```

Want to load the configuration file `~/.config/my-app/settings.cfg`?

```go
freedesktop.AppName = "my-app"
cfg := freedesktop.GetAppConfig("settings.cfg")
```

Get a List of Installed Applications
------------------------------------

Want to know what applications are installed on the user's system?

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

Find Application Icon
---------------------

Application icon lookup is also supported.

```go
// this will return something like /usr/share/icons/hicolor/48x48/apps/firefox.png
iconpath := freedesktop.AppIcon("firefox")

// use a method like this for a different size
smallicon := freedesktop.AppIconForSize("firefox", "32x32")
```

Open a URL or File in its Default Application
---------------------------------------------

```go
// open a url in the user's configured default browser
freedesktop.XdgOpen("http://www.my-app.com/")

// or use a file path to open it in its default application
freedesktop.XdgOpen("/path/to/file")
```



This package is very much a work in progress and is being developed ad hoc as I find more things to implement.
