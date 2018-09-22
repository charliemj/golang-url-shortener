# iron/go:dev is the alpine image with the go tools added
FROM iron/go:dev
WORKDIR /golang-url-shortener
# Set an env var that matches your github repo name, replace treeder/dockergo here with your repo name
ENV SRC_DIR=/go/src/github.com/charliemj/golang-url-shortener/
# Add the source code:
ADD . $SRC_DIR
# Build it:
RUN go install github.com/charliemj/golang-url-shortener/
ENTRYPOINT /go/bin/golang-url-shortener
EXPOSE 8080

FROM nginx
COPY public /usr/share/nginx/html

