FROM golang:1.19

RUN apt update
RUN apt install -y ffmpeg
RUN apt install -y gphoto2

WORKDIR /src

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

RUN go mod download

COPY . .

RUN go build

ENTRYPOINT ["./goo"]