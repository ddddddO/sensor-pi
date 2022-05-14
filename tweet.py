#!/usr/bin/env python3

import settings
import tweepy

auth = tweepy.OAuthHandler(settings.consumer_key, settings.consumer_secret)
auth.set_access_token(settings.token, settings.token_secret)

api = tweepy.API(auth)

title = 'ただいまの気温・気圧・湿度(屋内)'
location = '@多摩川付近'
temp = input()
pressure = input()
hum = input()
content = title + location + '\n' + temp + '\n' + pressure + '\n' + hum

image_path = '/home/pi/github.com/ddddddO/sensor-pi/plotter/pressure.png'

api.update_status_with_media(status=content, filename=image_path)
