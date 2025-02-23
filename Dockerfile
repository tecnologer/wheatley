FROM golang

WORKDIR /wheatley

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./pkg ./pkg
COPY ./wheatley.db ./wheatley.db

RUN go mod tidy

RUN go build -o wheatley cmd/main.go

CMD ["./wheatley"]
