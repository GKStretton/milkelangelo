FROM alpine:latest

RUN apk add avrdude
RUN apk add python3
RUN apk add py3-pip

RUN python3 -m pip install --upgrade pip
RUN python3 -m pip install pyserial
RUN python3 -m pip install paho-mqtt

WORKDIR /src

COPY . .

ENTRYPOINT ["python3", "./main.py"]