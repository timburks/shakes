all:
	mkdir -p ./cmd/shakes ./rpc ./gapic
	./tools/GENERATE-RPC.sh
	./tools/GENERATE-GRPC.sh
	./tools/GENERATE-GAPIC.sh
	./tools/GENERATE-CLI.sh
	go install ./...
