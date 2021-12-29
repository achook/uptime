FROM golang:1.17-alpine AS build

WORKDIR /uptime
COPY ./* .
RUN go get ./...
RUN go build -o uptime *.go

FROM alpine:latest AS production
WORKDIR /uptime

COPY --from=build /uptime/uptime .
COPY --from=build /uptime/service_account.json .

ENV GOOGLE_APPLICATION_CREDENTIALS="/uptime/service_account.json"
ENV PROJECT_ID="goroczas"

CMD ["./uptime"]