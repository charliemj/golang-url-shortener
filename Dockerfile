# golang image where workspace (GOPATH) configured at /go.
FROM golang:latest

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/charliemj/golang-url-shortener

# Build the golang-docker command inside the container.
RUN go install github.com/charliemj/golang-url-shortener

# Run the golang-docker command when the container starts.
ENTRYPOINT /go/bin/golang-url-shortener

FROM nginx:alpine
COPY default.conf /etc/nginx/conf.d/default.conf
COPY index.html /usr/share/nginx/html/index.html
COPY show.html /usr/share/nginx/html/show.html


# http server listens on port 8080.
EXPOSE 8080

## iron/go:dev is the alpine image with the go tools added
#FROM ubuntu
##FROM iron/go:dev
##WORKDIR /golang-url-shortener
### Set an env var that matches your github repo name, replace treeder/dockergo here with your repo name
##ENV SRC_DIR=/go/src/github.com/charliemj/golang-url-shortener/
### Add the source code:
##ADD . $SRC_DIR
### Build it:
##RUN go install github.com/charliemj/golang-url-shortener/
##ENTRYPOINT /go/bin/golang-url-shortener#
#

## Install Nginx#

## Add application repository URL to the default sources
## RUN echo "deb http://archive.ubuntu.com/ubuntu/ raring main universe" >> /etc/apt/sources.list#

## Update the repository
#RUN apt-get update#

## Install necessary tools
#RUN apt-get install -y vim wget dialog net-tools#

#RUN apt-get install -y nginx#

## Remove the default Nginx configuration file
#RUN rm -v /etc/nginx/nginx.conf#

## Copy a configuration file from the current directory
#ADD nginx.conf /etc/nginx/#

#RUN mkdir /etc/nginx/logs#

## Add a sample index file
#ADD index.html /www/data/
#COPY index.html /www/data/
#COPY show.html /www/data/#

## Append "daemon off;" to the beginning of the configuration
#RUN echo "daemon off;" >> /etc/nginx/nginx.conf#

#EXPOSE 8080#

#CMD ["nginx"]