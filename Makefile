.PHONY: build cleanup deploy

build:
	mkdir "./build" \
	&& GOOS=linux GOARCH=amd64 go build -o build/transcoder ./main.go \
	&& zip -j build/transcoder.zip build/transcoder

cleanup:
	rm -rf ./build

deploy:
	aws lambda \
		update-function-code \
		--profile default \
		--region ap-northeast-1 \
		--function-name transcoder \
		--zip-file fileb://./build/transcoder.zip \
		--publish