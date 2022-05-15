# MEMO:
# crontab -e
# 0 9,18 * * * sudo /home/pi/github.com/ddddddO/sensor-pi/env-bot/bme280_store.py > /tmp/bme_result && plotter && cat /tmp/bme_result | /home/pi/github.com/ddddddO/sensor-pi/env-bot/tweet.py
tweet:
	sudo ./bme280_store.py > /tmp/bme_result && plotter && cat /tmp/bme_result | ./tweet.py

# Create database: sqlite3 environment.sqlite3
# Create table: execute schema.sql TODO: using migration tool
conn:
	sqlite3 environment.sqlite3
