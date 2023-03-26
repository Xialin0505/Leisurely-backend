import pandas as pd
import wikipediaapi
import numpy as np
import re
import pke
import json

import pprint
import datetime
import threading
import random
#import yake

# from spacy.tokens import DocBin
import spacy
import requests

import nltk; nltk.download('wordnet'); nltk.download('omw-1.4')
from nltk.corpus import wordnet as wn

import socket
#import nltk; nltk.download('stopwords')
#from nltk.corpus import stopwords
#stop_words = stopwords.words('english')

#from thinc.api import Config
#from spacy.lang.en import English

#def train():
#    config = Config().from_disk("./train.cfg")
#    nlp = English.from_config(config)

class ExtractKeyword (threading.Thread):
    def __init__(self, keyword):
        threading.Thread.__init__(self)
        self.keyword = keyword
        self._return = None
    
    def run(self):
        wiki_wiki = wikipediaapi.Wikipedia(
            language='en',
            extract_format=wikipediaapi.ExtractFormat.WIKI
        )
        self._return = wiki_wiki.page(self.keyword).summary

    def join(self):
        threading.Thread.join(self)
        return self._return
      #print "Starting " + self.name
      #print_time(self.name, 5, self.counter)
      #print "Exiting " + self.name

def extractPerson(title):
        nlp = spacy.load("en_core_web_sm")

        text = title
        #text = text.replace("\'t", "")
        text = re.sub('[^A-Za-z0-9]+', ' ', text)

        # preprocess event title and description
        # @, -, :

        doc = nlp(text.strip())

        named_entities = set()

        for i in doc.ents:
                entry = str(i.lemma_).lower()
                text = text.replace(str(i).lower(), "")
                if i.label_ in ["PERSON"]:
                    named_entities.add(entry.title())

        return named_entities

def getSections(sections, text):
    for s in sections:
        text += getSections(s.sections, s.title + ": " + s.text)
        # uncomment it if only want first section
        #break
    return text

def parseEvent(title):
    #nlp = spacy.load("en_core_web_sm")
    #title = title.replace("\'s", "")
    #title = title.replace("\'t", "")
    text = re.split(r"[^a-zA-Z0-9\s\']", title)
    text = list(filter(('').__ne__, text))

    return text

def search(title):
    #print("starting search")
    parsed = parseEvent(title)
    parsed.extend(extractPerson(parsed[0]))

    #print(parsed)
    #print(person)

    wiki_wiki = wikipediaapi.Wikipedia(
            language='en',
            extract_format=wikipediaapi.ExtractFormat.WIKI
    )

    summaries = ""

    threads = []

    for keyword in parsed:
        thread = ExtractKeyword(keyword)
        thread.start()
        threads.append(thread)

    for t in threads:
        result = t.join()
        if result != None:
            summaries += result
    
    # try:

    #     for keyword in parsed:
    #         #print(keyword)
    #         p_wiki = wiki_wiki.page(keyword)
    #         summaries += p_wiki.summary

        # for entities in person:
        #     #print(entities)
        #     wikiText = ""
        #     p_wiki = wiki_wiki.page(entities)
        #     #wikiText = getSections(p_wiki.sections, wikiText)
        #     #print(wikiText)
        #     #print(p_wiki.summary)
        #     #return p_wiki.summary
        #     summaries += p_wiki.summary

    # except:
    #     print("Cannot find content in WIKI")

    #print("end search")
    return summaries

def extractKeyword(text):
    # initialize keyphrase extraction model, here TopicRank
    extractor = pke.unsupervised.TopicRank()

    # load the content of the document, here document is expected to be a simple 
    # test string and preprocessing is carried out using spacy
    extractor.load_document(input=text, language='en')

    # keyphrase candidate selection, in the case of TopicRank: sequences of nouns
    # and adjectives (i.e. `(Noun|Adj)*`)
    extractor.candidate_selection()

    # candidate weighting, in the case of TopicRank: using a random walk algorithm
    extractor.candidate_weighting()

    # N-best selection, keyphrases contains the 10 highest scored candidates as
    # (keyphrase, score) tuples
    keyphrases = extractor.get_n_best(n=5)
    return keyphrases

def tagEvent(event):

    # data = """{
    #     "events" : []
    # }"""

    # events = []

    # allEvents = json.loads(data)

    ## search from wiki
    text = search(event["Title"]) 
    text += event["Description"]

    ## extract keyword
    keyword = extractKeyword(text)

    # eachEvent = """{
    #     "eventTitle" : "title",
    #     "tags" : [],
    #     "score" : 0.0
    # }"""

    tags = []

    for key in keyword:
        # person = extractPerson(key[0])
        person = []
        if len(person) == 0:
            tags.extend(key[0].split())
        else:
            tags.append(key[0])

    tags = list(set(tags))
    
    #currentE = json.loads(eachEvent)
    event["EventTags"] = tags
    #currentE["eventTitle"] = t


    #events.append(currentE)
    #print(json.dumps(currentE, indent=4))
    
    #allEvents["events"] = events
    return event

