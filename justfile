gocache := env_var_or_default("GOCACHE", "/private/tmp/gogen-go-cache")
gomodcache := env_var_or_default("GOMODCACHE", "/private/tmp/gogen-go-mod-cache")
git_commit := `git rev-parse HEAD 2>/dev/null || echo unknown`
git_date := `git show -s --format='%ct' 2>/dev/null || echo 1970-01-01`
project_name := `/usr/bin/env -u GOROOT go list -m | awk -F/ '{print $NF}' || echo gogen`
ldflags := "-X main.GitCommit=" + git_commit + " -X main.GitDate=" + git_date

default: gogen

tidy:
    /usr/bin/env -u GOROOT GOCACHE={{gocache}} GOMODCACHE={{gomodcache}} go mod tidy

gogen: tidy
    /usr/bin/env -u GOROOT GOCACHE={{gocache}} GOMODCACHE={{gomodcache}} go build -v -ldflags "{{ldflags}}" -o {{project_name}} .

clean:
    rm -f gogen
    /usr/bin/env -u GOROOT GOCACHE={{gocache}} GOMODCACHE={{gomodcache}} go clean -cache -testcache

test: tidy
    /usr/bin/env -u GOROOT GOCACHE={{gocache}} GOMODCACHE={{gomodcache}} go test -v ./...

check:
    gofmt -w *.go
    /usr/bin/env -u GOROOT GOCACHE={{gocache}} GOMODCACHE={{gomodcache}} go test ./...
    /usr/bin/env -u GOROOT GOCACHE={{gocache}} GOMODCACHE={{gomodcache}} go vet ./...
    /usr/bin/env -u GOROOT GOCACHE={{gocache}} GOMODCACHE={{gomodcache}} go build ./...

lint: tidy
    golangci-lint run ./...

proto:
    test -f ./bin/compile.sh
    sh ./bin/compile.sh
