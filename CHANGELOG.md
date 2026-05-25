# Changelog

## [1.6.0](https://github.com/twangodev/lfm-cli/compare/v1.5.0...v1.6.0) (2026-05-25)


### Features

* add generic install.sh installer script ([fbee784](https://github.com/twangodev/lfm-cli/commit/fbee784102a86a593096f30391a90553d1b71429))
* add GoReleaser config for multi-platform distribution ([3d5131b](https://github.com/twangodev/lfm-cli/commit/3d5131bbd69917fa77bb3a58e4ffbb5de9b29383))


### Bug Fixes

* address review — persist-credentials:false, reject ARMv6, clarify signing wording ([db35dbb](https://github.com/twangodev/lfm-cli/commit/db35dbbe77c2e515dcaf7cbf9debfb664bdd1102))

## [1.5.0](https://github.com/twangodev/lfm-cli/compare/v1.4.0...v1.5.0) (2026-05-24)


### Features

* use Listening activity type for Discord presence ([2918584](https://github.com/twangodev/lfm-cli/commit/29185848c61b005bac5539bf8ed47835a4030386))

## [1.4.0](https://github.com/twangodev/lfm-cli/compare/v1.3.1...v1.4.0) (2026-05-24)


### Features

* port Discord RPC to go-discordrpc ([cc9cd59](https://github.com/twangodev/lfm-cli/commit/cc9cd59a9e9be6547745528167fd6b987b68c223))

## [1.3.1](https://github.com/twangodev/lfm-cli/compare/v1.3.0...v1.3.1) (2026-01-08)


### Bug Fixes

* reset timestamp to force presence update after Discord reconnection ([7aaeb8e](https://github.com/twangodev/lfm-cli/commit/7aaeb8e21c460dd65954318f7a3a6629d9a2275e))
* output filename format in build script ([9377bd3](https://github.com/twangodev/lfm-cli/commit/9377bd3b151f30210987f93bbefdf863efbb43d6))

## [1.3.0](https://github.com/twangodev/lfm-cli/compare/v1.2.3...v1.3.0) (2026-01-08)


### Features

* enhance Discord connection handling ([a1ae086](https://github.com/twangodev/lfm-cli/commit/a1ae0862dfa88a080f5723e716739355ff68fb78))
* support multi-platform builds for macOS, Windows, Linux, and FreeBSD ([c46e5b6](https://github.com/twangodev/lfm-cli/commit/c46e5b6d932a684591110cf181d534199e8a96fc))


### Bug Fixes

* remove unsupported arm64 build target ([26284d7](https://github.com/twangodev/lfm-cli/commit/26284d71609cb7d09b05460e01f4b056c46e44e8))


### Dependencies

* update Go version and dependencies ([6a3f495](https://github.com/twangodev/lfm-cli/commit/6a3f495432db0b02c6022e33738c69b875b98dc8))
* bump golang.org/x/net to 0.29.0 ([0f69a85](https://github.com/twangodev/lfm-cli/commit/0f69a85316e963de6a9b232ae4e87330c8158dca))
* bump github.com/urfave/cli/v2 to 2.27.3 ([5a61bf4](https://github.com/twangodev/lfm-cli/commit/5a61bf4bfd4af27adbf99b524b7c0418bde7a6fb))

## [1.2.3](https://github.com/twangodev/lfm-cli/compare/v1.2.2...v1.2.3) (2023-10-26)


### Dependencies

* bump golang.org/x/net to 0.17.0 ([4d75b62](https://github.com/twangodev/lfm-cli/commit/4d75b62958bbae76d52f9b150699007f8b1844f7))

## [1.2.2](https://github.com/twangodev/lfm-cli/compare/v1.2.1...v1.2.2) (2023-08-22)


### Bug Fixes

* correct last.fm capitalization ([6f7fbb1](https://github.com/twangodev/lfm-cli/commit/6f7fbb11a6b4d8f96ab01afdeced9839130552ab))


### Dependencies

* bump golang.org/x/net to 0.14.0 ([8551d68](https://github.com/twangodev/lfm-cli/commit/8551d68cc5e1d9eddb6cda4ffab7f3d9ca637e64))
* bump github.com/sirupsen/logrus to 1.9.3 ([484ea72](https://github.com/twangodev/lfm-cli/commit/484ea722c5dcf365ef317ad9929dd225e6e1d989))
* bump github.com/urfave/cli/v2 to 2.25.7 ([ef15c55](https://github.com/twangodev/lfm-cli/commit/ef15c55b10d060676b20685745d94f631e01d725))

## [1.2.1](https://github.com/twangodev/lfm-cli/compare/v1.2.0...v1.2.1) (2023-03-21)


### Dependencies

* bump golang.org/x/net to 0.7.0 ([ad49149](https://github.com/twangodev/lfm-cli/commit/ad491492f72eb7ec57b9cb1c479bb5758ccea6c5))

## [1.2.0](https://github.com/twangodev/lfm-cli/compare/v1.1.5...v1.2.0) (2022-12-12)


### Dependencies

* update to lfm-api v1.0.3 ([0fa796c](https://github.com/twangodev/lfm-cli/commit/0fa796c1ab9012800b075889bc954060b8efad44))

## [1.1.5](https://github.com/twangodev/lfm-cli/compare/v1.1...v1.1.5) (2022-11-25)


### Features

* add go-colorable and change log formatter ([7909d01](https://github.com/twangodev/lfm-cli/commit/7909d01c95a2b8fd642f5a862fae01946c749b18))
* migrate to lfm-api package ([3ec0038](https://github.com/twangodev/lfm-cli/commit/3ec003824bd0bd8c28d409011423c056a4112092))


### Dependencies

* update to lfm-api v1.0.2 ([e3fdfe8](https://github.com/twangodev/lfm-cli/commit/e3fdfe8f0be55de73efa45f1a835447be4fe636d))

## [1.1](https://github.com/twangodev/lfm-cli/compare/v1.0...v1.1) (2022-10-24)


### Bug Fixes

* close IPC socket upon login ([ab1b94e](https://github.com/twangodev/lfm-cli/commit/ab1b94ee254b0742476402ea83ce2b79675c228f))
* logging for IPC socket ([61871c9](https://github.com/twangodev/lfm-cli/commit/61871c9a88285b41d8d1a6522ca1c6a41c181fac))

## 1.0 (2022-10-24)


### Features

* add keepStatus functionality ([e56846f](https://github.com/twangodev/lfm-cli/commit/e56846fc1630efe59036b3bac7e3b1a85f714258))
