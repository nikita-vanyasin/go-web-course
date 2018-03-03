FROM scratch

ENV CONTENT_FOLDER_PATH /content
ENV CONNECTION_STRING root@/simple_video_server

VOLUME $CONTENT_FOLDER_PATH

COPY go-daemon /

CMD ["/go-daemon"]