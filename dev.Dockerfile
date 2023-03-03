FROM golang:alpine

WORKDIR /code

COPY . /code
RUN go mod tidy

CMD [ "go run ." ]