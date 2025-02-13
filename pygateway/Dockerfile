FROM alpine:latest

RUN apk add avrdude
RUN apk add python3
RUN apk add py3-pip

RUN python3 -m pip install --upgrade pip --break-system-packages
RUN python3 -m pip install pyserial --break-system-packages
RUN python3 -m pip install paho-mqtt --break-system-packages

WORKDIR /src

COPY . .

ENTRYPOINT ["python3", "./main.py"]