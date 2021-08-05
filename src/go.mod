module hehan.net/my/stockcmd

go 1.16

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/deckarep/golang-set v1.7.1
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/fatih/color v1.9.0
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/guptarohit/asciigraph v0.5.2 // indirect
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/jinzhu/now v1.1.2
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/juju/ansiterm v0.0.0-20210706145210-9283cdf370b5 // indirect
	github.com/juju/loggo v0.0.0-20210728185423-eebad3a902c4 // indirect
	github.com/juju/utils/v2 v2.0.0-20210305225158-eedbe7b6b3e2 // indirect
	github.com/klauspost/compress v1.10.11 // indirect
	github.com/levigross/grequests v0.0.0-20190908174114-253788527a1a
	github.com/manifoldco/promptui v0.7.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/olekukonko/tablewriter v0.0.5
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/ory/dockertest v3.3.5+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/rocketlaunchr/dataframe-go v0.0.0-20210422123815-aaa951b82e1b
	github.com/silenceper/pool v0.0.0-20200429081406-a659d818d9aa
	github.com/spf13/cobra v0.0.7
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/tealeg/xlsx v1.0.5 // indirect
	github.com/xitongsys/parquet-go v1.5.4 // indirect
	go.etcd.io/bbolt v1.3.6
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/exp v0.0.0-20210729172720-737cce5152fc // indirect
	golang.org/x/net v0.0.0-20210726213435-c6fcb2dbf985 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/text v0.3.6
	gonum.org/v1/gonum v0.7.0
	gotest.tools v2.2.0+incompatible
)

replace github.com/olekukonko/tablewriter => ./mydeps/tablewriter

replace github.com/silenceper/pool => ./mydeps/pool
