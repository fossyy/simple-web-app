FROM golang:alpine3.19 AS builder

WORKDIR /src
COPY . .

RUN go build -o /bin/hello ./main.go

FROM scratch

COPY --from=builder /bin/hello /bin/hello

LABEL maintainer="Bagas <bagas@fossy.my.id>"
LABEL description="Simple GO Web Application"

CMD ["/bin/hello"]
