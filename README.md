# KRunner GitLab backend

This package provides a [KRunner](http://blog.davidedmundson.co.uk/blog/cross-process-runners/) backend which will use a GitLab instance as a search backend. Currently only project names are searched for, this might change in the future.

## Installation

Clone this repository

```
git clone https://github.com/shochdoerfer/krunner-gitlab
```

Build the go application

```
cd
go install
```

Register the runner in KDE by storing a file called `krunner-gitlab.desktop` in `$HOME/.local/share/kservices5` and then restart the rkunner process.

```
[Desktop Entry]
Name=GitLab
Comment=GitLab KRunner
X-KDE-ServiceTypes=Plasma/Runner
Type=Service
Icon=internet-web-browser
X-KDE-PluginInfo-Author=Stephan Hochdörfer
X-KDE-PluginInfo-Email=S.Hochdoerfer@bitExpert.de
X-KDE-PluginInfo-Name=krunner-gitlab
X-KDE-PluginInfo-Version=1.0
X-KDE-PluginInfo-License=Apache-2.0
X-KDE-PluginInfo-EnabledByDefault=true
X-Plasma-API=DBus
X-Plasma-DBusRunner-Service=de.hochdoerfer.gitlab
X-Plasma-DBusRunner-Path=/krunner
```

To configure krunner-gitlab with the url and the access token for your own GitLab instance, create a file `$HOME/.krunner-gitlab.yaml` like this:

```
url: https://your-gitlab-server/api/v4
token: your-token
```

It is important to note that the url needs to point to the GitLab API url!

## Run the application

Run `krunner-gitlab` in your `go/bin` directory. Invoke KRunner and start searching for GitLab projects. 

## License

KRunner GitLab is released under the Apache 2.0 license.
