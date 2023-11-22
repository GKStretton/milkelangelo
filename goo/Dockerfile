FROM golang:1.19

RUN apt update && apt install -y \
    ffmpeg \
    gphoto2 \
    python3 \
    python3-pip

RUN python3 -m pip install opencv-python scikit-image matplotlib pyyaml pandas

WORKDIR /src

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./asol-protos ../asol-protos

RUN go mod download

COPY . .

RUN go build -race

ENTRYPOINT ["./goo"]