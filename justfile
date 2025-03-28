gengo := "go tool devtool gen"
gofumpt := "go tool gofumpt"

gen:
    {{ gengo }} ./pkg/apis/kubepkg/v1alpha1

test:
    go test -v -failfast ./...

fmt:
    {{ gofumpt }} -w -l .
