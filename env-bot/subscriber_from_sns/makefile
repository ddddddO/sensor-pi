# build:
# 	GOOS=linux GOARCH=amd64 go build -o main

# zip: build
# 	rm -f subscriber_from_sns.zip && zip -r subscriber_from_sns.zip main

zip_all: zip_humidity zip_pressure zip_temperature

build_humidity:
	cd humidity/ && GOOS=linux GOARCH=amd64 go build -o main

zip_humidity: build_humidity
	cd humidity/ && rm -f subscriber_from_sns_humidity.zip && zip -r subscriber_from_sns_humidity.zip main

build_pressure:
	cd pressure/ && GOOS=linux GOARCH=amd64 go build -o main

zip_pressure: build_pressure
	cd pressure/ && rm -f subscriber_from_sns_pressure.zip && zip -r subscriber_from_sns_pressure.zip main

build_temperature:
	cd temperature/ && GOOS=linux GOARCH=amd64 go build -o main

zip_temperature: build_temperature
	cd temperature/ && rm -f subscriber_from_sns_temperature.zip && zip -r subscriber_from_sns_temperature.zip main
