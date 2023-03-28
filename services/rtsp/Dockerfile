FROM aler9/rtsp-simple-server:v0.21.6 AS rtsp
FROM restreamio/gstreamer:latest-prod
COPY --from=rtsp /rtsp-simple-server /
COPY --from=rtsp /rtsp-simple-server.yml /

RUN apt update
RUN apt install v4l-utils -y
RUN apt install python3 python3-pip -y
RUN python3 -m pip install --upgrade pip
RUN python3 -m pip install pyyaml

ENTRYPOINT [ "/rtsp-simple-server" ]