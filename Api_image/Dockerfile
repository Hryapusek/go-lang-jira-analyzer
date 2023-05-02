FROM golang:1.20

WORKDIR /usr/src/api

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY ApiServer/go.mod ApiServer/go.sum ./
RUN go mod download && go mod verify

COPY ApiServer .
RUN go build -v ./cmd/main/gateway.go
RUN go build -v ./cmd/main/analytics.go
RUN go build -v ./cmd/main/resource.go

WORKDIR ./cmd/main/
CMD /usr/src/api/cmd/main/gateway & /usr/src/api/cmd/main/resource & /usr/src/api/cmd/main/analytics

EXPOSE 8000