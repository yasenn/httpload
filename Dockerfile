FROM golang:1.14-alpine
MAINTAINER https://github.com/yasenn
WORKDIR /opt/httpload
RUN go build
RUN chmod +x httpload
COPY . .
CMD ["./httpload"]
