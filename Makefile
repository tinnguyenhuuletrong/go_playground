build:
	go build -o bin/app

run: build
	./bin/app

test:
	go test -v ./...

run_beam_wc: build
	./bin/app --input gs://dataflow-samples/shakespeare/kinglear.txt --output outputs