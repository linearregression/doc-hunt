export CGO_ENABLED=1

compile:
	rm -rf build
	mkdir build
	git stash -u
	gox -osarch=linux/386 -output "build/{{.Dir}}_{{.OS}}_{{.Arch}}"
	gox -osarch=linux/amd64 -output "build/{{.Dir}}_{{.OS}}_{{.Arch}}"

version:
	git stash -u
	sed -i "s/[[:digit:]]\+\.[[:digit:]]\+\.[[:digit:]]\+/$(v)/g" file/version.go
	git add -A
	git commit -m "feat(version) : "$(v)
	git tag v$(v) master

fmt:
	find ! -path "./vendor/*" -name "*.go" -exec go fmt {} \;

gometalinter:
	gometalinter -D gotype --vendor --deadline=240s --dupl-threshold=200 -e '_string' -j 5 ./...

doc-hunt:
	doc-hunt check -e

run-tests:
	./test.sh

test-all: gometalinter run-tests doc-hunt

test-package:
	go test -race -cover -coverprofile=/tmp/doc-hunt github.com/antham/doc-hunt/$(pkg)
	go tool cover -html=/tmp/doc-hunt
