from diagrams import Cluster, Diagram, Edge

from diagrams.programming.language import Python, Go
from diagrams.onprem.compute import Server
from diagrams.generic.database import SQL
from diagrams.aws.compute import Lambda
from diagrams.aws.integration import SNS, SQS
from diagrams.saas.social import Twitter

with Diagram("ENV-Tweet-Bot", show=False, outformat="png"):
  bme280 = Server("bme280\n(sensor)")

  green_edge = Edge(color="darkgreen")
  brown_edge = Edge(color="brown")
  black_edge = Edge(color="black", style="bold")

  with Cluster("Raspberry Pi 4"):
    prog1 = Python("Get and Store\nsensor value")
    db = SQL("SQLite3")
    prog2 = Go("Fetch and Publish\nenv data\n(10 records)")
    prog3 = Python("Pull and Decode\nimage data\nand Post tweet")

  with Cluster("AWS"):
    sns = SNS("pressure/\nhumidity/\ntemperature data")

    with Cluster("Generate image"):
      svc_group = [Lambda("pressure"),
                  Lambda("humidity"),
                  Lambda("temperature")]

    sqs = SQS("encoded base64\nimage data")

  twitter = Twitter("Twitter")

  bme280 >> black_edge >> prog1 >> black_edge >> black_edge >> db >> black_edge >> prog2 >> black_edge >> sns
  sns >> green_edge >> svc_group
  svc_group >> green_edge >> sqs >> brown_edge >> prog3
  prog3 >> black_edge >> twitter
