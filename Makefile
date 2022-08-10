all:
	mkdir -p ./cmd/shakes ./rpc ./gapic
	./tools/GENERATE-RPC.sh
	./tools/GENERATE-GRPC.sh
	./tools/GENERATE-GAPIC.sh
	./tools/GENERATE-CLI.sh
	./tools/GENERATE-DOCS.sh
	go install ./...

clean:
	rm -rf cmd/shakes/*.go gapic/*.go rpc/*.go third_party envoy/proto.pb docs


protos:
	./tools/GENERATE-ENVOY-DESCRIPTORS.sh

cloudrun:
	gcloud run deploy --source . 

gateway:
	gcloud api-gateway api-configs create shakes-config --api=shakes --project=$(PROJECT) --grpc-files=proto.pb,api_config.yaml
	gcloud api-gateway gateways create shakes --api=shakes --api-config=shakes-config --location=us-west2 --project=$(PROJECT)
