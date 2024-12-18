# Changelog

## [0.2.2](https://github.com/abemedia/go-don/compare/v0.2.1...v0.2.2) (2024-12-18)


### Bug Fixes

* **encoding:** take the first encoding supported ([#205](https://github.com/abemedia/go-don/issues/205)) ([3f53fd0](https://github.com/abemedia/go-don/commit/3f53fd061cacbc748e891baf998a91a972a5b37d))

## [0.2.1](https://github.com/abemedia/go-don/compare/v0.2.0...v0.2.1) (2023-06-04)


### Bug Fixes

* content-length is 0 when responding with nil value ([#148](https://github.com/abemedia/go-don/issues/148)) ([9cf40b2](https://github.com/abemedia/go-don/commit/9cf40b2c8a36468d072c3b17cebb2dc745c0e520))


### Performance Improvements

* **encoding/text:** reduce allocs ([#142](https://github.com/abemedia/go-don/issues/142)) ([5759715](https://github.com/abemedia/go-don/commit/575971580296f0acaa929c6e849bbe707e580165))
* improve pool performance for pointer types ([#147](https://github.com/abemedia/go-don/issues/147)) ([d9464de](https://github.com/abemedia/go-don/commit/d9464deb560eac1500b51f45c5cfe64822efed63))

## [0.2.0](https://github.com/abemedia/go-don/compare/v0.1.4...v0.2.0) (2023-05-14)


### ⚠ BREAKING CHANGES

* remove Empty type (use any instead) ([#135](https://github.com/abemedia/go-don/issues/135))
* move encoding logic to sub-package ([#111](https://github.com/abemedia/go-don/issues/111))

### Features

* **encoding:** support more media type aliases ([#115](https://github.com/abemedia/go-don/issues/115)) ([d88115c](https://github.com/abemedia/go-don/commit/d88115c058e6d81c9fd0ec1d27d55bd44b4cf8e6))
* **encoding:** support protocol buffers ([#117](https://github.com/abemedia/go-don/issues/117)) ([ace6006](https://github.com/abemedia/go-don/commit/ace600620fbe9c67e56ecfb1b7394536cc1da0a4))
* **encoding:** support toml ([#114](https://github.com/abemedia/go-don/issues/114)) ([e95b4ae](https://github.com/abemedia/go-don/commit/e95b4aed2a43c5bd87dbf3bb4591faf0d0fd3c97))


### Bug Fixes

* no content issues on error or pointer to Empty ([#125](https://github.com/abemedia/go-don/issues/125)) ([fa50d36](https://github.com/abemedia/go-don/commit/fa50d363e872d51baeed84cb516d6d4a45fc345b))


### Performance Improvements

* **encoding/protobuf:** cache reflection results ([#138](https://github.com/abemedia/go-don/issues/138)) ([99e0cea](https://github.com/abemedia/go-don/commit/99e0cea46d5e42e91dda63bfdd365835161a9a03))
* **encoding:** improve encoding performance ([#113](https://github.com/abemedia/go-don/issues/113)) ([a541544](https://github.com/abemedia/go-don/commit/a541544614d07121266a2ebf1eebfd75b9d7541d))
* reuse requests to reduce allocs ([#127](https://github.com/abemedia/go-don/issues/127)) ([827209b](https://github.com/abemedia/go-don/commit/827209bca6cfa7a91c414f6bced4a10308d9573f))


### Code Refactoring

* move encoding logic to sub-package ([#111](https://github.com/abemedia/go-don/issues/111)) ([8f50031](https://github.com/abemedia/go-don/commit/8f50031717f53348d31619b96411dcbf60e1e6fc))
* remove Empty type (use any instead) ([#135](https://github.com/abemedia/go-don/issues/135)) ([72848e8](https://github.com/abemedia/go-don/commit/72848e8389c67f4443a1f99fc1e4a8610c831b65))

## [0.1.4](https://github.com/abemedia/go-don/compare/v0.1.3...v0.1.4) (2023-05-09)


### Bug Fixes

* error parsing GET request ([#108](https://github.com/abemedia/go-don/issues/108)) ([8298a2c](https://github.com/abemedia/go-don/commit/8298a2c7a8d46858420fccbbe39909db71838b38))

## [0.1.3](https://github.com/abemedia/go-don/compare/v0.1.2...v0.1.3) (2023-05-08)


### Features

* api serve ([#98](https://github.com/abemedia/go-don/issues/98)) ([602b24c](https://github.com/abemedia/go-don/commit/602b24c5220bee9955d30ec38e7fbc8b41aa2e10))
* implement errors interfaces & improve tests ([#91](https://github.com/abemedia/go-don/issues/91)) ([0a282f0](https://github.com/abemedia/go-don/commit/0a282f0fc2fbe289a89fd9cc0ba94939108fb205))


### Bug Fixes

* decoder panics on tag not found ([#100](https://github.com/abemedia/go-don/issues/100)) ([3a73c35](https://github.com/abemedia/go-don/commit/3a73c35dd996e1035360733d4b60d52b88c3243b))


### Performance Improvements

* improve encoding performance ([#104](https://github.com/abemedia/go-don/issues/104)) ([9dbebfa](https://github.com/abemedia/go-don/commit/9dbebfa81db3277efd964d6d8fa9f1755ef9683a))

## [0.1.2](https://github.com/abemedia/go-don/compare/v0.1.1...v0.1.2) (2023-05-05)


### Features

* improve error handling ([#85](https://github.com/abemedia/go-don/issues/85)) ([3f976fc](https://github.com/abemedia/go-don/commit/3f976fca67e518b9c786c4af32c46586fd5cdc06))


### Bug Fixes

* panic on interface request ([#82](https://github.com/abemedia/go-don/issues/82)) ([8e83bd6](https://github.com/abemedia/go-don/commit/8e83bd692db5569b36426b112d4d243cc106968a))

## [0.1.1](https://github.com/abemedia/go-don/compare/v0.1.0...v0.1.1) (2023-01-14)


### Features

* better error handling, minor refactor ([#58](https://github.com/abemedia/go-don/issues/58)) ([0de3fc3](https://github.com/abemedia/go-don/commit/0de3fc32deb4692a7e768f1f650122b664785810))
* **encoding/text:** support marshaler & stringer, improve performance ([#60](https://github.com/abemedia/go-don/issues/60)) ([b7bffe8](https://github.com/abemedia/go-don/commit/b7bffe81d2ca0651a78d694462e6684df211f0ca))

## 0.1.0 (2022-08-29)


### Bug Fixes

* **encoding/text:** int conversion on 32bit ([#32](https://github.com/abemedia/go-don/issues/32)) ([3e469fe](https://github.com/abemedia/go-don/commit/3e469fe24189849d25e24395500eca23d6043a96))
* use error's marshaler, circular ref in Handler, lint issues ([#28](https://github.com/abemedia/go-don/issues/28)) ([f8b32ea](https://github.com/abemedia/go-don/commit/f8b32eaa0150d96a6ce186f2bdf41ef0e90a39e0))
