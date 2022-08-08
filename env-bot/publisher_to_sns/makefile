run:
	DSN="../environment.sqlite3" go run main.go

build_4_pi:
	CC=arm-linux-gnueabihf-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 go build -o publisher_to_sns
