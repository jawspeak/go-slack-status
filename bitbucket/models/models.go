package models

type Author struct {
	Role string `json:"role,omitempty"`
	User *User  `json:"user,omitempty"`
}

type Comment struct {
	Author      *User      `json:"author,omitempty"`
	Comments    []*Comment `json:"comments,omitempty"`
	CreatedDate int64      `json:"createdDate,omitempty"`
	ID          int64      `json:"id,omitempty"`
	Text        string     `json:"text,omitempty"`
	UpdatedDate int64      `json:"updatedDate,omitempty"`
	Version     int64      `json:"version,omitempty"`
}

type CommitsResponse struct {
	IsLastPage    bool  `json:"isLastPage,omitempty"`
	Limit         int64 `json:"limit,omitempty"`
	NextPageStart int64 `json:"nextPageStart,omitempty"`
	Size          int64 `json:"size,omitempty"`
	Start         int64 `json:"start,omitempty"`
}

type Link struct {
	Rel string `json:"rel,omitempty"`
	URL string `json:"url,omitempty"`
}

/** PR stuff */
type PullRequest struct {
	Author      *Author           `json:"author,omitempty"`
	CreatedDate int64             `json:"createdDate,omitempty"`
	Description string            `json:"description,omitempty"`
	ID          int64             `json:"id,omitempty"`
	Links       *PullRequestLinks `json:"links,omitempty"`
	State       string            `json:"state,omitempty"`
	Title       string            `json:"title,omitempty"`
	UpdatedDate int64             `json:"updatedDate,omitempty"`
	Version     int64             `json:"version,omitempty"`
}

type PullRequestLinks struct {
	Self []*PullRequestSelfItems0 `json:"self,omitempty"`
}

type PullRequestSelfItems0 struct {
	Href string `json:"href,omitempty"`
}

type PullRequestActivitiesResponse struct {
	IsLastPage    bool                   `json:"isLastPage,omitempty"`
	Limit         int64                  `json:"limit,omitempty"`
	NextPageStart int64                  `json:"nextPageStart,omitempty"`
	Size          int64                  `json:"size,omitempty"`
	Start         int64                  `json:"start,omitempty"`
	Values        []*PullRequestActivity `json:"values,omitempty"`
}

type PullRequestActivity struct {
	/** COMMENTED, OPENED, APPROVED, RESCOPED, MERGED, DECLINED, UNAPPROVED, some others i think too */
	Action      string   `json:"action,omitempty"`
	Comment     *Comment `json:"comment,omitempty"`
	CreatedDate int64    `json:"createdDate,omitempty"`
	ID          int64    `json:"id,omitempty"`
	User        *User    `json:"user,omitempty"`
}

type PullRequestResponse struct {
	Author      *Author `json:"author,omitempty"`
	CreatedDate int64   `json:"createdDate,omitempty"`
	Description string  `json:"description,omitempty"`
	ID          int64   `json:"id,omitempty"`
	State       string  `json:"state,omitempty"`
	Title       string  `json:"title,omitempty"`
	UpdatedDate int64   `json:"updatedDate,omitempty"`
	Version     int64   `json:"version,omitempty"`
}

type PullRequestsResponse struct {
	IsLastPage    bool           `json:"isLastPage,omitempty"`
	Limit         int64          `json:"limit,omitempty"`
	NextPageStart int64          `json:"nextPageStart,omitempty"`
	Size          int64          `json:"size,omitempty"`
	Start         int64          `json:"start,omitempty"`
	Values        []*PullRequest `json:"values,omitempty"`
}

type User struct {
	DisplayName  string `json:"displayName,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
	ID           int64  `json:"id,omitempty"`
	Link         *Link  `json:"link,omitempty"`
	Name         string `json:"name,omitempty"`
	Slug         string `json:"slug,omitempty"`
}
