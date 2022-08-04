FROM golang:latest

LABEL maintainer="Quique <hello@bookApi.com>"

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

ENV PORT 4001

EXPOSE 4001
# RUN go install github.com/gobuffalo/pop/v6/soda@latest

RUN cd api && ls && go build && ls && cd .. && ls

# RUN cd migrations && echo 'alias soda="~/go/bin/soda"' >> ~/.bashrc &&  soda migrate && cd ..


# RUN go build

# RUN cd ..

CMD [ "./api/api" ]