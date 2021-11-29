# https://dev.to/plutov/docker-and-go-modules-3kkn
FROM golang:latest
ENV http_proxy 10.158.100.6:8080
ENV https_proxy 10.158.100.6:8080

RUN apt-get install curl
COPY kubectl ./
RUN mv ./kubectl /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl

# RUN JQ_URL="https://circle-downloads.s3.amazonaws.com/circleci-images/cache/linux-amd64/jq-latest" \
#   && curl --silent --show-error --location --fail --retry 3 --output /usr/bin/jq $JQ_URL \
#   && chmod +x /usr/bin/jq \
#   && jq --version

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ingenious ingenious.go
ENV http_proxy=
ENV https_proxy=

#ENTRYPOINT ["/app/vyom"]
