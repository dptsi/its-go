# Changelog

All notable changes to this project will be documented in this file. See [commit-and-tag-version](https://github.com/absolute-version/commit-and-tag-version) for commit guidelines.

## [1.9.8](https://github.com/dptsi/its-go/compare/v1.9.7...v1.9.8) (2025-09-25)


### Features

* Group sentry transaction by path ([dc3fd21](https://github.com/dptsi/its-go/commit/dc3fd2159142a8ad51ae259a13cd80748ced58a5))

## [1.9.7](https://github.com/dptsi/its-go/compare/v1.9.6...v1.9.7) (2025-09-23)


### Features

* Use Sentry Hub from `sentrygin` to Allow Event Enrichments ([556f30d](https://github.com/dptsi/its-go/commit/556f30d8fe77751c25b5b2a7cb4c7d2e0d758d4b))

## [1.9.6](https://github.com/dptsi/its-go/compare/v1.9.5...v1.9.6) (2025-09-22)


### Features

* integrate Sentry for error tracking ([af1f13a](https://github.com/dptsi/its-go/commit/af1f13a4dfb8c5ee037e7af6a9523e38bbd126ae))

## [1.9.5](https://github.com/dptsi/its-go/compare/v1.9.4...v1.9.5) (2025-09-18)


### Bug Fixes

* fix adding reg_id to user model ([64d9ee2](https://github.com/dptsi/its-go/commit/64d9ee26afc084ac656f4c4c5694728976a7b7bb))

## [1.9.4](https://github.com/dptsi/its-go/compare/v1.9.3...v1.9.4) (2025-09-18)


### Features

* add reg_id to user model ([95c35c0](https://github.com/dptsi/its-go/commit/95c35c012b4d85645a410f597f8cbdd2c0f39b2c))

## [1.9.3](https://github.com/dptsi/its-go/compare/v1.9.2...v1.9.3) (2025-07-03)


### Bug Fixes

* cannot set custom firestore database id ([5628058](https://github.com/dptsi/its-go/commit/5628058e8cfc4723094d45fa2b646f3a96956706))

## [1.9.2](https://github.com/dptsi/its-go/compare/v1.9.1...v1.9.2) (2024-11-05)

## [1.9.1](https://github.com/dptsi/its-go/compare/v1.9.0...v1.9.1) (2024-09-19)

## [1.9.0](https://github.com/dptsi/its-go/compare/v1.8.1...v1.9.0) (2024-07-17)


### Features

* **firestore:** add firestore service ([d98ba73](https://github.com/dptsi/its-go/commit/d98ba73d4f0e343f2042506fd18146eed33b29ce))

## [1.8.1](https://github.com/dptsi/its-go/compare/v1.8.0...v1.8.1) (2024-07-01)


### Bug Fixes

* **event:** missing error feedback ([602c0b8](https://github.com/dptsi/its-go/commit/602c0b8425f7a377e1fe9d660fac0c3f0372e883))

## [1.8.0](https://github.com/dptsi/its-go/compare/v1.7.1...v1.8.0) (2024-07-01)


### Features

* **app:** add graceful shutdown feature ([86f76fa](https://github.com/dptsi/its-go/commit/86f76fa519e8ca2e15eb2a4aed519e670980264d))

## [1.7.1](https://github.com/dptsi/its-go/compare/v1.7.0...v1.7.1) (2024-06-27)


### Bug Fixes

* **event:** prevent repetitive event name ([42d9776](https://github.com/dptsi/its-go/commit/42d97763a5c97af36dad855c4e4210ee6693cf5e))

## [1.7.0](https://github.com/dptsi/its-go/compare/v1.6.1...v1.7.0) (2024-06-24)


### Features

* **app:** don't return base url in infinite scroll & table advance ([e1586f2](https://github.com/dptsi/its-go/commit/e1586f26ee5792535989bdadb2ae5c2a06423eb3))

## [1.6.1](https://github.com/dptsi/its-go/compare/v1.6.0...v1.6.1) (2024-06-24)


### Bug Fixes

* **oidc:** calling oidc endpoints taking too long ([8d43c42](https://github.com/dptsi/its-go/commit/8d43c42728dd56e961ed66af4a7168a775e93c1d))

## [1.6.0](https://github.com/dptsi/its-go/compare/v1.5.5...v1.6.0) (2024-06-12)


### Features

* **auth:** unit org ([0307d77](https://github.com/dptsi/its-go/commit/0307d777977c7cfdbff8e374f9b6b014eb4cf4c0))

## [1.5.5](https://github.com/dptsi/its-go/compare/v1.5.4...v1.5.5) (2024-06-09)


### Bug Fixes

* initiating oidc instance taking too long ([75ad146](https://github.com/dptsi/its-go/commit/75ad1463699af1b68c5d245ab684f0d720821534))

## [1.5.4](https://github.com/dptsi/its-go/compare/v1.5.3...v1.5.4) (2024-05-31)


### Bug Fixes

* **app:** query params not being preserved in infinite scroll response ([36faa16](https://github.com/dptsi/its-go/commit/36faa169fa5b517167c78f4f493a1d49b507b9e8))

## [1.5.3](https://github.com/dptsi/its-go/compare/v1.5.2...v1.5.3) (2024-04-24)


### Bug Fixes

* **app:** missing infinite scroll response by base url ([754f61c](https://github.com/dptsi/its-go/commit/754f61ceb6bde745be529d20a0ef19b412e84de1))

## [1.5.2](https://github.com/dptsi/its-go/compare/v1.5.1...v1.5.2) (2024-04-24)


### Bug Fixes

* role permissions not being mapped ([eb5777e](https://github.com/dptsi/its-go/commit/eb5777e09192526b2a44f2b39e47ce684b9ebed1))

## [1.5.1](https://github.com/dptsi/its-go/compare/v1.5.0...v1.5.1) (2024-03-16)


### Features

* **http:** use exact path match rather than regex for csrf exclusion ([71bffbf](https://github.com/dptsi/its-go/commit/71bffbffb52abb0427ada2246a0c68f3514777fb))

## [1.5.0](https://github.com/dptsi/its-go/compare/v1.4.3...v1.5.0) (2024-02-28)


### Features

* custom group to role mapping ([2a92332](https://github.com/dptsi/its-go/commit/2a92332c6f4c5e299bd7d8216c722ebda9a7a4ad))

## [1.4.3](https://github.com/dptsi/its-go/compare/v1.4.2...v1.4.3) (2024-02-28)


### Bug Fixes

* **sessions:** wrong naming ([a292d6d](https://github.com/dptsi/its-go/commit/a292d6d3db498b3b67dc7768903aa62017669de2))

## [1.4.2](https://github.com/dptsi/its-go/compare/v1.4.1...v1.4.2) (2024-02-28)


### Features

* **sessions:** firestore storage ([8d4155e](https://github.com/dptsi/its-go/commit/8d4155ee7a186ed509e816be3118c8cbcf9f6d99))


### Bug Fixes

* **session:** slow tx performance & can't delete regenerated session ([952a95d](https://github.com/dptsi/its-go/commit/952a95d861a18d072c7e4c02707c5a3e7a3796da))

## [1.4.1](https://github.com/dptsi/its-go/compare/v1.4.0...v1.4.1) (2024-02-28)


### Bug Fixes

* incorrect uuid error when sid is not valid uuid ([4f771df](https://github.com/dptsi/its-go/commit/4f771df8ebc0602febb99d20949614b47dfd30a8))

## [1.4.0](https://github.com/dptsi/its-go/compare/v1.3.2...v1.4.0) (2024-02-16)


### Features

* **activity log:** service ([40d9a83](https://github.com/dptsi/its-go/commit/40d9a8379717f315350e1f705df406aa9b5f9196))
* **auth:** impersonator_id information ([f3a0c93](https://github.com/dptsi/its-go/commit/f3a0c93411d67048097313e306bc8c28d91e9492))
* **logging:** service ([3b5c23c](https://github.com/dptsi/its-go/commit/3b5c23c6a1f1a298e1dd740b70420e1889f715a5))


### Bug Fixes

* unused logging info ([eff9b41](https://github.com/dptsi/its-go/commit/eff9b41a3524bce689ed728442e5f97468b1d146))
* **web:** missing logging ([47261dc](https://github.com/dptsi/its-go/commit/47261dcac37b062b19008d12a44de13483b73550))

## [1.3.2](https://github.com/dptsi/its-go/compare/v1.3.1...v1.3.2) (2024-02-16)


### Bug Fixes

* **sso:** missing phone and group ([353eddb](https://github.com/dptsi/its-go/commit/353eddb1c10ddb666e5ad6b806c1dbae25b73cac))

## [1.3.1](https://github.com/dptsi/its-go/compare/v1.3.0...v1.3.1) (2024-01-11)


### Bug Fixes

* session expired without checking last active time ([9f36160](https://github.com/dptsi/its-go/commit/9f361601b87d929ff51341897fdfc803a1b9364a))

## [1.3.0](https://github.com/dptsi/its-go/compare/v1.2.1...v1.3.0) (2024-01-09)


### Features

* **script:** add created file path ([fd40294](https://github.com/dptsi/its-go/commit/fd40294a386fdf3a02ee44c6eb3537cb3e3a21de))

## [1.2.1](https://github.com/dptsi/its-go/compare/v1.2.0...v1.2.1) (2024-01-09)

## [1.2.0](https://github.com/dptsi/its-go/compare/v1.1.1...v1.2.0) (2024-01-09)


### Features

* **script:** add success message ([4b194a6](https://github.com/dptsi/its-go/commit/4b194a6dc9c6bb6e045a9d9f1e1e50588737a16d))
* **script:** make command ([8baac03](https://github.com/dptsi/its-go/commit/8baac030e861eaef798b9cc1c9a7ef41afb3c521))
* **script:** make controller ([11e48bb](https://github.com/dptsi/its-go/commit/11e48bb064d95938b0b71248b125d46c195963e1))
* **script:** make entity ([55332de](https://github.com/dptsi/its-go/commit/55332deb47d5e147035a575348888e44743f32ca))
* **script:** make event ([57d9c58](https://github.com/dptsi/its-go/commit/57d9c5882d9a3894f69b04cd6f2f194f4946d8e5))
* **script:** make query object ([082efda](https://github.com/dptsi/its-go/commit/082efda8736fc4572cf064d61b1a85c627a93fdf))
* **script:** make repository ([f38249c](https://github.com/dptsi/its-go/commit/f38249cae1eaf5042e68920da14714de36d7ba9d))
* **script:** make value object ([61ac4b7](https://github.com/dptsi/its-go/commit/61ac4b74afb1588442df40855c02c6afa8786156))
* **script:** prevent user from accidentally add suffix ([564253e](https://github.com/dptsi/its-go/commit/564253ecaddcfa1319208d59f6a782f67d996139))


### Bug Fixes

* **script:** file overwritten if exists ([00765f9](https://github.com/dptsi/its-go/commit/00765f9033fe90a9baa8e7fd2873115f641b8292))

## [1.1.1](https://github.com/dptsi/its-go/compare/v1.1.0...v1.1.1) (2024-01-04)

## [1.1.0](https://github.com/dptsi/its-go/compare/v1.0.0...v1.1.0) (2024-01-04)


### Features

* remove active role from auth logic ([b8f5d20](https://github.com/dptsi/its-go/commit/b8f5d20f11bb1da37e7d837f4b3feebb7e832063))

## 1.0.0 (2024-01-04)


### Features

* add app instance 1210e21
* **app:** event, errors, and pagination 9fac6ef
* **auth:** services 3c4f9f3
* context in app and module service e1f1b3b
* **crypt:** encrypt/decrypt service 340038c
* csrf cookie route a63d140
* customizable csrf config 01e4bc1
* customizable port ba8348e
* **entra,sso:** auth provider 45b32da
* event handler fadc119
* extendable middleware f31b585
* middleware service 53f9bf3
* **oidc:** oidc library 9a1bb43
* remove module specific bind function 4117249
* **script:** base service b74a520
* **script:** extend auth, events, and middleware boilerplate 502897e
* **script:** generate app key 5666da1
* **script:** make module 8381da2


### Bug Fixes

* auth guard not found 9e7afec
* compile error in base go f442c1e
* csrf cookie not showed in swagger 703341c
* error message not descriptive 8792c06
* invalidated session not being deleted 29922ae
* middleware group bug d92dc37
* missing where clause in session efe2f67
* session lifetime not in minutes 57fbb88
* session provider not scalable 4bd5adf
* wrong method in route template d815093
