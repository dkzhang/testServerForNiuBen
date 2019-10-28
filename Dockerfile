
FROM golang

RUN go get github.com/kataras/iris

WORKDIR /go/src
RUN git clone https://github.com/dkzhang/testServerForNiuBen.git #git20191028

WORKDIR /go/src/testServerForNiuBen
RUN go build ./main.go

EXPOSE 8080

CMD ./main