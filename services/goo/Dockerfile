FROM golang:1.19

RUN apt update
RUN apt install -y ffmpeg
RUN apt install -y gphoto2
RUN apt install -y python3 python3-pip

RUN python3 -m pip install opencv-python scikit-image matplotlib pyyaml pandas

WORKDIR /src

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go mod download

COPY . .

RUN go build -race

ENTRYPOINT ["./goo"]