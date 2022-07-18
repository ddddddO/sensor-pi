#!/usr/bin/env python3.8

import base64
import time
import json
import os
import settings
import tweepy
import boto3

# ref: https://boto3.amazonaws.com/v1/documentation/api/latest/guide/sqs-example-sending-receiving-msgs.html#receive-and-delete-messages-from-a-queue

# Create SQS client
sqs = boto3.client('sqs')

queue_url = 'https://sqs.ap-northeast-1.amazonaws.com/820544363308/filedata_to_tweeter'

def receive_message() -> (str, str, float, str):
  # Receive message from SQS queue
  # TODO: 何度か実行しないと取得できないことがままある。なので何度か実行する作りにする必要がある
  response = sqs.receive_message(
    QueueUrl=queue_url,
    AttributeNames=[
      'SentTimestamp'
    ],
    MaxNumberOfMessages=1,
    MessageAttributeNames=[
      'All'
    ],
    VisibilityTimeout=0,
    WaitTimeSeconds=3
  )

  message = response['Messages'][0] # TODO: ここでKeyErrorがでたら再度取得するようにしてもいいかも
  receipt_handle = message['ReceiptHandle']
  body = message['Body']
  j_body =  json.loads(body)
  environment = j_body['environment']

  env = environment[0]
  type = env['type']
  latest_value = env['latest']['value']
  encoded = env['encoded']

  return receipt_handle, type, latest_value, encoded


def decode_to_file(name, encoded):  
  decoded = base64.b64decode(encoded)
  file_path = '/tmp/%s.png' % name

  with open(file_path, "wb") as f:
    f.write(decoded)


def delete_message(type, receipt_handle):
  # Delete received message from queue
  sqs.delete_message(
    QueueUrl=queue_url,
    ReceiptHandle=receipt_handle
  )
  print('received and deleted message: %s' % type)


if __name__ == '__main__':

	received_temperature = False
	received_pressure = False
	received_humidity = False

	while True:
		if received_temperature and received_pressure and received_humidity:
			break

		try:
			receipt_handle, type, latest_value, encoded = receive_message()

			if type == 'temperature':
				latest_value_t = latest_value
				received_temperature = True

			if type == 'pressure':
				latest_value_p = latest_value
				received_pressure = True

			if type == 'humidity':
				latest_value_h = latest_value
				received_humidity = True

			decode_to_file(type, encoded)
			delete_message(type, receipt_handle)

			time.sleep(1)
			continue
		except Exception as err:
			# TODO: 厳密にエラーハンドリングする
			print('in err')
			print(err)
			time.sleep(3)
			continue

	try:
		title = 'ただいまの気温・気圧・湿度'
		location = '@多摩川付近(屋内)'
		temp = "temperature : {:-6.2f} ℃".format(latest_value_t)
		pressure = "pressure : %7.2f hPa" % (latest_value_p)
		hum = "humidity : %6.2f ％" % (latest_value_h)
		content = title + location + '\n' + temp + '\n' + pressure + '\n' + hum

		auth = tweepy.OAuthHandler(settings.consumer_key, settings.consumer_secret)
		auth.set_access_token(settings.token, settings.token_secret)
		api = tweepy.API(auth)

		image_path = '/tmp/pressure.png'
		status = api.update_status_with_media(status=content, filename=image_path)

		image_path = '/tmp/temperature.png'
		status = api.update_status_with_media(in_reply_to_status_id=status.id, status='', filename=image_path)

		image_path = '/tmp/humidity.png'
		status = api.update_status_with_media(in_reply_to_status_id=status.id, status='', filename=image_path)

	except Exception as err:
		# TODO: error handling
		print('Exception!: {err}'.format(err=err))
	finally:
		pass
