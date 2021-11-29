FROM golang:latest
ENV http_proxy 10.158.100.6:8080
ENV https_proxy 10.158.100.6:8080
RUN apt-get install curl

RUN JQ_URL="https://circle-downloads.s3.amazonaws.com/circleci-images/cache/linux-amd64/jq-latest" \
  && curl --silent --show-error --location --fail --retry 3 --output /usr/bin/jq $JQ_URL \
  && chmod +x /usr/bin/jq \
  && jq --version

ENV http_proxy= 
ENV https_proxy= 
WORKDIR /app
COPY mediationxapp2 ./
EXPOSE 9002
# CMD ["./testapp2"]
