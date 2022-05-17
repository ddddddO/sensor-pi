import datetime
import sqlite3
from pypika import Query, Table, Order

class Repository:
	def __init__(self, dsn) -> None:
		self.conn = sqlite3.connect(dsn)

	def fetch(self) -> (float, float, float):
		environment = Table('environment')
		q = Query.from_(environment)\
			.select(environment.temperature, environment.pressure, environment.humidity)\
			.orderby(environment.date, order=Order.desc)\
			.limit(1)

		cur = self.conn.cursor()
		cur.execute(str(q))
		row = cur.fetchone()
		t, p, h = row[0], row[1], row[2]
		return t, p, h

	def store(self, temperature, pressure, humidity):
		date = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
		environment = Table('environment')
		q = Query.into(environment)\
			.columns(environment.date, environment.temperature, environment.pressure, environment.humidity)\
			.insert(date, temperature, pressure, humidity)

		cur = self.conn.cursor()
		cur.execute(str(q))
		self.conn.commit()

	def close(self):
		self.conn.close()
