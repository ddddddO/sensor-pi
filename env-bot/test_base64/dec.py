import base64

data = input()

encoded = base64.b64decode(data)
file_path = "./decoded.png"

with open(file_path, "wb") as f:
  f.write(encoded)