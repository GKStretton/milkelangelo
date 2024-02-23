FROM golang:1.21

RUN apt update
RUN apt install -y gphoto2
RUN apt install -y python3 python3-pip

# todo: change to venv
RUN python3 -m pip install --break-system-packages opencv-python-headless scikit-image matplotlib pyyaml
RUN python3 -m pip install --break-system-packages --pre betterproto==2.0.0b5

WORKDIR /src

COPY ./dslrcapture/go.mod ./go.mod
COPY ./dslrcapture/go.sum ./go.sum
COPY ./goo ../goo
COPY ./asol-protos ../asol-protos
COPY ./tools ../tools
COPY ./resources ../resources

RUN go mod download

COPY ./dslrcapture .

RUN go build -race

ENTRYPOINT ["./dslrcapture"]