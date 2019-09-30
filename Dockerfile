
FROM golang

RUN go get github.com/kataras/iris

RUN git clone https://github.com/dkzhang/testServerForNiuBen.git  && \
 cd testServerForNiuBen  && \
 go build ./main.go

CMD ./main