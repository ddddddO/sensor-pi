#!/usr/bin/env python3

from bme280 import BME280
import datetime
import sqlite3

class Repository:
	def __init__(self, dsn, file) -> None:
		self.conn = sqlite3.connect(dsn)
		self.file = file

	def store(self, temperature, pressure, humidity):
		cur = self.conn.cursor()
		date = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
		cur.execute('insert into environment(date, temperature, pressure, humidity) values(?, ?, ?, ?)', (date, temperature, pressure, humidity))
		self.conn.commit()

	def close(self):
		self.conn.close()

	def storeFile(self, temperature, pressure, humidity):
		f = open(self.file, 'w')
		f.write("temp : {:-6.2f} ℃\n".format(temperature))
		f.write("pressure : %7.2f hPa\n" % (pressure))
		f.write("hum : %6.2f ％" % (humidity))
		f.close()

if __name__ == '__main__':
	try:
		bme280 = BME280()
		bme280.get_calib_param()
		bme280.readData()
		t, p, h = bme280.result()

		dsn = '/home/pi/github.com/ddddddO/sensor-pi/env-bot/environment.sqlite3'
		file = '/tmp/bme_result'
		repo = Repository(dsn, file)
		repo.store(t, p, h)
		repo.storeFile(t, p, h)
		repo.close()
	except KeyboardInterrupt:
		pass
