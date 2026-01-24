# 注意：下方縮排必須是 Tab 鍵，不能是空格！

.PHONY: fmt run tidy

fmt:
	go fmt ./...

run:
	go run ./cmd/api-server/main.go

tidy:
	go mod tidy -v