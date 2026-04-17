.PHONY: generate build run dev clean

# templ 코드 생성 (*.templ -> *_templ.go)
generate:
	templ generate

# 빌드
build: generate
	go build -o huginn-go .

# 서버 실행 (기본 포트 3001)
run: build
	./huginn-go

# 개발 모드: templ watch + go run (파일 변경시 자동 재시작)
dev:
	templ generate --watch &
	PORT=3001 go run main.go

# 정리
clean:
	rm -f huginn-go
	find . -name '*_templ.go' -delete

# Go 모듈 업데이트
tidy:
	go mod tidy
