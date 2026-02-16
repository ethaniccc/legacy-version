module github.com/ethaniccc/legacy-version

go 1.25.1

require (
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/df-mc/dragonfly v0.10.11-0.20260211105526-1ec432bade04
	github.com/df-mc/worldupgrader v1.0.20
	github.com/go-gl/mathgl v1.2.0
	github.com/google/uuid v1.6.0
	github.com/hashicorp/go-version v1.8.0
	github.com/pelletier/go-toml v1.9.5
	github.com/rogpeppe/go-internal v1.14.1
	github.com/samber/lo v1.49.1
	github.com/sandertv/gophertunnel v1.54.0
	github.com/segmentio/fasthash v1.0.3
	golang.org/x/exp v0.0.0-20260211191109-2735e65f0518
	golang.org/x/image v0.25.0
	golang.org/x/oauth2 v0.35.0
)

require (
	github.com/coreos/go-oidc/v3 v3.17.0 // indirect
	github.com/df-mc/go-playfab v1.0.0 // indirect
	github.com/df-mc/go-xsapi v1.0.1 // indirect
	github.com/df-mc/goleveldb v1.1.9 // indirect
	github.com/df-mc/jsonc v1.0.5 // indirect
	github.com/go-jose/go-jose/v4 v4.1.3 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.17.0 // indirect
	github.com/sandertv/go-raknet v1.15.0 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/zaataylor/cartesian v0.0.0-20221028053253-3b3244d82727 // indirect
	golang.org/x/mod v0.33.0 // indirect
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/text v0.34.0 // indirect
)

replace (
	github.com/df-mc/dragonfly => github.com/lumineproxy/dragonfly v0.0.0-20260101232342-710c0a1931ba
	github.com/df-mc/go-nethernet => github.com/lumineproxy/go-nethernet v0.0.0-20251201013730-45766a8ec674
	github.com/sandertv/go-raknet => github.com/lumineproxy/go-raknet v0.0.0-20260102015805-b139417c724e
)
