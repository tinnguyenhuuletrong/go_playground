build:
	go build -o bin/app

run: build
	./bin/app

test:
	go test -v ./...

run_beam_wc: build
	./bin/app --input gs://dataflow-samples/shakespeare/kinglear.txt --output outputs

build_protoc:
	# Need install first
	# $ go install github.com/twitchtv/twirp/protoc-gen-twirp@latest
	# $ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	protoc --twirp_out=./ --go_out=./ ./grpc_play/notes/service.proto

build_cel_protoc:
	protoc --go_out=./ ./goo_cel_play/type.proto