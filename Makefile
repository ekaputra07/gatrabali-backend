fmt:
	go fmt ./...

test:
	go test -cover -race -count=1 ./...

run_main:
	MINIFLUX_HOST=http://example.com MINIFLUX_USER=user MINIFLUX_PASS=password go run main/main.go

deploy_http:
	gcloud functions deploy Feeds --region=us-central1 --runtime=go111 --trigger-http
	gcloud functions deploy News --region=us-central1 --runtime=go111 --trigger-http
	gcloud functions deploy CategorySummaries --region=us-central1 --runtime=go111 --trigger-http

deploy_pubsub:
	gcloud functions deploy SyncData --runtime=go111 --trigger-topic=SyncData --env-vars-file=env.yaml
