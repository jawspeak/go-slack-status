package bitbucket
import (
	"github.com/jawspeak/go-slack-status/bitbucket/requestsresponses"
	"net/http"
	"github.com/jawspeak/go-slack-status/config"
	"github.com/golang/glog"
	"encoding/json"
	"github.com/jawspeak/go-slack-status/bitbucket/models"
	"io/ioutil"
	"net/http/httputil"
	"fmt"
	"net/url"
	"bytes"
)

type BitbucketClient struct {
	conf *config.Config
}

func NewBitbucketClient(conf *config.Config) *BitbucketClient {
	return &BitbucketClient{conf: conf}
}

func (c *BitbucketClient) GetPullRequests(params requestsresponses.GetPullRequestsParams) (*requestsresponses.GetPullRequestsOK, error) {
	path := fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/pull-requests", params.Project, params.Repo)
	q := url.Values{}
	if params.Order != nil { q.Set("order", *params.Order) }

	if params.Limit != nil { q.Set("limit", fmt.Sprintf("%d", *params.Limit)) }
	if params.Start != nil { q.Set("start", fmt.Sprintf("%d", *params.Start)) }
	if params.State != nil { q.Set("state", *params.State) }
	if params.Username1 != nil { q.Set("username.1", *params.Username1) }
	if params.Username2 != nil { q.Set("username.2", *params.Username2) }
	if params.Role1 != nil { q.Set("role.1", *params.Role1) }
	if params.Role2 != nil { q.Set("role.2", *params.Role2) }
	if len(q) > 0 {
		path = path + "?" + q.Encode()
	}
	resp, err := c.doRequest(path)
	if err != nil {
		return &requestsresponses.GetPullRequestsOK{}, err
	}
	model := &models.PullRequestsResponse{}
	unmarshal(resp, model)
	return &requestsresponses.GetPullRequestsOK{Payload: model }, nil
}

func (c *BitbucketClient) GetPullRequestActivities(params requestsresponses.GetPullRequestActivitiesParams) (*requestsresponses.GetPullRequestActivitiesOK, error) {
	path := fmt.Sprintf("/rest/api/1.0/projects/%s/repos/%s/pull-requests/%d/activities", params.Project, params.Repo, params.PullRequestID)
	q := url.Values{}
	if params.Limit != nil { q.Set("limit", fmt.Sprintf("%d", *params.Limit)) }
	if len(q) > 0 {
		path = path + "?" + q.Encode()
	}
	resp, err := c.doRequest(path)
	if err != nil {
		return &requestsresponses.GetPullRequestActivitiesOK{}, err
	}
	model := new(models.PullRequestActivitiesResponse)
	unmarshal(resp, &model)
	return &requestsresponses.GetPullRequestActivitiesOK{Payload: model }, nil
}

func (c *BitbucketClient) doRequest(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", "https://" + c.conf.StashHost + path, bytes.NewBuffer(nil))
	if err != nil {
		glog.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.conf.StashUsername, c.conf.StashPassword)

	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Info("Sending request:", string(dump))
	resp, err := http.DefaultClient.Do(req)

	dump, err = httputil.DumpResponse(resp, true)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Info("Response:", string(dump))

	if err != nil {
		glog.Fatal(err)
	}

	if resp.StatusCode != 200 {
		glog.Warning("non-200 response code ", resp)
		return &http.Response{}, requestsresponses.NewAPIError("todo <name me>", resp, resp.StatusCode)
	}

	return resp, nil
}

func unmarshal(resp *http.Response, model interface{}) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Fatal(err)
	}
	err = json.Unmarshal(b, model)
	if err != nil {
		glog.Fatal(err)
	}
}
