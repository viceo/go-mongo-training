FROM golang
WORKDIR /go/src/app/github.com/viceo/go-mongodb
COPY . .

RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

CMD [ "watcher" ]

