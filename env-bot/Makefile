# MEMO: github.com/yohamta/dagu/.github/workflows/release.yaml を参考にコマンド実行後、GOOS=linux GOARCH=arm GOARM=7 go build -o dagu の後に、ラズパイへバイナリを転送する
# update_dagu:
# 	go install github.com/yohamta/dagu/cmd@latest && mv /home/pi/go/bin/cmd /home/pi/go/bin/dagu

# crontab -e
# 0 9,18 * * * /home/pi/go/bin/dagu start /home/pi/github.com/ddddddO/sensor-pi/env-bot/dag.yaml 1> /dev/null
tweet:
	dagu start dag.yaml

# Create database: sqlite3 environment.sqlite3
# Create table: execute schema.sql TODO: using migration tool
conn:
	sqlite3 environment.sqlite3

receive_bt:
	python3 bluetooth_receiver.py &

# mh_z19の動作確認
test_mhz_uart:
	sudo python3 -m mh_z19

test_mhz_pwm:
	python3 -m mh_z19 --pwm --pwm_gpio 12

test_bme:
	sudo python3 ../bme280.py
