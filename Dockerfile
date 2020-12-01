FROM golang:alpine as base

WORKDIR /build
RUN apk update && apk upgrade && apk add --no-cache bash git openssh
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/gorilla/schema

COPY . .

RUN go build -o main .
WORKDIR /app
RUN cp /build/main .

FROM alpine
COPY --from=base /app/main /
ENV RANDOM_ORG_API_KEY=your_randomorg_api_key
EXPOSE 80
ENTRYPOINT [ "/main" ]
