
FROM golang

RUN go get github.com/kataras/iris

WORKDIR /go/src
RUN git clone https://github.com/dkzhang/testServerForNiuBen.git

WORKDIR /go/src/testServerForNiuBen
RUN go build ./main.go

CMD ./main