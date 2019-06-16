format:
	go fmt ./...

test:
	go test -cover -race -count=1 ./...

deploy_http:
	gcloud functions deploy Hello --region=us-central1 --runtime=go111 --trigger-http

deploy_pubsub:
	gcloud functions deploy SyncData --runtime=go111 --trigger-topic=SyncData
