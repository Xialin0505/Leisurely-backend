FROM golang:latest

RUN apt-get update && apt-get upgrade -y &&\
  apt-get install -y git &&\
  apt-get -y install python3-pip &&\
  pip3 install numpy &&\
  pip3 install spacy &&\
  pip3 install pandas &&\
  pip3 install Wikipedia-API &&\
  pip3 install git+https://github.com/boudinfl/pke.git &&\
  pip3 install nltk &&\
  python3 -m spacy download en_core_web_sm
  # make openssh-client

RUN apt-get install -y vim

WORKDIR /app

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh &&\
  chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

CMD air

#FROM postgres:latest
#RUN psql -U postgres -c 'create database leisurely;'