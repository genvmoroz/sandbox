module sandbox

go 1.19

require (
	github.com/90poe/atlantic/crewing-domain-service v0.145.0
	github.com/90poe/chartering-domain-service v0.15.1
	github.com/90poe/performance/vessel-performance-information-service v1.1.1
	github.com/90poe/port-information-domain-service/v3 v3.10.5
	github.com/90poe/service-chassis/grpc/v3 v3.8.0
	github.com/90poe/service-chassis/m2m/v2 v2.0.1
	github.com/90poe/vessel-information-domain-service/v5 v5.64.2
	github.com/90poe/vessel-itinerary-domain-service/v3 v3.41.0
	github.com/90poe/voyage-monitor/reports-service v0.5.15
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/kellydunn/golang-geo v0.7.0
	github.com/xuri/excelize/v2 v2.7.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/90poe/operational-effectiveness/vessel-position-service/v2 v2.5.1 // indirect
	github.com/90poe/ps/m2m-auth-idp-service/api/v2 v2.0.2 // indirect
	github.com/90poe/service-chassis/authorisation/v2 v2.1.3 // indirect
	github.com/90poe/service-chassis/context v1.3.6 // indirect
	github.com/90poe/service-chassis/correlation v1.4.1 // indirect
	github.com/90poe/service-chassis/jsonx v1.3.0 // indirect
	github.com/90poe/service-chassis/logging v1.3.3 // indirect
	github.com/90poe/service-chassis/logging/v2 v2.2.1 // indirect
	github.com/90poe/service-chassis/logging/v4 v4.0.0-beta // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/gbrlsnchs/jwt/v2 v2.0.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/goccy/go-json v0.9.6 // indirect
	github.com/golang/geo v0.0.0-20181008215305-476085157cff // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/kylelemons/go-gypsy v1.0.0 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.1 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx v1.2.21 // indirect
	github.com/lestrrat-go/option v1.0.0 // indirect
	github.com/lib/pq v1.10.4 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.33.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.3 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/xuri/efp v0.0.0-20220603152613-6918739fd470 // indirect
	github.com/xuri/nfp v0.0.0-20220409054826-5e722a1d9e22 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.31.0 // indirect
	go.opentelemetry.io/otel v1.7.0 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20220714211235-042d03aeabc9 // indirect
	google.golang.org/grpc v1.49.0 // indirect
)

replace (
	cloud.google.com/go => cloud.google.com/go v0.102.1
	github.com/90poe/m2m-auth-idp-service => github.com/90poe/ps/m2m-auth-idp-service v0.2.1
	github.com/90poe/platform-remarks-service/v2 => github.com/90poe/ps/platform-remarks-service/v2 v2.19.2
	github.com/90poe/service-chassis/http => github.com/90poe/service-chassis/http v1.4.3
	github.com/90poe/user-service/v4 => github.com/90poe/ps/user-service/v4 v4.11.2
	github.com/golangci/golangci-lint => github.com/golangci/golangci-lint v1.46.2
)
