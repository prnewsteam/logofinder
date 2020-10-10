FROM golang:1.15

# install imagick
RUN apt-get update
RUN apt-get install imagemagick libmagickwand-dev librsvg2-bin -y
RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
RUN dpkg -i google-chrome-stable_current_amd64.deb; apt-get -fy install

ADD . /go/github.com/prnewsteam/logofinder
WORKDIR /go/github.com/prnewsteam/logofinder

RUN go mod download
RUN go install

EXPOSE 8099

ENTRYPOINT ["/go/bin/logofinder"]