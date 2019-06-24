### Post to Slack Atlassian Bitbucket (Stash) status updates

No significant external dependencies.

Check out your config repo into the `config-repo` folder which is ignored by this git repo. Then your team can have an internal private repo for collaborating crontabs and what teams get what notifications.

### Development

I work with a separate GOPATH for open source work. In the folder holding the gopath src file, I have a GOPATH_source_me.sh script with the following contents:

```bash
echo "you should source me, not execute me! run >> 'source GOPATH_source_me'"
export GOPATH=/Users/jaw/me/Development/go-opensource-work
export PATH=$PATH:$GOPATH/bin
echo $GOPATH
```

That topmost folder contains:

```
>ls -a /Users/jaw/me/Development/go-opensource-work
.                       bin
..                      pkg
.idea                   src
GOPATH_source_me.sh     
```

When I want to start working on this, I source that `. GOPATH_source_me.sh`.

I use IntelliJ with the golang plugin.

I run `go install ./...` to install this binary in bin/ (and see my PATH edit above so it would include that bin/ folder).

See the `Makefile` for generating, building, and running.

I scp this script to my server I eventually run this on. See the crontab work in the capistrano script.

### Here's what a message looks like from this slackbot
As shown via the slack message [previewer](https://github.com/jawspeak/go-slack-status/blob/master/slack/request.go)

![sample message](sample%20slackbot%20msg.png)
