#!/usr/local/bin/python3.8

import sys
sys.path.append('/home/pi/.local/lib/python3.8/site-packages')

from pypika import Query, Table, Order
import mh_z19
import os
import datetime
import sqlite3

class Repository:
	def __init__(self, dsn) -> None:
		self.conn = sqlite3.connect(dsn)

	def fetch(self) -> float:
		mh_z19 = Table('mh_z19')
		q = Query.from_(mh_z19)\
			.select(mh_z19.co2)\
			.orderby(mh_z19.date, order=Order.desc)\
			.limit(1)

		cur = self.conn.cursor()
		cur.execute(str(q))
		row = cur.fetchone()
		c = row[0]
		return c

	def store(self, co2):
		date = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
		mh_z19 = Table('mh_z19')
		q = Query.into(mh_z19)\
			.columns(mh_z19.date, mh_z19.co2)\
			.insert(date, co2)

		cur = self.conn.cursor()
		cur.execute(str(q))
		self.conn.commit()

	def close(self):
		self.conn.close()


ret = mh_z19.read_from_pwm()
co2 = ret['co2'] # TODO: KeyError handling. retryする？

dsn = os.getenv('DSN')
repo = Repository(dsn)

repo.store(co2)