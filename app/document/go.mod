module github.com/MoScenix/industrial-fault-tree-ai/app/document

go 1.25.5

replace (
	github.com/MoScenix/industrial-fault-tree-ai/common => ../../common
	github.com/MoScenix/industrial-fault-tree-ai/rpc_gen => ../../rpc_gen
	github.com/apache/thrift => github.com/apache/thrift v0.13.0
)

require (
	github.com/MoScenix/industrial-fault-tree-ai/common v0.0.0-00010101000000-000000000000
	github.com/MoScenix/industrial-fault-tree-ai/rpc_gen v0.0.0-00010101000000-000000000000
	github.com/cloudwego/eino v0.7.19
	github.com/cloudwego/eino-ext/components/embedding/openai v0.0.0-20260331071634-4f359694d2d9
	github.com/cloudwego/eino-ext/components/indexer/milvus v0.0.0-20260331071634-4f359694d2d9
	github.com/cloudwego/eino-ext/components/retriever/milvus v0.0.0-20260331071634-4f359694d2d9
	github.com/cloudwego/kitex v0.15.4
	github.com/joho/godotenv v1.5.1
	github.com/kitex-contrib/obs-opentelemetry/logging/logrus v0.0.0-20251121033812-f6c3e41f13e9
	github.com/kr/pretty v0.3.1
	github.com/ledongthuc/pdf v0.0.0-20250511090121-5959a4027728
	github.com/milvus-io/milvus-sdk-go/v2 v2.4.2
	go.uber.org/zap v1.27.1
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/validator.v2 v2.0.1
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/armon/go-metrics v0.3.10 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/bytedance/gopkg v0.1.3 // indirect
	github.com/bytedance/sonic v1.14.2 // indirect
	github.com/bytedance/sonic/loader v0.4.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/cloudwego/configmanager v0.2.3 // indirect
	github.com/cloudwego/dynamicgo v0.7.1 // indirect
	github.com/cloudwego/eino-ext/libs/acl/openai v0.1.15 // indirect
	github.com/cloudwego/fastpb v0.0.5 // indirect
	github.com/cloudwego/frugal v0.3.0 // indirect
	github.com/cloudwego/gopkg v0.1.8 // indirect
	github.com/cloudwego/localsession v0.2.1 // indirect
	github.com/cloudwego/netpoll v0.7.2 // indirect
	github.com/cloudwego/runtimex v0.1.1 // indirect
	github.com/cloudwego/thriftgo v0.4.3 // indirect
	github.com/cockroachdb/errors v1.9.1 // indirect
	github.com/cockroachdb/logtags v0.0.0-20211118104740-dabe8e521a4f // indirect
	github.com/cockroachdb/redact v1.1.3 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/eino-contrib/jsonschema v1.0.3 // indirect
	github.com/evanphx/json-patch v0.5.2 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/getsentry/sentry-go v0.12.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/pprof v0.0.0-20240727154555-813a5fbdbec8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/goph/emperror v0.17.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/hashicorp/consul/api v1.20.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/jhump/protoreflect v1.8.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kitex-contrib/monitor-prometheus v0.2.0 // indirect
	github.com/kitex-contrib/registry-consul v0.2.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/meguminnnnnnnnn/go-openai v0.1.2 // indirect
	github.com/milvus-io/milvus-proto/go-api/v2 v2.4.10-0.20240819025435-512e3b98866a // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nikolalohinski/gonja v1.5.3 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.17.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/slongfield/pyfmt v0.0.0-20220222012616-ea85ff4c361f // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/tidwall/gjson v1.17.3 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	github.com/yargevad/filepathx v1.0.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/arch v0.15.0 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.73.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
