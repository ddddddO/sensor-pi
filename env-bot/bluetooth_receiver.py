# ref: https://monomonotech.jp/kurage/raspberrypi/daiso_btshutter.html
from asyncio import subprocess
import evdev
import subprocess
import time

while True:
    try:
        # ls /dev/input でevent番号要確認
        device = evdev.InputDevice('/dev/input/event0')
        print(device)

        for event in device.read_loop():
            if event.type == evdev.ecodes.EV_KEY:
                if event.value == 1: # 0:KEYUP, 1:KEYDOWN
                    print(event.code)

                    if event.code == evdev.ecodes.KEY_VOLUMEUP:
                        print('Received!')
                        subprocess.Popen(['/home/pi/go/bin/dagu', 'start', '/home/pi/github.com/ddddddO/sensor-pi/env-bot/dag.yaml'])
    except:
        print('Retry...')
        time.sleep(1)
