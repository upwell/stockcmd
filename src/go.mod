module hehan.net/my/stockcmd

go 1.13

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/containerd/continuity v0.0.0-20200228182428-0f16d7a0959c // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/fatih/color v1.9.0
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/jinzhu/now v1.1.1
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/levigross/grequests v0.0.0-20190908174114-253788527a1a
	github.com/manifoldco/promptui v0.7.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/ory/dockertest v3.3.5+incompatible // indirect
	github.com/pkg/errors v0.8.1
	github.com/rocketlaunchr/dataframe-go v0.0.0-20200331233158-9ec5d1a942f2
	github.com/spf13/cobra v0.0.7
	github.com/stretchr/testify v1.4.0
	go.etcd.io/bbolt v1.3.2
	go.uber.org/zap v1.10.0
	golang.org/x/text v0.3.2
	gonum.org/v1/gonum v0.6.2
	gotest.tools v2.2.0+incompatible
)

replace github.com/olekukonko/tablewriter => ./mydeps/tablewriter
