PROTO=./protos/example.proto
PB_INCLUDES=-I./protos
GO_OUT=./pb


# all: sync_protos go python
all: go

go: go_pb

go_pb: $(PROTO)
	protoc $(PB_INCLUDES) $(PROTO) --go_out=plugins=grpc:$(GO_OUT)

clean:
	rm $(PB_OUT)/*