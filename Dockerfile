#FROM golang:alpine
#WORKDIR /app
#
#RUN go version
#ENV GOPATH=/
#
#COPY ./ ./
#
#RUN go mod download
#RUN go build -o app ./cmd/app/main.go
#
#CMD ["./app"]

FROM golang:alpine AS builder

WORKDIR /app

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o server ./cmd/app/main.go

FROM alpine:latest

COPY --from=builder /app/server ./
COPY --from=builder /app/migrations ./migrations
COPY family-flow-app-4900a-firebase-adminsdk-fbsvc-9e2d528ca1.json ./family-flow-app-4900a-firebase-adminsdk-fbsvc-9e2d528ca1.json
COPY family-flow-app-4900a-firebase-adminsdk-fbsvc-4aa2ce52fc.json ./family-flow-app-4900a-firebase-adminsdk-fbsvc-4aa2ce52fc.json
COPY config/config.yaml ./config/config.yaml
COPY .aws/ .aws/

# COPY .aws/ ./
# COPY local.env/ ./

EXPOSE 8080

CMD ["./server", "--port", "8080"]