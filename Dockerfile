FROM golang:1.23

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY connection ./connection
COPY cmd ./cmd

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=arm

RUN go build -o /bin/polity ./cmd/polity 
RUN go build -o /bin/polityd ./cmd/polityd 

EXPOSE 9005/udp
CMD ["/bin/polityd"]
