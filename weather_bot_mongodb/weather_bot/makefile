runbot:
	go run cmd/main.go
buildbot:
	go build -o main cmd/main.go
runbotwith:
	go run cmd/with_ctx/main.go
mock:
	~/go/bin/mockgen -source=internal/repositories/user_repo.go -destination=./internal/mock/male_mock.go -package=mock
