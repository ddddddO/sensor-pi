from diagrams import Cluster, Diagram

from diagrams.programming.language import Python, Go
from diagrams.onprem.compute import Server
from diagrams.generic.database import SQL
from diagrams.aws.compute import Lambda
from diagrams.aws.integration import SNS, SQS
from diagrams.saas.social import Twitter

with Diagram("Next ENV-Tweet-Bot", direction="TB", show=False, outformat="png"):
  bme280 = Server("bme280\n(sensor)")
  mh_z19 = Server("mh-z19\n(sensor)")
  sensor_c = Server("sensor c")
  
  with Cluster("Raspberry pi"):
    prog_a = Python("Get and Publish\nsensor value")
    prog_b = Python("Get and Publish\nsensor value")
    prog_c = Python("Get and Publish\nsensor value")
    prog1 = Go("Pull and Store\nsensor data")
    db = SQL("sqlite3")
    prog2 = Go("Fetch and Publish\nenv data\n(10 records)")
    prog3 = Python("Pull and Decode\nimage data\nand Post tweet")

  with Cluster("AWS"):
    sqs0 = SQS("sensor values")
    sns = SNS("pressure/\nhumidity/\ntemperature data")

    with Cluster("Generate image"):
      svc_group = [Lambda("pressure"),
                  Lambda("humidity"),
                  Lambda("temperature")]

    sqs1 = SQS("encoded base64\nimage data")

  twitter = Twitter("twitter")

  mh_z19 >> prog_a >> sqs0
  bme280 >> prog_b >> sqs0
  sensor_c >> prog_c >> sqs0
  sqs0 >> prog1
  prog1 >> db >> prog2
  prog2 >> sns >> svc_group
  svc_group >> sqs1 >> prog3
  prog3 >> twitter
