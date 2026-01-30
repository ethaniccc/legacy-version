module github.com/ethaniccc/legacy-version

go 1.25.1

require (
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/df-mc/dragonfly v0.10.9
	github.com/df-mc/worldupgrader v1.0.20
	github.com/go-gl/mathgl v1.2.0
	github.com/google/uuid v1.6.0
	github.com/hashicorp/go-version v1.8.0
	github.com/pelletier/go-toml v1.9.5
	github.com/rogpeppe/go-internal v1.14.1
	github.com/samber/lo v1.49.1
	github.com/sandertv/gophertunnel v1.52.0
	github.com/segmentio/fasthash v1.0.3
	golang.org/x/exp v0.0.0-20251209150349-8475f28825e9
	golang.org/x/image v0.25.0
	golang.org/x/oauth2 v0.34.0
)

require (
	github.com/df-mc/go-nethernet v0.0.0-20250326113854-da40ae9a1339 // indirect
	github.com/df-mc/go-playfab v1.0.0 // indirect
	github.com/df-mc/go-xsapi v1.0.1 // indirect
	github.com/df-mc/goleveldb v1.1.9 // indirect
	github.com/getsentry/sentry-go v0.40.0 // indirect
	github.com/go-jose/go-jose/v4 v4.1.3 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/lumineproxy/log v0.0.0-20251025091855-37df1c504536 // indirect
	github.com/lumineproxy/socks5udp v1.0.1-0.20251223034400-3f3b0cbbe69b // indirect
	github.com/muhammadmuzzammil1998/jsonc v1.0.0 // indirect
	github.com/pion/datachannel v1.5.10 // indirect
	github.com/pion/dtls/v3 v3.0.9 // indirect
	github.com/pion/ice/v4 v4.1.0 // indirect
	github.com/pion/interceptor v0.1.42 // indirect
	github.com/pion/logging v0.2.4 // indirect
	github.com/pion/mdns/v2 v2.1.0 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/rtcp v1.2.16 // indirect
	github.com/pion/rtp v1.8.27 // indirect
	github.com/pion/sctp v1.8.41 // indirect
	github.com/pion/sdp/v3 v3.0.17 // indirect
	github.com/pion/srtp/v3 v3.0.9 // indirect
	github.com/pion/stun/v3 v3.0.2 // indirect
	github.com/pion/transport/v3 v3.1.1 // indirect
	github.com/pion/turn/v4 v4.1.3 // indirect
	github.com/pion/webrtc/v4 v4.1.8 // indirect
	github.com/sandertv/go-raknet v1.14.3-0.20250823121252-325aeea25d25 // indirect
	github.com/wlynxg/anet v0.0.5 // indirect
	github.com/zaataylor/cartesian v0.0.0-20221028053253-3b3244d82727 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/mod v0.31.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)

replace (
	github.com/df-mc/dragonfly => github.com/lumineproxy/dragonfly v0.0.0-20260101232342-710c0a1931ba
	github.com/df-mc/go-nethernet => github.com/lumineproxy/go-nethernet v0.0.0-20251201013730-45766a8ec674
	github.com/sandertv/go-raknet => github.com/lumineproxy/go-raknet v0.0.0-20260102015805-b139417c724e
	github.com/sandertv/gophertunnel => github.com/lumineproxy/gophertunnel v0.0.0-20260102091803-5fb43dd59e0c
)
