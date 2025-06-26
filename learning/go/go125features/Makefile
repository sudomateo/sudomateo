CPU_QUOTA ?= 100%

.PHONY: go-version-124
go-version-124:
	go mod edit -go=1.24.4

.PHONY: go-version-125
go-version-125:
	go mod edit -go=1.25rc1

.PHONY: gomaxprocs
gomaxprocs:
	go build ./example/gomaxprocs
	sudo systemd-run --slice user --scope --unit gomaxprocs ./gomaxprocs

.PHONY: gomaxprocs-set-cpu-quota
gomaxprocs-set-cpu-quota:
	sudo systemctl set-property gomaxprocs.scope CPUQuota=$(CPU_QUOTA)

.PHONY: jsonv2-general
jsonv2-general:
	@GOEXPERIMENT=jsonv2 go run ./example/jsonv2/general

.PHONY: jsonv2-tags_inline
jsonv2-tags_inline:
	@GOEXPERIMENT=jsonv2 go run ./example/jsonv2/tags_inline

.PHONY: jsonv2-tags_unknown
jsonv2-tags_unknown:
	@GOEXPERIMENT=jsonv2 go run ./example/jsonv2/tags_unknown

.PHONY: jsonv2-custom
jsonv2-custom:
	@GOEXPERIMENT=jsonv2 go run ./example/jsonv2/custom

.PHONY: synctest
synctest:
	go test -v ./example/gomaxprocs
