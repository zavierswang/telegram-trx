.PHONY: all build clean run check cover lint docker help

dateTime=`date +%F_%T`
ARCH="linux-amd64"

all: build

build:
	xgo -targets=linux/amd64 -ldflags="-w -s" -out=./build/telegram-trx -pkg=cmd/telegram-trx/main.go .
	upx ./build/telegram-trx-${ARCH}
	tar czf build/telegram-trx_${dateTime}.tar.gz \
		build/telegram-trx-${ARCH} \
 		assets \
 		template \
 		telegram-trx.yaml.example

