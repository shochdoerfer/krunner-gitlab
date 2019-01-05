package main

import (
	"os/exec"
	"fmt"
	"time"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/xanzy/go-gitlab"
	"github.com/spf13/viper"
)

// see for definition https://cgit.kde.org/krunner.git/plain/src/data/org.kde.krunner1.xml
const intro = `
<node>
  <interface name="org.kde.krunner1">
    <method name="Actions">
      <annotation name="org.qtproject.QtDBus.QtTypeName.Out0" value="RemoteActions" />
      <arg name="matches" type="a(sss)" direction="out" />
    </method>
    <method name="Run">
      <arg name="matchId" type="s" direction="in"/>
      <arg name="actionId" type="s" direction="in"/>
    </method>
    <method name="Match">
      <arg name="query" type="s" direction="in"/>
      <annotation name="org.qtproject.QtDBus.QtTypeName.Out0" value="RemoteMatches"/>
      <arg name="matches" type="a(sssuda{sv})" direction="out"/>
    </method>
  </interface>` + introspect.IntrospectDataString + `</node>`

type matchOut struct {
	ID, Text, IconName string
	Type               int32
	Relevance          float64
	Properties         map[string]interface{}
}

// http://blog.davidedmundson.co.uk/blog/cross-process-runners/
type runner struct {
	client *gitlab.Client
}

func (r runner) Actions() ([]string, *dbus.Error) {
	return make([]string, 0), nil
}

func (r runner) Match(query string) ([]matchOut, *dbus.Error) {
	matches := make([]matchOut, 0)
	
	if len(query) < 4 {
		return matches, nil
	}

	opt := &gitlab.ListProjectsOptions{
		Search: gitlab.String(query), 
		Simple: gitlab.Bool(true),
		Archived: gitlab.Bool(false),
		Statistics: gitlab.Bool(false), 
		ListOptions: gitlab.ListOptions{
			PerPage: 5,
			Page:    1,
		},
	}
	
	projects, _, err := r.client.Projects.ListProjects(opt)
	if err != nil {
		return matches, nil
	}
	
	matches = make([]matchOut, len(projects))
	time := time.Now()
	for i := 0; i < len(projects); i++ {
		// the longer the last acitivity has been, the least important is the project
		relevance := (float64(time.Unix()) / float64(projects[i].LastActivityAt.Unix()))
		if relevance > 1 {
			relevance = 1 - ((relevance - float64(1)) * float64(10))
		}
		
		matches[i].ID = projects[i].WebURL
		matches[i].Text = projects[i].Name
		matches[i].IconName = "internet-web-browser"
		matches[i].Type = 100
		matches[i].Relevance = relevance
	}

	return matches, nil
}

func (r runner) Run(matchId string, actionId string) *dbus.Error {
	cmd := "x-www-browser"
	args := []string{matchId}
	exec.Command(cmd, args...).Run();
	
	return nil
}

func main() {
	// read config
	viper.SetConfigName(".krunner-gitlab")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	
	// configure GitLab client connection
	client := gitlab.NewClient(nil, fmt.Sprintf("%v", viper.Get("token")))
	client.SetBaseURL(fmt.Sprintf("%v", viper.Get("url")))

	// connect to Session DBUS
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	reply, err := conn.RequestName("de.hochdoerfer.gitlab", dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		panic("Name de.hochdoerfer.gitlab already taken")
	}

	// create & export runner instance
	f := runner{client}
	conn.Export(f, "/krunner", "org.kde.krunner1")
	conn.Export(introspect.Introspectable(intro), "/krunner", "org.freedesktop.DBus.Introspectable")
	fmt.Println("Listening on de.hochdoerfer.gitlab/krunner...")
	select {}
}
