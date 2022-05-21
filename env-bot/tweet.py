#!/usr/bin/env python3

import os
import settings
import tweepy
from bme280_repository import Repository

if __name__ == '__main__':
	try:
		dsn = os.getenv('DSN')
		repo = Repository(dsn)
		t, p, h = repo.fetch()

		title = 'ただいまの気温・気圧・湿度(屋内)'
		location = '@多摩川付近'
		temp = "temp : {:-6.2f} ℃".format(t)
		pressure = "pressure : %7.2f hPa" % (p)
		hum = "hum : %6.2f ％" % (h)
		content = title + location + '\n' + temp + '\n' + pressure + '\n' + hum

		auth = tweepy.OAuthHandler(settings.consumer_key, settings.consumer_secret)
		auth.set_access_token(settings.token, settings.token_secret)
		api = tweepy.API(auth)

		image_path = os.getenv('PRESSURE_IMAGE_PATH')
		api.update_status_with_media(status=content, filename=image_path)
	except Exception as err:
		# TODO: error handling
		print('Exception!: {err}'.format(err=err))
	finally:
		repo.close()
