GENGO = go run ./internal/cmd/tool gen

gen:
	$(GENGO) ./pkg/apis/kubepkg/v1alpha1

test:
	go test -failfast ./...