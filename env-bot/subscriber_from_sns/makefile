build:
	GOOS=linux GOARCH=amd64 go build -o main

zip: build
	rm subscriber_from_sns.zip && zip -r subscriber_from_sns.zip main