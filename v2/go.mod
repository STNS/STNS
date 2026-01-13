module github.com/STNS/STNS/v2

go 1.24.0

replace github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.0+incompatible

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/aws/aws-sdk-go-v2/config v1.32.6
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.20.20
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.19.15
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.53.5
	github.com/aws/aws-sdk-go-v2/service/s3 v1.95.0
	github.com/facebookgo/pidfile v0.0.0-20150612191647-f242e2999868
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/iancoleman/strcase v0.3.0
	github.com/jpillora/ipfilter v1.2.9
	github.com/labstack/echo/v4 v4.13.4
	github.com/labstack/gommon v0.4.2
	github.com/lestrrat/go-server-starter v0.0.0-20180220115249-6ac0b358431b
	github.com/lib/pq v1.10.9
	github.com/nmcclain/ldap v0.0.0-20210720162743-7f8d1e44eeba
	github.com/stretchr/testify v1.11.1
	github.com/tredoe/osutil v1.5.0
	github.com/urfave/cli v1.22.17
	go.etcd.io/etcd/client/v2 v2.305.26
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/redis.v5 v5.2.9
)

require (
	github.com/aws/aws-sdk-go-v2 v1.41.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.4 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.6 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.32.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.5 // indirect
	github.com/aws/smithy-go v1.24.0 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/facebookgo/atomicfile v0.0.0-20151019160806-2de1f203e7d5 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/nmcclain/asn1-ber v0.0.0-20170104154839-2661553a0484 // indirect
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.16.0 // indirect
	github.com/phuslu/iploc v1.0.20230201 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.etcd.io/etcd/api/v3 v3.5.26 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.26 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
