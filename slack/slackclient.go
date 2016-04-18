package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/jawspeak/go-slack-status/bitbucket/cache"
	"github.com/jawspeak/go-slack-status/config"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type SlackClient struct {
	conf                  *config.Config
	cachedata             *cache.Data
	now                   time.Time
	yesterdayStart        time.Time
	yesterdayEnd          time.Time
	locationOfServerTimes *time.Location
}

var colorPallet = [2][1]string{
	//  [...]string{"#d9850e", "#680b0b", "#eaabf3", "#222455", "#3fe6d6"}, // http://www.color-hex.com/color-palette/17537
	//	[3]string{"#90a7b4", "#a3c4d8", "#bcdce5"}, // http://www.color-hex.com/color-palette/17727
	//	[3]string{"#d5f0f4", "#82d2de", "#30b4c9"}, // http://www.color-hex.com/color-palette/17675
	//	playing with various color schemes. here i'll alternate with two greys.
	[1]string{"#ddd"},
	[1]string{"#666"},
}

const bullet = "• "
const check = "✓ ~"
const MERGED = "MERGED"
const SHORT_MMM_D = "Jan 2"

func NewSlackClient(conf *config.Config, cacheData *cache.Data) *SlackClient {
	now := time.Now()
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		glog.Fatal(err)
	}
	var delta time.Duration
	// Expect this will run in a crontab in the mornings pacific time, ignore holidays.
	if now.In(loc).Weekday() == time.Monday {
		delta, err = time.ParseDuration("-72h")
	} else {
		delta, err = time.ParseDuration("-24h")
	}
	if err != nil {
		glog.Fatal(err)
	}
	yesterday := now.Add(delta)
	yesterdayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, loc)
	yesterdayEnd := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 0, 0, loc)

	return &SlackClient{
		conf:                  conf,
		cachedata:             cacheData,
		now:                   now,
		yesterdayStart:        yesterdayStart,
		yesterdayEnd:          yesterdayEnd,
		locationOfServerTimes: loc, // our bitbucket server is in pacific time
	}
}

