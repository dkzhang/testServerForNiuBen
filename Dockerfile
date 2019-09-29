
FROM golang

ENV GO111MODULE on
RUN go get github.com/kataras/iris@master

RUN git clone https://github.com/dkzhang/testServerForNiuBen.git  && \
 go build main.go

CMD ./main