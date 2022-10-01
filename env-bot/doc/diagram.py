from diagrams import Cluster, Diagram, Edge

from diagrams.programming.language import Python, Go
from diagrams.onprem.compute import Server
from diagrams.generic.database import SQL
from diagrams.aws.compute import Lambda
from diagrams.aws.integration import SNS, SQS
from diagrams.saas.social import Twitter
from diagrams.aws.iot import IotButton
from diagrams.gcp.devtools import Scheduler

with Diagram("ENV-Tweet-Bot", show=False, outformat="png"):
  bme280 = Server("bme280\n(sensor)")
  mh_z19 = Server("mh_z19\n(sensor)")
  bt_remote_controller = IotButton("Bluetooth remocon\n(not AWS IotButton)")

  green_edge = Edge(color="darkgreen")
  brown_edge = Edge(color="brown")
  black_edge = Edge(color="black", style="bold")

  with Cluster("Raspberry Pi 4"):
    cron = Scheduler("Execute Dagu\nat 9 and 18\n(not GCP Scheduler)")
    bt_receiver = Python("Receive bluetooth \nand Execute Dagu")

    with Cluster("Managed by Dagu"):
      prog1 = Python("Get and Store\nsensor value")
      prog4 = Python("Get and Store\nsensor value")
      db = SQL("SQLite3")
      prog2 = Go("Fetch and Publish\nenv data\n(10 records)")
      prog3 = Python("Pull and Decode\nimage data\nand Post tweet")

  with Cluster("AWS"):
    sns = SNS("pressure/\nhumidity/\ntemperature and co2 data")

    with Cluster("Generate image"):
      svc_group = [Lambda("pressure"),
                  Lambda("humidity"),
                  Lambda("temperature"),
                  Lambda("co2")]

    sqs = SQS("encoded base64\nimage data")

  twitter = Twitter("Twitter")

  cron
  bt_remote_controller >> bt_receiver
  bme280 >> black_edge >> prog1 >> black_edge >> db
  mh_z19 >> black_edge >> prog4 >> black_edge >> db
  db >> black_edge >> prog2 >> black_edge >> sns
  sns >> green_edge >> svc_group
  svc_group >> green_edge >> sqs >> brown_edge >> prog3
  prog3 >> black_edge >> twitter
