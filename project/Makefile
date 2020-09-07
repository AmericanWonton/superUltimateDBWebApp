run:
	go run *.go

gomod-exp:
	export GO111MODULE=on

gobuild:
	go build
dockerbuild:
	docker build -t suberdb .
dockerbuildandpush:
	docker build -t suberdb .
	docker tag suberdb americanwonton/suberdb
	docker push americanwonton/suberdb
dockerrun:
	docker run -it -p 8080:8080 suberdb
dockerrundetached:
	docker run -d -p 8080:8080 suberdb
dockerrunitvolume:
	docker run -it -p 8080:8080 -v photo-images:/static/images suberdb
dockerrundetvolume:
	docker run -d -p 8080:8080 -v photo-images:/static/images suberdb
dockertagimage:
	docker tag suberdb americanwonton/suberdb
dockerimagepush:
	docker push americanwonton/suberdb
dockerallpush:
	docker tag suberdb americanwonton/suberdb
	docker push americanwonton/suberdb
dockerseeshell:
	docker run -it suberdb sh