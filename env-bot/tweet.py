#!/usr/bin/env python3

import settings
import tweepy

path = '/tmp/bme_result'
f = open(path, 'r')
lines = f.readlines()
f.close()

title = 'ただいまの気温・気圧・湿度(屋内)'
location = '@多摩川付近'
temp = lines[0]
pressure = lines[1]
hum = lines[2]
content = title + location + '\n' + temp + '\n' + pressure + '\n' + hum

auth = tweepy.OAuthHandler(settings.consumer_key, settings.consumer_secret)
auth.set_access_token(settings.token, settings.token_secret)
api = tweepy.API(auth)

image_path = '/home/pi/github.com/ddddddO/sensor-pi/env-bot/plotter/pressure.png'
api.update_status_with_media(status=content, filename=image_path)
