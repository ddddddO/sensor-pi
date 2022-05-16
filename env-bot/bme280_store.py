#!/usr/bin/env python3

from bme280 import BME280
import datetime
import sqlite3
from pypika import Query, Table

class Repository:
	def __init__(self, dsn) -> None:
		self.conn = sqlite3.connect(dsn)

	def store(self, temperature, pressure, humidity):
		cur = self.conn.cursor()
		date = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')

		environment = Table('environment')
		q = Query.into(environment)\
			.columns(environment.date, environment.temperature, environment.pressure, environment.humidity)\
			.insert(date, temperature, pressure, humidity)

		cur.execute(str(q))
		self.conn.commit()

	def close(self):
		self.conn.close()

if __name__ == '__main__':
	try:
		bme280 = BME280()
		bme280.get_calib_param()
		bme280.read_data()
		t, p, h = bme280.result()

		dsn = '/home/pi/github.com/ddddddO/sensor-pi/env-bot/environment.sqlite3'
		repo = Repository(dsn)
		repo.store(t, p, h)
		repo.close()
	except KeyboardInterrupt:
		pass
