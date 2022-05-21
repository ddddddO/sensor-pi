update_dagu:
	go install github.com/yohamta/dagu/cmd@latest && mv /home/pi/go/bin/cmd /home/pi/go/bin/dagu

# crontab -e
# 0 9,18 * * * /home/pi/go/bin/dagu start /home/pi/github.com/ddddddO/sensor-pi/env-bot/dag.yaml 1> /dev/null
tweet:
	dagu start dag.yaml

# Create database: sqlite3 environment.sqlite3
# Create table: execute schema.sql TODO: using migration tool
conn:
	sqlite3 environment.sqlite3
