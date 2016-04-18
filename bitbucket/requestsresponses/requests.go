package requestsresponses

/*GetPullRequestActivitiesParams contains all the parameters to send to the API endpoint
for the get pull request activities operation typically these are written to a http.Request
*/
type GetPullRequestActivitiesParams struct {
	/*Limit
	  Probably defaults to 25. It is a best practice to check the limit attribute on the response to see what limit has been applied.
	*/
	Limit         *int64
	Project       string
	PullRequestID int64
	Repo          string
}

type GetPullRequestsParams struct {
	/*Limit
	  Probably defaults to 25. It is a best practice to check the limit attribute on the response to see what limit has been applied.
	*/
	Limit *int64
	/*Order
	  NEWEST is as in newest first.
	*/
	Order   *string
	Project string
	Repo    string
	Role1   *string
	Role2   *string
	/*Start
	  The count of the result to start with, inclusive (I think).
	*/
	Start *int64
	/*State
	  You probably want to include this in, and probably as ALL to see everything (which you won't by default).
	*/
	State     *string
	Username1 *string
	Username2 *string
}
