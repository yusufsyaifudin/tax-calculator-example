FROM golang:1.11.1-alpine

# since we need Makefile command
RUN apk add --update make
RUN apk add --update git
RUN mkdir -p /go/src/github.com/yusufsyaifudin/tax-calculator-example
ADD . /go/src/github.com/yusufsyaifudin/tax-calculator-example
WORKDIR /go/src/github.com/yusufsyaifudin/tax-calculator-example

# delete the vendor directory to make sure that this docker images can successfully build from scratch without using host dependency
RUN rm -rf vendor

# run Makefile command
RUN make install

# binary will be located in /go/src/github.com/yusufsyaifudin/tax-calculator-example/out/tax-calculator-server
CMD ["/go/src/github.com/yusufsyaifudin/tax-calculator-example/out/tax-calculator-server"]
