# lukasmalkmus/expensify-go

> An Expensify API client. - by **[Lukas Malkmus]**

[![Go Reference][gopkg_badge]][gopkg]
[![Go Workflow][go_workflow_badge]][go_workflow]
[![Coverage Status][coverage_badge]][coverage]
[![Go Report][report_badge]][report]
[![Latest Release][release_badge]][release]
[![License][license_badge]][license]

---

## Table of Contents

1. [Introduction](#introduction)
1. [Usage](#usage)
1. [Contributing](#contributing)
1. [License](#license)

## Introduction

_expensify-go_ is an opinionated client library for the Expensify API. I created
it in order to add expenses which makes it the only method currently supported.

## Usage

### Installation

```bash
go get github.com/lukasmalkmus/expensify-go
```

### Usage

```go
// Get credentials from https://www.expensify.com/tools/integrations.
client, err := expensify.NewClient("XXX-REPLACE-ME-XXX", "XXX-REPLACE-ME-XXX")
if err != nil {
    // Handle error!
}

expense := &expensify.Expense{
    Merchant: "Apple Inc.",
    Created:  expensify.NewTime(time.Now()),
    Amount:   99,
    Currency: "USD",
}

res, err := client.Expense.Create(context.TODO(), "you@example.com", []*expensify.Expense{exp})
if err != nil {
    // Handle error!
}

fmt.Println(res[0].TransactionID)
```

## Contributing

Feel free to submit PRs or to fill Issues. Every kind of help is appreciated.

## License

© Lukas Malkmus, 2020

Distributed under MIT License (`The MIT License`).

See [LICENSE](LICENSE) for more information.

[![License Status][license_status_badge]][license_status]

<!-- Links -->

[Lukas Malkmus]: https://github.com/lukasmalkmus

<!-- Badges -->

[gopkg]: https://pkg.go.dev/github.com/lukasmalkmus/expensify-go
[gopkg_badge]: https://img.shields.io/badge/doc-reference-007d9c?logo=go&logoColor=white&style=flat-square
[go_workflow]: https://github.com/lukasmalkmus/expensify-go/actions?query=workflow%3Ago
[go_workflow_badge]: https://img.shields.io/github/workflow/status/lukasmalkmus/expensify-go/go?style=flat-square&ghcache=unused
[coverage]: https://codecov.io/gh/lukasmalkmus/expensify-go
[coverage_badge]: https://img.shields.io/codecov/c/github/lukasmalkmus/expensify-go.svg?style=flat-square&ghcache=unused
[report]: https://goreportcard.com/report/github.com/lukasmalkmus/expensify-go
[report_badge]: https://goreportcard.com/badge/github.com/lukasmalkmus/expensify-go?style=flat-square&ghcache=unused
[release]: https://github.com/lukasmalkmus/expensify-go/releases/latest
[release_badge]: https://img.shields.io/github/release/lukasmalkmus/expensify-go.svg?style=flat-square&ghcache=unused
[license]: https://opensource.org/licenses/MIT
[license_badge]: https://img.shields.io/github/license/lukasmalkmus/expensify-go.svg?color=blue&style=flat-square&ghcache=unused
[license_status]: https://app.fossa.com/projects/git%2Bgithub.com%2Flukasmalkmus%expensify-go
[license_status_badge]: https://app.fossa.com/api/projects/git%2Bgithub.com%2Flukasmalkmus%expensify-go.svg?type=large&ghcache=unused
