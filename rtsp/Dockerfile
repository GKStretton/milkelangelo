FROM aler9/rtsp-simple-server AS rtsp
FROM restreamio/gstreamer:latest-prod
COPY --from=rtsp /rtsp-simple-server /
COPY --from=rtsp /rtsp-simple-server.yml /
ENTRYPOINT [ "/rtsp-simple-server" ]