package main

import (
	"github.com/jawspeak/go-slack-status/bitbucket"
	"github.com/jawspeak/go-slack-status/config"
	"github.com/jawspeak/go-slack-status/slack"

	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/jawspeak/go-slack-status/bitbucket/cache"
	"os"
	"strings"
	"time"
)

type teams []string

func (t *teams) String() string {
	return fmt.Sprint(*t)
}
func (teams *teams) Set(value string) error {
	for _, t := range strings.Split(value, ",") {
		*teams = append(*teams, t)
	}
	return nil
}

// Works well if you have a crontab and want different teams to run at different times, and some to use the cache
// while others to re-fetch the data. This string matches the name it the config
// You can pass in multiple "-team abc -team def" etc to run multiple teams at once.
var teamsFlag teams

func init() {
	flag.Var(&teamsFlag, "team", "(multiple allowed) One <stats_for_these_teams.team_name> from config.json to only send out for those")
}

func main() {
	var modeFlag = flag.String("mode", "", "[cached|live] to use the local cache of data or re-fetch it. if [live] it will exit after fetching.")
	var confFilePath = flag.String("conf", "./config.json", "path to config.json")
	flag.Parse()
	conf := config.Setup(confFilePath)
	glog.Info("starting. loaded config")

	var cachedata *cache.Data
	switch *modeFlag {
	case "cache":
		cachedata = cache.LoadJson("./pr-cached-data.json")
	case "live":
		cachedata = bitbucket.NewFetcher(conf).FetchData()
		glog.Info("loaded PR activity from stash")
		// Fetching is slow, so store for playing with it later.
		cachedata.SaveJson("pr-cached-data.json")
		glog.Warning("FINISHED saving data successfully. Exiting. Use -mode=cache to run again")
		os.Exit(0)
	case "":
		fallthrough
	default:
		glog.Error("You need to pass in -mode")
		flag.Usage()
		os.Exit(1)
	}

	slackClient := slack.NewSlackClient(conf, cachedata)
	if len(teamsFlag) == 0 {
		runAllTeams(conf, slackClient, cachedata)
	} else {
		teamsMap := conf.TeamsMap()
		for _, requestedTeamFlag := range teamsFlag {
			if _, ok := teamsMap[requestedTeamFlag]; !ok {
				glog.Errorf("Requested team (%s) not found (refer to config.json for teams)", requestedTeamFlag)
				flag.Usage()
				os.Exit(1)
			}
		}
		for _, requestedTeamFlag := range teamsFlag {
			slackClient.PingSlackWebhook(0, teamsMap[requestedTeamFlag])
		}
	}
	glog.Info("work done\n")
}

func runAllTeams(conf *config.Config, slackClient *slack.SlackClient, cachedata *cache.Data) {
	for i, team := range conf.StatsForTheseTeams {
		slackClient.PingSlackWebhook(i, team)
		time.Sleep(time.Second * 30) // ping each team one at a time
	}
}
