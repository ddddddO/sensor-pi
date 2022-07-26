from diagrams import Cluster, Diagram, Edge

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

  green_edge = Edge(color="darkgreen")
  brown_edge = Edge(color="brown")
  black_edge = Edge(color="black", style="bold")

  with Cluster("Raspberry Pi 4"):
    prog_a = Python("Get and Publish\nsensor value")
    prog_b = Python("Get and Publish\nsensor value")
    prog_c = Python("Get and Publish\nsensor value")
    prog1 = Go("Pull and Store\nsensor data")
    db = SQL("SQLite3")
    prog2 = Go("Fetch and Publish\nenv data\n(10 records)")
    prog3 = Python("Pull and Decode\nimage data\nand Post tweet")

  sensor_d = Server("sensor d")
  sensor_e = Server("sensor e")

  with Cluster("Host x"):
    prog_d = Python("Get and Publish\nsensor value")
    prog_e = Python("Get and Publish\nsensor value")

  with Cluster("AWS"):
    sqs0 = SQS("sensor values")
    sns = SNS("sensor values")

    with Cluster("Generate image"):
      svc_group = [Lambda("pressure"),
                  Lambda("humidity"),
                  Lambda("temperature"),
                  Lambda("co2"),
                  Lambda("sensor c"),
                  Lambda("sensor d"),
                  Lambda("sensor e")]

    sqs1 = SQS("encoded base64\nimage data")

  twitter = Twitter("Twitter")

  mh_z19 >> black_edge >> prog_a >> green_edge >> sqs0
  bme280 >> black_edge >> prog_b >> green_edge >> sqs0
  sensor_c >> black_edge >> prog_c >> green_edge >> sqs0

  sensor_d >> black_edge >> prog_d >> green_edge >> sqs0
  sensor_e >> black_edge >> prog_e >> green_edge >> sqs0

  sqs0 >> brown_edge >> prog1
  prog1 >> black_edge >> db >> black_edge >> prog2
  prog2 >> black_edge >> sns >> green_edge >> svc_group
  svc_group >> green_edge >> sqs1 >> brown_edge >> prog3
  prog3 >> black_edge >> twitter
