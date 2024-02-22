module hehan.net/my/stockcmd

go 1.18

require (
	github.com/deckarep/golang-set v1.8.0
	github.com/fatih/color v1.13.0
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/iancoleman/strcase v0.2.0
	github.com/jinzhu/now v1.1.5
	github.com/levigross/grequests v0.0.0-20190908174114-253788527a1a
	github.com/manifoldco/promptui v0.9.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/errors v0.9.1
	github.com/rocketlaunchr/dataframe-go v0.0.0-20211025052708-a1030444159b
	github.com/silenceper/pool v1.0.0
	github.com/spf13/cobra v1.5.0
	github.com/stretchr/testify v1.7.0
	go.etcd.io/bbolt v1.3.6
	go.uber.org/zap v1.21.0
	golang.org/x/text v0.8.0
	gonum.org/v1/gonum v0.11.0
	gotest.tools v2.2.0+incompatible
)

require (
	github.com/apache/arrow/go/arrow v0.0.0-20211112161151-bc219186db40 // indirect
	github.com/apache/thrift v0.16.0 // indirect
	github.com/chzyer/readline v1.5.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/goccy/go-json v0.9.8 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/guptarohit/asciigraph v0.5.5 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/juju/clock v1.0.0 // indirect
	github.com/juju/errors v1.0.0 // indirect
	github.com/juju/loggo v1.0.0 // indirect
	github.com/juju/utils/v2 v2.0.0-20210305225158-eedbe7b6b3e2 // indirect
	github.com/klauspost/compress v1.15.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/opencontainers/runc v1.1.12 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rocketlaunchr/mysql-go v1.1.3 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/upwell/go-wcwidth v0.0.3 // indirect
	github.com/xitongsys/parquet-go v1.6.2 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/exp v0.0.0-20220706164943-b4a6d9510983 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/olekukonko/tablewriter => ./mydeps/tablewriter

replace github.com/silenceper/pool => ./mydeps/pool
