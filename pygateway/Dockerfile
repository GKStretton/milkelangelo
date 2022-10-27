FROM python:3.9

RUN python -m pip install --upgrade pip
RUN python -m pip install pyserial
RUN python -m pip install paho-mqtt

WORKDIR /src

COPY . .

ENTRYPOINT ["python3", "./main.py"]