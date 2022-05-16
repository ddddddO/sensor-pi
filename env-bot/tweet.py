#!/usr/bin/env python3

import settings
import tweepy
import sqlite3
from pypika import Query, Table, Order

dsn = '/home/pi/github.com/ddddddO/sensor-pi/env-bot/environment.sqlite3'
conn = sqlite3.connect(dsn)
environment = Table('environment')
q = Query.from_(environment)\
    .select(environment.temperature, environment.pressure, environment.humidity)\
    .orderby(environment.date, order=Order.desc)\
    .limit(1)

cur = conn.cursor()
cur.execute(str(q))
row = cur.fetchone()
conn.close()

title = 'ただいまの気温・気圧・湿度(屋内)'
location = '@多摩川付近'
temp = "temp : {:-6.2f} ℃\n".format(row[0])
pressure = "pressure : %7.2f hPa\n" % (row[1])
hum = "hum : %6.2f ％" % (row[2])
content = title + location + '\n' + temp + pressure + hum

auth = tweepy.OAuthHandler(settings.consumer_key, settings.consumer_secret)
auth.set_access_token(settings.token, settings.token_secret)
api = tweepy.API(auth)

image_path = '/home/pi/github.com/ddddddO/sensor-pi/env-bot/plotter/pressure.png'
api.update_status_with_media(status=content, filename=image_path)
