.PHONY: build run test race

build:
	go build .

run: build
	./counting_ip $(ARGS)

test:
	python3 ./generate_test_data.py $(GENERATED_IPS)
	go test -v .

race:
	python3 ./generate_test_data.py $(GENERATED_IPS)
	go test -race -v .