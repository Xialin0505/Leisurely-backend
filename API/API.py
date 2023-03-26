import requests
import json
from threading import Thread
import time
import numpy as np
import matplotlib.pyplot as plt


class myThread (Thread):
    def __init__(self):
        Thread.__init__(self)
        self.duration = 0
    def run(self):
        start = time.time()
        generatePlan()
        end = time.time()
        self.duration = end - start
        

def generatePlan():
    URL = "http://Leisurely-lb-80-1889731165.ca-central-1.elb.amazonaws.com:80/v1/generatePlanFree/2"

    data =  {
        "startTime": 14.00,
        "endTime": 23.00,
        "date": "Mar 5",
        "location": "Kitchener",
        "country": "Canada",
        "transport": "Transit"
    }

    response = requests.post(URL, data)
    #print(json.dumps(response.json(), indent=4))
    #print(response)

def runTest(numThread):
    threads = []
    for i in range(numThread):
        threads.append(myThread())

    for i in range(numThread):
        threads[i].start()

    for i in range(numThread):
        threads[i].join()

    timeSpend = 0
    for i in range(numThread):
        timeSpend += threads[i].duration

    return round(timeSpend/numThread, 2)

if __name__ == "__main__":
    numReq = [1, 2, 5, 10]
    for i in numReq:
        sum = 0
        for j in range(5):
            value = runTest(i)
            sum += value

        print("average response time:", sum/5, "s", "for", i, "requests")


    