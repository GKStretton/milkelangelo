FROM python:3.9

RUN python -m pip install --upgrade pip
RUN python -m pip install paho-mqtt
RUN python -m pip install pyyaml

WORKDIR /src

COPY . .

ENTRYPOINT ["python3", "./main.py"]