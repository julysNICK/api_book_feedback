FROM golang:latest

LABEL maintainer="Quique <hello@bookApi.com>"

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

ENV PORT 4001

RUN cd api && ls && go build && ls && cd .. && ls

# RUN go build

# RUN cd ..

CMD [ "./api/api" ]