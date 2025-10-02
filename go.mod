module github.com/kaiachain/kaia

go 1.23.7

require (
	github.com/Shopify/sarama v1.26.4
	github.com/VictoriaMetrics/fastcache v1.12.2
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137
	github.com/aristanetworks/goarista v0.0.0-20191001182449-186a6201b8ef
	github.com/aws/aws-sdk-go v1.34.28
	github.com/bt51/ntpclient v0.0.0-20140310165113-3045f71e2530
	github.com/cespare/cp v1.0.0
	github.com/clevergo/websocket v1.0.0
	github.com/consensys/gnark-crypto v0.18.0
	github.com/crate-crypto/go-eth-kzg v1.3.0
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.8.0
	github.com/dgraph-io/badger v1.6.0
	github.com/docker/docker v25.0.6+incompatible
	github.com/edsrzf/mmap-go v1.0.0
	github.com/ethereum/c-kzg-4844/v2 v2.1.0
	github.com/fatih/color v1.16.0
	github.com/go-redis/redis/v7 v7.4.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/go-stack/stack v1.8.1
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.5.4
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb
	github.com/gorilla/websocket v1.5.0
	github.com/hashicorp/golang-lru v1.0.2
	github.com/holiman/uint256 v1.3.2
	github.com/huin/goupnp v1.3.0
	github.com/influxdata/influxdb v1.8.3
	github.com/jackpal/go-nat-pmp v1.0.2
	github.com/jinzhu/gorm v1.9.15
	github.com/julienschmidt/httprouter v1.3.0
	github.com/linxGnu/grocksdb v1.10.0
	github.com/mattn/go-colorable v0.1.14
	github.com/mattn/go-isatty v0.0.20
	github.com/naoina/toml v0.1.2-0.20170918210437-9fafd6967416
	github.com/newrelic/go-agent/v3 v3.11.0
	github.com/otiai10/copy v1.0.1
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58
	github.com/peterh/liner v1.1.1-0.20190123174540-a2c9a5303de7
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.15.0
	github.com/prometheus/prometheus v2.1.0+incompatible
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563
	github.com/rjeczalik/notify v0.9.3
	github.com/rs/cors v1.7.0
	github.com/steakknife/bloomfilter v0.0.0-20180922174646-6819c0d2a570
	github.com/steakknife/hamming v0.0.0-20180906055917-c99c65617cd3 // indirect
	github.com/stretchr/testify v1.10.0
	github.com/supranational/blst v0.3.14
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/urfave/cli/v2 v2.27.5
	github.com/valyala/fasthttp v1.40.0
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.36.0
	golang.org/x/net v0.38.0
	golang.org/x/sys v0.31.0
	golang.org/x/tools v0.31.0
	google.golang.org/grpc v1.56.3
	gopkg.in/DataDog/dd-trace-go.v1 v1.42.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
	gopkg.in/fatih/set.v0 v0.1.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20200619000410-60c24ae608a6
	gotest.tools v2.2.0+incompatible
)

require (
	github.com/btcsuite/btcd/btcec/v2 v2.3.4
	github.com/cockroachdb/pebble v1.1.5
	github.com/dop251/goja v0.0.0-20231014103939-873a1496dc8e
	github.com/google/uuid v1.6.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/satori/go.uuid v1.2.0
	github.com/tyler-smith/go-bip32 v1.0.0
	github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4 v1.4.1
	go.uber.org/mock v0.5.0
	golang.org/x/exp v0.0.0-20240318143956-a85f2c67cd81
	golang.org/x/sync v0.12.0
)

require (
	github.com/AndreasBriese/bbloom v0.0.0-20190306092124-e2d15f34fcf9 // indirect
	github.com/BurntSushi/toml v1.4.1-0.20240526193622-a339e1f7089c // indirect
	github.com/DataDog/datadog-agent/pkg/obfuscate v0.0.0-20211129110424-6491aa3bf583 // indirect
	github.com/DataDog/datadog-go v4.8.2+incompatible // indirect
	github.com/DataDog/datadog-go/v5 v5.0.2 // indirect
	github.com/DataDog/sketches-go v1.2.1 // indirect
	github.com/DataDog/zstd v1.4.5 // indirect
	github.com/FactomProject/basen v0.0.0-20150613233007-fe3947df716e // indirect
	github.com/FactomProject/btcutilecc v0.0.0-20130527213604-d3a63a5752ec // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.20.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/fifo v0.0.0-20240606204812-0bbfbd93a7ce // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0 // indirect
	github.com/dgraph-io/ristretto v0.1.0 // indirect
	github.com/dgryski/go-farm v0.0.0-20190423205320-6a90982ecee2 // indirect
	github.com/dlclark/regexp2 v1.7.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/fjl/gencodec v0.1.1 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/garslo/gogen v0.0.0-20170306192744-1d203ffc1f61 // indirect
	github.com/getsentry/sentry-go v0.27.0 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.2.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20230207041349-798e818bf904 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.16.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/otiai10/mint v1.2.4 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/tinylib/msgp v1.1.2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go4.org/intern v0.0.0-20211027215823-ae77deb06f29 // indirect
	go4.org/unsafe/assume-no-moving-gc v0.0.0-20220617031537-928513b29760 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/jcmturner/aescts.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/dnsutils.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/gokrb5.v7 v7.5.0 // indirect
	gopkg.in/jcmturner/rpc.v1 v1.1.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.5.0 // indirect
	inet.af/netaddr v0.0.0-20220617031823-097006376321 // indirect
)
