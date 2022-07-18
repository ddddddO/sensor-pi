from diagrams import Cluster, Diagram

from diagrams.programming.language import Python, Go
from diagrams.onprem.compute import Server
from diagrams.generic.database import SQL
from diagrams.aws.compute import Lambda
from diagrams.aws.integration import SNS, SQS
from diagrams.saas.social import Twitter

with Diagram("ENV-Tweet-Bot", show=False, outformat="png"):
  bme280 = Server("bme280\n(sensor)")
  
  with Cluster("Raspberry pi"):
    prog1 = Python("Get and Store\nsensor value")
    db = SQL("sqlite3")
    prog2 = Go("Fetch and Publish\nenv data\n(10 records)")
    prog3 = Python("Pull and Decode\nimage data\nand Post tweet")

  with Cluster("AWS"):
    sns = SNS("pressure/\nhumidity/\ntemperature data")

    with Cluster("Generate image"):
      svc_group = [Lambda("pressure"),
                  Lambda("humidity"),
                  Lambda("temperature")]

    sqs = SQS("encoded base64\nimage data")

  twitter = Twitter("twitter")

  bme280 >> prog1 >> db >> prog2
  prog2 >> sns >> svc_group
  svc_group >> sqs >> prog3
  prog3 >> twitter
