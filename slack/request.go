package slack

type IncomingWebhook struct {
	Text            string `json:"text,omitempty"`
	Attachments     []Attachment `json:"attachments,omitempty"`
	// link_names=1 to @ mention people when we ping them
	LinkNames       int `json:"link_names,omitempty"`
	unfurlLinks     bool `json:"unfurl_links,omitempty"`
	IconEmoji       string `json:"icon_emoji,omitempty"`
	RobotName       string `json:"username,omitempty"`
	ChannelWithHash string `json:"channel,omitempty"`
}

type Attachment struct {
	MarkdownIn    []string `json:"mrkdwn_in,omitempty"`
	Fallback      string `json:"fallback,omitempty"`
	ColorHex      string `json:"color,omitempty"`
	AuthorName    string `json:"author_name,omitempty"`
	AuthorIconUrl string `json:"author_icon,omitempty"`
	Fields        []Field `json:"fields,omitempty"`
}

type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool `json:"short,omitempty"`
}


/*
Example data as we tested in the message format explorer: 

{
	"text": "What the Server team did 'yesterday':",
    "attachments": [
        {
            "mrkdwn_in": ["fields"],
            "fallback": "Jonathan merged 4 PRs, and reviewed 6 PRs, created 4 PRs, commented 12x in 4 PRs",
            "color": "blue",
            "author_name": "Jonathan",
            "author_icon": "http://placehold.it/15x15",
            "fields": [
                {
                    "title": "2 Created PRs",
                    "value": "• <www.example.com|title of what the PR was> (0 💬, 1👤, go)\n• <www.example.com|title of what the PR was> (2 💬, 2👤, web)",
                    "short": true
                },
                {
                    "title": "3 Merged PRs",
                    "value": "✓ ~<www.example.com|title of what the PR was> (12 💬, 2 👤, 5 days)~\n✓ ~<www.example.com|title of what the PR was> (1 💬, 1 👤, 8 hours)~\n✓ ~<www.example.com|title of what the PR was> (1 💬, 1 👤, 2 hours)~\n",
                    "short": true
                },
                {
                    "title": "5 new comments in 2 PRs",
                    "value": "• +1 💬 <www.example.com|title of what the PR was> (dge)\n• +4 💬 <www.example.com|title of what the PR was> (alec)\n",
                    "short": true
                },
		        {
                    "title": "2 outstanding PRs",
                    "value": "• <www.example.com|title of what the PR was> (12 💬, 2 👤, 5 days)\n• +4 💬 <www.example.com|title of what the PR was> (0 💬, 0👤, 1 days)\n",
                    "short": true
                }
				
            ]
        },
        {
            "mrkdwn_in": ["fields"],
            "fallback": "SampleName2 merged 4 PRs, and reviewed 6 PRs, created 4 PRs, commented 12x in 4 PRs",
            "color": "purple",
            "author_name": "SampleName2",
            "author_icon": "http://placehold.it/15x15",
            "fields": [
                {
                    "title": "0 Created PRs",
                    "value": "Last created 2 days ago",
                    "short": true
                },
                {
                    "title": "0 Merged PRs",
                    "value": "Last merged 4 days ago",
                    "short": true
                },
                {
                    "title": "0 new comments",
                    "value": "Last commented 8 days ago",
                    "short": true
                },
		        {
                    "title": "0 outstanding PRs",
                    "value": "Last created 12 days ago",
                    "short": true
                }
				
            ]
        }
   ]
}

*/
