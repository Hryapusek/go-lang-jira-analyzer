FROM golang:1.20

WORKDIR /usr/src/connector

COPY JiraConnector .
RUN go mod download && go mod verify

WORKDIR /usr/src/connector/cmd
RUN go build -v ./main.go

EXPOSE 8003

CMD /usr/src/connector/cmd/main