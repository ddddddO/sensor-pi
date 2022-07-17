import base64
import time
import json
import boto3

# ref: https://boto3.amazonaws.com/v1/documentation/api/latest/guide/sqs-example-sending-receiving-msgs.html#receive-and-delete-messages-from-a-queue

# Create SQS client
sqs = boto3.client('sqs')

queue_url = 'https://sqs.ap-northeast-1.amazonaws.com/820544363308/filedata_to_tweeter'

def receive_message() -> (
  str,
  str, float, str,
  str, float, str,
  str, float, str,
  ):
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

  # print('response: %s' % response)

  message = response['Messages'][0] # TODO: ここでKeyErrorがでたら再度取得するようにしてもいいかも
  receipt_handle = message['ReceiptHandle']
  body = message['Body']
  j_body =  json.loads(body)
  environment = j_body['environment']

  pressure = environment[0]
  type_p = pressure['type']
  latest_value_p = pressure['latest']['value']
  encoded_p = pressure['encoded']

  temperature = environment[1]
  type_t = temperature['type']
  latest_value_t = temperature['latest']['value']
  encoded_t = temperature['encoded']

  humidity = environment[2]
  type_h = humidity['type']
  latest_value_h = humidity['latest']['value']
  encoded_h = humidity['encoded']

  return (receipt_handle,
    type_p, latest_value_p, encoded_p,
    type_t, latest_value_t, encoded_t,
    type_h, latest_value_h, encoded_h,
  )


def decode_to_file(name, encoded):  
  decoded = base64.b64decode(encoded)
  file_path = './%s.png' % name

  with open(file_path, "wb") as f:
    f.write(decoded)


def delete_message(type, receipt_handle):
  # Delete received message from queue
  sqs.delete_message(
    QueueUrl=queue_url,
    ReceiptHandle=receipt_handle
  )
  print('received and deleted message: %s' % type)


while True:
  try:
    receipt_handle, type_p, latest_value_p, encoded_p, type_t, latest_value_t, encoded_t, type_h, latest_value_h, encoded_h = receive_message()

    decode_to_file(type_p, encoded_p)
    decode_to_file(type_t, encoded_t)
    decode_to_file(type_h, encoded_h)

    delete_message('all', receipt_handle)

    break
  except Exception as err:
    # TODO: 厳密にエラーハンドリングする
    print('in err')
    print(err)
    time.sleep(3)
    continue
