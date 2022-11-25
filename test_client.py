import socket
import time
import datetime

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM) # UDP
while True:
    sendme = "Hello, world the time is " + str(datetime.datetime.now())
    print(sendme)
    sock.sendto(bytes(sendme, "utf-8"), ("127.0.0.1", 10000))
    time.sleep(3)