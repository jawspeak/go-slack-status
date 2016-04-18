package config

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	filepath "path/filepath"
	"regexp"
	"strings"
)

type Config struct {
	StashHost              string          `json:"stash_host"`
	StashUsername          string          `json:"stash_username"`
	StashPassword          string          `json:"stash_password"`
	Projects               []projectConfig `json:"stash_projects"`
	LookBackDays           int             `json:"look_back_days"`
	StatsForTheseTeams     []Team          `json:"stats_for_these_teams"`
	IgnoreCommentUsernames []string        `json:"ignore_comments_from_usernames"`
}
type Team struct {
	TeamName                   string   `json:"team_name"`
	SlackIncomingWebhookUrl    string   `json:"slack_incoming_webhook_url"`
	SlackNotifyPeopleOnPosting bool     `json:"slack_notify_people_on_posting"`
	SlackRobotName             string   `json:"robot_name"`
	SlackRobotEmoji            string   `json:"robot_emoji"`
	SlackChannelOverride       string   `json:"slack_channel_override"`
	Members                    []string `json:"members_ldap_names"`
}
type projectConfig struct {
	Project string   `json:"project"`
	Repos   []string `json:"repos"`
}

func (c *Config) TeamsMap() map[string]Team {
	teamsMap := make(map[string]Team)
	for _, t := range c.StatsForTheseTeams {
		teamsMap[t.TeamName] = t
	}
	return teamsMap
}

func (c Config) AllMembersLdaps() []string {
	ldaps := make(map[string]bool)
	for _, team := range c.StatsForTheseTeams {
		for _, m := range team.Members {
			ldaps[m] = true
		}
	}
	keys := make([]string, 0, len(ldaps))
	for key, _ := range ldaps {
		keys = append(keys, key)
	}
	return keys
}

func validateRequiredField(field string, configValue *string) {
	if configValue == nil || len(*configValue) == 0 {
		fmt.Println("Required field unset in config.json: ", field)
		os.Exit(1)
	}
}

func validateNoCommasInTeamNames(teams *[]Team) {
	for _, t := range *teams {
		if strings.Contains(t.TeamName, ",") {
			fmt.Println("Commas not allowed in TeamNames in config.json: ", t.TeamName)
			os.Exit(1)
		}
	}
}

func parseJsonFileStripComments(path string, conf interface{}) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		glog.Fatal(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Println("read file: ", path, " working dir: ", wd, " absolute path: ", abspath)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		glog.Fatal(err)
	}
	commentStripper := regexp.MustCompile("(?s)[^:]//.*?\n|/\\*.*?\\*/")
	file = commentStripper.ReplaceAll(file, nil)

	if err := json.Unmarshal(file, &conf); err != nil {
		glog.Fatal(err)
	}
	glog.Info("comment filtered config.json file contents:", string(file))
}

func Setup(pathToConfig *string) *Config {
	var conf Config
	parseJsonFileStripComments(*pathToConfig, &conf)
	validateRequiredField("host", &conf.StashHost)
	validateRequiredField("username", &conf.StashUsername)
	validateNoCommasInTeamNames(&conf.StatsForTheseTeams)
	if &conf.StashPassword == nil || len(conf.StashPassword) == 0 {
		fmt.Printf("Enter your password for %s: ", conf.StashUsername)
		bytePassword, err := terminal.ReadPassword(0)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println()
		conf.StashPassword = string(bytePassword)
	}
	validateRequiredField("password", &conf.StashPassword)

	fmt.Println("using stash host:", conf.StashHost)
	return &conf
}
