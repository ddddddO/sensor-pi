from bme280 import BME280
import datetime
import sqlite3

class Repository:
	def __init__(self, dsn) -> None:
		self.conn = sqlite3.connect(dsn)

	def store(self, temperature, pressure, humidity):
		cur = self.conn.cursor()
		date = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
		cur.execute('insert into environment(date, temperature, pressure, humidity) values(?, ?, ?, ?)', (date, temperature, pressure, humidity))
		self.conn.commit()

	def close(self):
		self.conn.close()

if __name__ == '__main__':
	try:
		bme280 = BME280()
		bme280.get_calib_param()
		bme280.readData()
		t, p, h = bme280.result()

		dsn = '/home/pi/github.com/ddddddO/sensor-pi/environment.sqlite3'
		repo = Repository(dsn)
		repo.store(t, p, h)
		repo.close()
	except KeyboardInterrupt:
		pass
