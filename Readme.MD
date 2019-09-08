# Mailchain

[![Build Status](https://travis-ci.com/mailchain/mailchain.svg?branch=master&style=for-the-badge)](https://travis-ci.com/mailchain/mailchain.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/mailchain/mailchain/badge.svg?branch=master)](https://coveralls.io/github/mailchain/mailchain?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mailchain/mailchain)](https://goreportcard.com/report/github.com/mailchain/community) [![Join the chat at https://gitter.im/Mailchain/mailchain](https://badges.gitter.im/Mailchain/mailchain.svg)](https://gitter.im/Mailchain/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Introduction

Mailchain enables blockchain-based email-like messaging with plain or rich text and attachment capabilities. Using blockchain protocols and decentralized storage, Mailchain delivers a simple, secure, messaging experience.

The majority of blockchain protocols provide no standard way to handle messaging. Some account holders have sent messages as an encrypted or unencrypted string included in transaction data. Some applications work around the problem by asking users to link another method of contact (e.g. email address, phone number etc.) to an application. This compromises anonymity by asking users to link or reveal an identity. It also increases exposure to security risks, and relies on additional centralized services.

Mailchain provides a solution to these problems and a generalized solution to email message on blockchain protocols.

## Repo layout

This is the primary repo for the Mailchain development.

The complete list of Mailchain repositories is:

* [mailchain](mailchain): Main Mailchain repo, API and CLI
* [mailchain-web](mailchain-web): Inbox web client, the repository for https://inbox.mailchain.xyy/
* [docs.mailchain.xyz](docs.mailchain.xyz): Documentation hosted at https://docs.mailchain.xyz
* [mailchain-specification](mailchain-specification): Original specification and design documents

## Quickstart and documentation

You can run Mailchain on most operating systems in a matter of seconds. See
the [Mailchain Getting Started Guide]([getting-started](https://docs.mailchain.xyz/getting-started)) for how.

For more comprehensive documentation, start with the [Mailchain
docs][[linkerd-docs](https://docs.mailchain.xyz/concepts/overview)]. (The doc source code is available in the
[docs](docs.mailchain.xyz) repo.)

## Working in this repo

We :heart: pull requests! See [`CONTRIBUTING.md`][CONTRIBUTING] for info on
contributing changes.

## Get involved

* Join our community discussion on [Mailchain Gitter](https://gitter.im/Mailchain/community).
* Follow [@mailchain_xyz](https://twitter.com/mailchain_xyz) on Twitter.
* Contact us via regular [email](mailto:team@mailchain.xyz).

## Code of conduct

This project is for everyone. We ask that our users and contributors take a few
minutes to review our [code of conduct][CODE_OF_CONDUCT].

## Checksums

To check the application you downloaded matches the build SHA256 checksum.

For example, on a Mac:

``` bash
shasum -a 256 path_to_file/mailchain_...
```

The output should match the checksum found in checksums.txt

## License

Copyright 2019, Finobo. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
these files except in compliance with the License. You may obtain a copy of the
License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
