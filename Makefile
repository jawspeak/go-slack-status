build:
	go build ./...

run: build
	go run main.go -mode=cache -stderrthreshold=INFO -team='Server' -team='Mobile - iOS'
