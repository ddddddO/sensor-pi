# sudo apt-get update
# sudo apt install i2c-tools
i2c:
	i2cdetect -y 1

# sudo apt-get update
# sudo apt install -y python-smbus
# sudo pip install smbus2
bme280:
	 sudo python bme280_python2.py

# sudo pip3 install mh_z19
mhz19:
	sudo python3 -m mh_z19
