all:
	mkdir -p ./cmd/shakespeare ./rpc ./gapic
	./tools/GENERATE-RPC.sh
	./tools/GENERATE-GRPC.sh
	./tools/GENERATE-GAPIC.sh
	./tools/GENERATE-CLI.sh
	./tools/GENERATE-DOCS.sh
	go install ./...

clean:
	rm -rf cmd/shakespeare/*.go gapic/*.go rpc/*.go third_party envoy/proto.pb docs

protos:
	./tools/GENERATE-ENVOY-DESCRIPTORS.sh

cloudrun:
	gcloud run deploy --source . 

gateway:
	gcloud api-gateway api-configs create shakespeare-config --api=shakespeare --project=$(PROJECT) --grpc-files=proto.pb,api_config.yaml
	gcloud api-gateway gateways create shakespeare --api=shakespeare --api-config=shakespeare-config --location=us-west2 --project=$(PROJECT)
