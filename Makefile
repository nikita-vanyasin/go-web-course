

env CGO_ENABLED=0 GOOS=linux go build -o go-daemon
env CGO_ENABLED=0 GOOS=linux go build -o go-server


# docker build -t=go-web-course-server .
# docker run -v ${CONTENT_FOLDER_PATH}:/content  --net=host -d --name=go-server-1 go-web-course-server


# docker build -t=go-web-course-daemon .
# docker run -v ${CONTENT_FOLDER_PATH}:/content  --net=host -d --name=go-daemon-1 go-web-course-daemon