def reduceTag(taglist, preference):   
    # print("ReduceTag ", len(taglist))
    score = 0

    for tag in taglist: 
        try:
            tagsynset = wn.synsets(tag)[0]
        except:
            continue

        for p in preference:
            # print(p[0])
            # print(tagsynset)
            sim = tagsynset.path_similarity(p[0])
            
            #print(sim)
            if sim != None:
                score += sim * p[1]

    # for synset in wn.synsets(tag):
    #     for lemma in synset.lemma_names():
    #         for j in range(len(preference)):
    #             if lemma == preference[j][0]:
    #                 print(lemma)
    #                 return preference[j][0]
    
    return score

def scoreEvent(event, preference):
    tags = event["EventTags"]
    score = 0
    #newTags = []
    total_count = 0

    psyn = []

    try:
        for p in preference:
            if p["EventType"] != 1:
                tagWord = p["Description"].split(" ")
                for w in tagWord:
                    try:
                        psyn.append((wn.synsets(w)[0], p["Count"]))
                    except:
                        continue
            else:
                for tag in tags:
                    if tag == p["Description"]:
                        score += p["Count"]
                
            total_count += p["Count"]


        score += reduceTag(tags,psyn)
        if (total_count != 0):
            score = score / total_count
        else:
            score = 1
        #newTags.append(result)

        # for tag in newTags:
        #     for p in preference:
        #         total_count += p[1]
        #         if tag == p[0]:
        #             score += p[1]

        event["Score"] = score
        #event["tags"] = tags
    except:
        return None

    return event

def scoreEventD(event, preference):
    events = []
    for e in event:
        tags = e["EventTags"]
        score = 0
        #newTags = []
        total_count = 0

        psyn = []

        for p in preference:
            if p["EventType"] != 1:
                tagWord = p["Description"].split(" ")
                for w in tagWord:
                    try:
                        psyn.append((wn.synsets(w)[0], p["Count"]))
                    except:
                        continue
            else:
                for tag in tags:
                    if tag == p["Description"]:
                        score += p["Count"]
                
            total_count += p["Count"]


        score += reduceTag(tags,psyn)
        if (total_count != 0):
            score = score / total_count
        else:
            score = 1

        if (score != 0 and len(tags) != 0):
            e["Score"] = score
            #event["tags"] = tags
            events.append(e)
        #newTags.append(result)

        # for tag in newTags:
        #     for p in preference:
        #         total_count += p[1]
        #         if tag == p[0]:
        #             score += p[1]

    # if len(events) <= 10:
    #     return events
    
    # return random.sample(events, 10)
    return events

def tagEventD(event):

    # data = """{
    #     "events" : []
    # }"""

    # events = []

    # allEvents = json.loads(data)

    ## search from wiki
    text = search(event[0]) 

    ## extract keyword
    keyword = extractKeyword(text)

    eachEvent = """{
        "eventTitle" : "title",
        "EventTags" : [],
        "score" : 0.0
    }"""

    tags = []

    for key in keyword:
        #person = extractPerson(key[0])
        person = []
        if len(person) == 0:
            tags.extend(key[0].split())
        else:
            tags.append(key[0])

    tags = list(set(tags))
    
    currentE = json.loads(eachEvent)
    currentE["EventTags"] = tags
    currentE["eventTitle"] = event[0]


    #events.append(currentE)
    #print(json.dumps(currentE, indent=4))
    
    #allEvents["events"] = events
    return currentE

# def server():
#     host = '127.0.0.1'
#     port = 4000

#     s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
#     s.bind((host, port))
#     s.listen(5)

#     print("start listening")

#     while True:
#         conn, addr = s.accept()
#         events = recvData(conn)
#         print(json.dumps(events, indent=4))
#         events = scoreEvents(events)
    

def main():
    title = ["""Joe Hisaishi with the TSO in Concert""", """TD Toronto Jazz Festival"""]
    preference = [["music", 1, 2], ["restaurant",1, 2], ["joe biden", 1, 1], ["cat", 1, 3]]
    #preference = [["festival", 1, 2], ["music",1, 2], ["joe hisaishi", 1, 1], ["outdoor", 1, 3]]
    #result = tagEventD(title)
    #result = scoreEventD(result, preference)
    #print(json.dumps(result, indent=4))


if __name__ == "__main__":
   main()