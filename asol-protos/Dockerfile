FROM golang:1.19
WORKDIR /src

RUN apt-get update && apt-get install -y unzip gcc-multilib

# Install Python and pip
RUN apt-get install -y python3 python3-pip python3-venv

# Install python-betterproto
# (maybe should use venv instead of break pkg, but it's a container so not too important)
RUN pip3 install --break-system-packages "betterproto[compiler]"
RUN pip3 install --break-system-packages betterproto


# Install nanopb
ARG NANOPB_VERSION=0.4.6
RUN wget https://jpa.kapsi.fi/nanopb/download/nanopb-${NANOPB_VERSION}-linux-x86.tar.gz && \
    tar -xvf nanopb-${NANOPB_VERSION}-linux-x86.tar.gz -C /opt

# Install protoc generator for Go.
ARG PROTOC_GEN_GO_VERSION=1.28.1
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VERSION}

# Install protoc
ARG PROTOC_VERSION=21.9
RUN PROTOC_ZIP=protoc-${PROTOC_VERSION}-linux-x86_64.zip && \
    curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/$PROTOC_ZIP && \
    unzip -o $PROTOC_ZIP -d /usr/local bin/protoc && \
    unzip -o $PROTOC_ZIP -d /usr/local 'include/*'  && \
    rm -f $PROTOC_ZIP

# Install Node.js and npm
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
RUN apt-get install -y nodejs

# Install protoc-gen-ts
# Hard to do json conversions with this library
# RUN npm install -g ts-protoc-gen protoc-gen-js
RUN npm install -g ts-proto

# Copy src code.
COPY . .
