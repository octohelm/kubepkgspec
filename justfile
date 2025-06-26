gengo := "go tool devtool gen"
gofumpt := "go tool gofumpt"

gen:
    {{ gengo }} ./pkg/apis/kubepkg/v1alpha1

test:
    go test -count=1 -v -failfast ./...

fmt:
    {{ gofumpt }} -w -l .
