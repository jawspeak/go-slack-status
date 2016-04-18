build:
	go build ./...

run: build
	go run main.go -mode=cache -stderrthreshold=INFO -team='Server' -team='Mobile - iOS'

_build_linux:
	env GOOS=linux GOARCH=amd64 go build -v github.com/jawspeak/go-slack-status
deploy: _build_linux
	scp go-slack-status projectmonitor-a-01.corp.squareup.com:slack-status-tool/go-slack-status
	scp config/deploy/crontab projectmonitor-a-01.corp.squareup.com:slack-status-tool/crontab.tmp
	scp config.json projectmonitor-a-01.corp.squareup.com:slack-status-tool
	ssh projectmonitor-a-01.corp.squareup.com 'crontab -r; true'
	ssh projectmonitor-a-01.corp.squareup.com 'crontab slack-status-tool/crontab.tmp && rm -f slack-status-tool/crontab.tmp'
	rm ./go-slack-status
