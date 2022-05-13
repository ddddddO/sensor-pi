#!/usr/bin/env python3

import settings
import twitter

title = 'ただいまの気温・気圧・湿度(屋内)'
location = '@多摩川付近'
temp = input()
pressure = input()
hum = input()
content = title + location + '\n' + temp + '\n' + pressure + '\n' + hum

auth = twitter.OAuth(
    consumer_key=settings.consumer_key,
    consumer_secret=settings.consumer_secret,
    token=settings.token,
    token_secret=settings.token_secret)
t = twitter.Twitter(auth=auth)
t.statuses.update(status=content)