func (c *SlackClient) PingSlackWebhook(i int, team config.Team) {
	glog.Infof("pinging on slack team: %s", team.TeamName)
	wh := c.buildRequestForTeam(i, team)

	b, err := json.Marshal(wh)
	if err != nil {
		glog.Fatal(err)
	}

	glog.Infof("Sending request, with jsonbody=%s", string(b))
	req, err := http.NewRequest("POST", team.SlackIncomingWebhookUrl, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glog.Fatal(err)
	}
	defer resp.Body.Close()
	glog.Infof("Response: status=%s, headers=%s", resp.Status, resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Infof("Response: body=%s", string(body))
}

func (c *SlackClient) buildRequestForTeam(i int, team config.Team) IncomingWebhook {
	mrkdn := make([]string, 0)
	mrkdn = append(mrkdn, "fields")
	attachments := make([]Attachment, 0)

	for j, ldapName := range team.Members {
		// "yesterday's" work
		createdPrs := make(map[int64]cache.PullRequest)
		mergedPrs := make(map[int64]cache.PullRequest)
		outstandingPrs := make(map[int64]cache.PullRequest)
		prsCommentedIn := make(map[int64]cache.PullRequest)
		commentsInPrs := make(map[int64]cache.PrInteraction)

		for _, pr := range c.cachedata.PullRequests {
			if pr.AuthorLdap == ldapName {
				prCreated := time.Unix(pr.CreatedDateTime, 0)
				if prCreated.After(c.yesterdayStart) && prCreated.Before(c.yesterdayEnd) {
					createdPrs[pr.PullRequestId] = pr
				}
				prUpdated := time.Unix(pr.UpdatedDateTime, 0)
				if pr.State == "MERGED" && prUpdated.After(c.yesterdayStart) && prUpdated.Before(c.yesterdayEnd) {
					mergedPrs[pr.PullRequestId] = pr
				}
				if pr.State == "OPEN" {
					outstandingPrs[pr.PullRequestId] = pr
				}
			}
			for _, prInteraction := range pr.Comments {
				interactionCreated := time.Unix(prInteraction.CreatedDateTime, 0)
				if prInteraction.AuthorLdap == ldapName && interactionCreated.After(c.yesterdayStart) && interactionCreated.Before(c.yesterdayEnd) {
					prsCommentedIn[pr.PullRequestId] = pr
					commentsInPrs[prInteraction.RefId] = prInteraction
				}
			}
		}

		fields := make([]Field, 0)
		if len(createdPrs) == 0 && len(mergedPrs) == 0 && len(outstandingPrs) == 0 && len(prsCommentedIn) == 0 {
			fields = append(fields, Field{
				Value: "No Detected Activity (report to jaw@ as it may be a bug)",
				Short: true,
			})
		} else {
			c.addCreatedPrs(&fields, &createdPrs)
			c.addMergedPrs(&fields, &mergedPrs)
			c.addComments(&fields, &commentsInPrs, &prsCommentedIn)
			c.addOutstandingPrs(&fields, &outstandingPrs)
			// TODO can indicate the PRs each person Approved, too
		}

		colorI := i % len(colorPallet)
		colorJ := j % len(colorPallet[colorI])
		attachments = append(attachments, Attachment{
			MarkdownIn: mrkdn,
			Fallback: fmt.Sprintf("%s: %d created, %d merged, commented %dx in %d PRs, (%d outstanding)",
				ldapName, len(createdPrs), len(mergedPrs), len(commentsInPrs), len(prsCommentedIn), len(outstandingPrs)),
			ColorHex:      colorPallet[colorI][colorJ],
			AuthorName:    ldapName,
			AuthorIconUrl: "",
			Fields:        fields,
		})
	}

	linkNamesLookup := make(map[bool]int)
	linkNamesLookup[true] = 1
	return IncomingWebhook{
		Text:            fmt.Sprintf("▼ ▼ ▼ What *%s* team did yesterday (%s) ▼ ▼ ▼ Virtual standup %s", team.TeamName, c.yesterdayStart.Format(SHORT_MMM_D), "@jaw"),
		Attachments:     attachments,
		LinkNames:       linkNamesLookup[team.SlackNotifyPeopleOnPosting],
		unfurlLinks:     false,
		IconEmoji:       team.SlackRobotEmoji,
		RobotName:       team.SlackRobotName,
		ChannelWithHash: team.SlackChannelOverride,
	}
}

func (c *SlackClient) addCreatedPrs(fields *[]Field, createdPrs *map[int64]cache.PullRequest) {
	const createdFmt = "<%s|%s> %s/%s"
	//%d :speech_balloon:, %d :bust_in_silhouette:,

	value := make([]string, 0)
	for _, e := range *createdPrs {
		var buff bytes.Buffer
		if e.State == MERGED {
			buff.WriteString(check)
		} else {
			buff.WriteString(bullet)
		}
		buff.WriteString(fmt.Sprintf(createdFmt, e.SelfUrl, elipses(e.Title), e.Repo, e.Project))
		if e.State == MERGED {
			buff.WriteString("~")
		}
		value = append(value, buff.String())
	}
	*fields = append(*fields, Field{
		Title: fmt.Sprintf("%d Created", len(*createdPrs)),
		Value: strings.Join(value, "\n"),
		Short: true,
	})
}

func (c *SlackClient) addMergedPrs(fields *[]Field, mergedPrs *map[int64]cache.PullRequest) {
	const mergedFmt = "<%s|%s> %d :speech_balloon:, %d :bust_in_silhouette:, %s/%s"

	value := make([]string, 0)
	for _, e := range *mergedPrs {
		var buff bytes.Buffer
		if e.State == MERGED {
			buff.WriteString(check)
		} else {
			buff.WriteString(bullet)
		}
		buff.WriteString(fmt.Sprintf(mergedFmt, e.SelfUrl, elipses(e.Title), e.CommentCount,
			len(e.CommentsByAuthorLdap), e.Repo, e.Project))
		if e.State == MERGED {
			buff.WriteString("~")
		}
		value = append(value, buff.String())
	}
	*fields = append(*fields, Field{
		Title: fmt.Sprintf("%d Merged", len(*mergedPrs)),
		Value: strings.Join(value, "\n"),
		Short: true,
	})
}

func (c *SlackClient) addComments(fields *[]Field, commentsInPrs *map[int64]cache.PrInteraction,
	prsCommentedIn *map[int64]cache.PullRequest) {
	const commentsFmt = "<%s|%s> +%d :speech_balloon:, %s/%s"

	value := make([]string, 0)
	for _, e := range *prsCommentedIn {
		var buff bytes.Buffer
		if e.State == MERGED {
			buff.WriteString(check)
		} else {
			buff.WriteString(bullet)
		}

		buff.WriteString(fmt.Sprintf(commentsFmt, e.SelfUrl, elipses(e.Title), len(*commentsInPrs),
			e.Repo, e.Project))
		if e.State == MERGED {
			buff.WriteString("~")
		}
		value = append(value, buff.String())
	}
	var title string
	if len(*commentsInPrs) > 0 {
		title = fmt.Sprintf("%d Comments in %d PRs", len(*commentsInPrs), len(*prsCommentedIn))
	} else {
		title = "0 Comments"
	}
	*fields = append(*fields, Field{
		Title: title,
		Value: strings.Join(value, "\n"),
		Short: true,
	})
}

func (c *SlackClient) addOutstandingPrs(fields *[]Field, outstandingPrs *map[int64]cache.PullRequest) {
	const outstandingFmt = "<%s|%s> %d :speech_balloon:, %d :bust_in_silhouette:, %s/%s (%s days old)"

	value := make([]string, 0)
	for _, e := range *outstandingPrs {
		var buff bytes.Buffer
		buff.WriteString(bullet)
		// TODO also verify no unapprovals, or look which is latest.
		// if len(e.ApprovalsByAuthorLdap) > 0 {
		// buff.WriteString(":white_check_mark: ") // PR is approved, needs merging
		// }
		days := fmt.Sprintf("%.1f", time.Now().Sub(time.Unix(e.CreatedDateTime, 0)).Hours()/24)
		buff.WriteString(fmt.Sprintf(outstandingFmt, e.SelfUrl, elipses(e.Title), e.CommentCount,
			len(e.CommentsByAuthorLdap), e.Repo, e.Project, days))
		value = append(value, buff.String())
	}
	*fields = append(*fields, Field{
		Title: fmt.Sprintf("%d Outstanding", len(*outstandingPrs)),
		Value: strings.Join(value, "\n"),
		Short: true,
	})
}

func elipses(title string) string {
	const maxLen = 40
	if len(title) > maxLen {
		return title[:maxLen] + "…"
	}
	return title
}
