FROM golang:1.14.4

RUN mkdir /go/src/shops/backend

ADD . /go/src/shops/backend

WORKDIR /go/src/shops/backend

ENV GO111MODULE=on

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /shops .

EXPOSE 8080

FROM scratch

COPY --from=builder /shops ./

ENTRYPOINT ["./shops" docker build -t my-golang-app .]
