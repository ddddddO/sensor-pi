# ref: https://monomonotech.jp/kurage/raspberrypi/daiso_btshutter.html
import evdev
# import subprocess
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
                        # subprocess.Popen(['mpg321', 'decision3.mp3'])
                        print('Up!')

                    if event.code == evdev.ecodes.KEY_ENTER:
                        # subprocess.Popen(['amixer', 'sset', 'PCM', '10%-', '-M'])
                        print('Down!')
    except:
        print('Retry...')
        time.sleep(1)