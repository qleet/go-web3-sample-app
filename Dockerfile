FROM golang:1.19-alpine as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN export CGO_ENABLED=0; go build -a -o go-web3-sample-app main.go

FROM scratch
COPY --from=builder /app/go-web3-sample-app /
COPY --from=builder /app/*.html /
EXPOSE 8080
ENTRYPOINT [ "/go-web3-sample-app" ]
