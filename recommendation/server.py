from textwrap import indent
from recommendation import *
import socket
import threading
import time

class processEvent (threading.Thread):
    def __init__(self, threadID, event, pref):
        threading.Thread.__init__(self)
        self.threadID = threadID
        self.event = event
        self.pref = pref
        self._return = None
   
    def run(self):
        #print ("Starting " + str(self.threadID))
        #start_time = datetime.datetime.now()
        self._return = tagEvent(self.event)
        self._return = scoreEvent(self._return, self.pref)
        # end_time = datetime.datetime.now()
        # print(end_time - start_time)

    def join(self):
        threading.Thread.join(self)
        return self._return

class ThreadedServer(object):
    def __init__(self, host, port):
        self.host = host
        self.port = port
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        self.sock.bind((self.host, self.port))

    def listen(self):
        self.sock.listen(5)
        print("start listening")

        while True:
            client, address = self.sock.accept()
            client.settimeout(60)
            threading.Thread(target = self.listenToClient,args = (client,address)).start()

    def listenToClient(self, client, address):
        while True:
            events, prefs = recvData(client)
            if (events != None):
                #print(json.dumps(events, indent=4))
                events = scoreEvents(events, prefs)
                sendEvent(client, events)
                break
            time.sleep(10)
        
        client.close()


def scoreEvents(events, prefs):

    data = """{
            "Events": []
        }"""

    allEvent = json.loads(data)

    result = []
    threadList = []

    i = 0 
    for event in events:
        thread = processEvent(i, event, prefs)
        threadList.append(thread)
        i = i+1

    for thread in threadList:
        thread.start()

    for thread in threadList:
        e = thread.join()
        if e != None and len(e["EventTags"]) != 0:
            result.append(e)

    #result = scoreEventD(result, prefs)
    allEvent["Events"] = result

    return allEvent

def recvData(conn):
    preflen = int.from_bytes(conn.recv(4), byteorder='big')

    preference = conn.recv(preflen).decode("utf-8")
    preferences = json.loads(preference)

    msglen = int.from_bytes(conn.recv(4), byteorder='big')

    fragments = []
    while msglen > 0: 
        chunk = conn.recv(10000)
        msglen -= len(chunk)
        
        fragments.append(chunk)

    if len(fragments) == 0:
        return None

    data = b''.join(fragments).decode("utf-8") 
        
    events = json.loads(data)
    # print(json.dumps(events, indent=4))
    return [events, preferences]

def sendEvent(conn, events):
    print(json.dumps(events, indent=4))

    dataInByte = bytes(json.dumps(events),encoding="utf-8")
    datalen = len(dataInByte)
    #print(datalen)
    # conn.sendall(datalen.to_bytes(8, byteorder='big'))
    conn.sendall(dataInByte)

def server():
    host = '127.0.0.1'
    port = 4000

    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.bind((host, port))
    s.listen(5)

    wn.ensure_loaded() 
    print("start listening")

    while True:
        ## make this accept multi-thread
        conn, addr = s.accept()
        events = recvData(conn)
        #print(json.dumps(events, indent=4))
        events = scoreEvents(events)
        sendEvent(conn, events)

def main():
    f = open('../setting.json')
    data = json.load(f)

    hostName = data['BackendRecom_Name']
    print(hostName)

    host = hostName
    port = 4000
    #server()
    ThreadedServer(host,port).listen()

if __name__ == "__main__":
    main()