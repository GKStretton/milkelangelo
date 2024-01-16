FROM golang:1.21.5

RUN apt update
RUN apt install -y python3 python3-setuptools python3-pip python3-venv
RUN apt install -y libjpeg-dev zlib1g-dev
RUN apt install -y imagemagick
RUN apt install -y libgl1
RUN sed -i '/<policy domain="path" rights="none" pattern="\@\*"/d' /etc/ImageMagick-6/policy.xml


RUN python3 --version

WORKDIR /src

# setup python env
COPY ./tools/requirements.txt ../tools/requirements.txt
# maybe a venv would be better but this is faster
RUN cd ../tools && python3 -m pip install --break-system-packages -r requirements.txt
COPY ./tools ../tools

COPY ./runner/go.mod ./go.mod
COPY ./runner/go.sum ./go.sum
COPY ./goo ../goo
COPY ./asol-protos ../asol-protos
COPY ./scripts ../scripts
COPY ./resources ../resources

RUN go mod download

COPY ./runner .

RUN go build

ENTRYPOINT ["./runner"]