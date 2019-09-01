FROM golang:1.12.9 as builder
LABEL maintainer="Karthik Raveendran <karthik.panicker@gmail.com>"
WORKDIR /etsello
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest 
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /etsello/main .
COPY --from=builder /etsello/.env .
COPY --from=builder /etsello/templates ./templates
EXPOSE 80
CMD ["./main"]