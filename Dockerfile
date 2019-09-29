
FROM golang

ENV GO111MODULE on
RUN go get github.com/kataras/iris@master

RUN git clone https://github.com/dkzhang/testServerForNiuBen.git

CMD go run main.go