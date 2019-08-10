fmt:
	go fmt ./...

test:
	go test -cover -race -count=1 ./...

run_main:
	MINIFLUX_HOST=http://example.com MINIFLUX_USER=user MINIFLUX_PASS=password go run main/main.go

deploy_http_handler:
	gcloud functions deploy Hello --runtime=go111 --trigger-http

deploy_pubsub_handler:
	gcloud functions deploy SyncData --runtime=go111 --trigger-topic=SyncData --env-vars-file=env.yaml
	gcloud functions deploy SendPushNotification --runtime=go111 --trigger-topic=PushNotification

deploy_doc_handler:
	gcloud functions deploy NotifyCategorySubscribers-main --entry-point NotifyCategorySubscribers --runtime=go111 --trigger-event=providers/cloud.firestore/eventTypes/document.create --trigger-resource="projects/gatrabali/databases/(default)/documents/entries/{entryId}"
	gcloud functions deploy NotifyCategorySubscribers-kriminal --entry-point NotifyCategorySubscribers --runtime=go111 --trigger-event=providers/cloud.firestore/eventTypes/document.create --trigger-resource="projects/gatrabali/databases/(default)/documents/kriminal/{entryId}"
	gcloud functions deploy NotifyCategorySubscribers-baliunited --entry-point NotifyCategorySubscribers --runtime=go111 --trigger-event=providers/cloud.firestore/eventTypes/document.create --trigger-resource="projects/gatrabali/databases/(default)/documents/baliunited/{entryId}"
	gcloud functions deploy NotifyCategorySubscribers-balebengong --entry-point NotifyCategorySubscribers --runtime=go111 --trigger-event=providers/cloud.firestore/eventTypes/document.create --trigger-resource="projects/gatrabali/databases/(default)/documents/balebengong/{entryId}"
	gcloud functions deploy AggregateEntryResponses --entry-point AggregateEntryResponses --runtime=go111 --trigger-event=providers/cloud.firestore/eventTypes/document.write --trigger-resource="projects/gatrabali/databases/(default)/documents/entry_responses/{responseId}"