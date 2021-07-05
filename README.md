taggit
======

[![Go Report Card](https://goreportcard.com/badge/oss.indeed.com/go/taggit)](https://goreportcard.com/report/oss.indeed.com/go/taggit)
[![Build Status](https://travis-ci.org/modprox/taggit.svg?branch=master)](https://travis-ci.org/modprox/taggit) 
[![GoDoc](https://godoc.org/oss.indeed.com/go/taggit?status.svg)](https://godoc.org/oss.indeed.com/go/taggit)
[![NetflixOSS Lifecycle](https://img.shields.io/osslifecycle/modprox/taggit.svg)](OSSMETADATA)
[![GitHub](https://img.shields.io/github/license/modprox/taggit.svg)](LICENSE)

# Project Overview

Command `taggit` is a git repository workflow tool for publishing semver releases

# Getting Started

The `oss.indeed.com/go/taggit` CLI tool can be installed by running:

```bash
$ go install oss.indeed.com/go/taggit@latest
```

Example usage:

#### help

Use `help` to have `taggit` display proper usage arguments.

```bash
$ taggit help
Usage:  <flags> <subcommand> <subcommand args>

Subcommands:
	flags            describe all known top-level flags
	help             describe subcommands and their syntax
	list             List tagged versions.
	major            Create an incremented major version
	minor            Create an incremented minor version
	patch            Create an incremented patch version
	zero             Create initial v0.0.0 tag


Use " flags" for a list of top-level flags
```

#### list

Use `list` to have `taggit` list all of the semver compatible tags associated
with the repository.

```bash
$ taggit list
v0.0.0 |= v0.0.0
v0.0.1 |= v0.0.1
v0.3.3 |= v0.3.3
v0.4.0 |= v0.4.0 v0.4.0-rc2 v0.4.0-rc1
v0.5.0 |= v0.5.0
v0.9.0 |= v0.9.0
v1.0.0 |= v1.0.0-rc.2 v1.0.0-rc.1
```
The output is organized into two sections. On the left are the base versions of
the set of tags. On the right side is a list of tags associated with the base tag version.
The order of the list is decreasing based on semver semantics.

#### zero

Use `zero` to create the very first `v0.0.0` tag associated with a repository.
The version increment commands will not work until the very first tag already exists.

```bash
$ taggit zero
taggit: created tag: v0.0.0
```

#### patch

Use `patch` to increment the patch level of the current latest version.

```bash
$ taggit list
v0.0.0 ✓ | v0.0.0
$ taggit patch
taggit: created tag: v0.0.1
```

To publish a pre-release version, provide the pre-release version as an argument
to the `patch` command.

```bash
$ taggit patch beta1
taggit: created tag: v0.0.3-beta1
```

If the latest tag happens to be a pre-release version, the patch level is not
incremented. Instead, a tag with no pre-release extension is published at that
tag patch level.

```bash
$ taggit list
v0.0.0 ✓ | v0.0.0
v0.0.1 ✓ | v0.0.1
v0.0.2   | v0.0.2-alpha1
$ taggit patch
taggit: created tag: v0.0.2
```

#### minor

Use `minor` to increment the minor level of the current latest version.

```bash
$ taggit list
v0.0.0 ✓ | v0.0.0
v0.0.1 ✓ | v0.0.1
v0.0.2 ✓ | v0.0.2 v0.0.2-alpha1
v0.0.3 ✓ | v0.0.3 v0.0.3-beta1 
$ taggit minor
taggit: created tag: v0.1.0
```

To publish a pre-release version, provide the pre-release version as an argument
to the `minor` command.

```bash
$ taggit minor beta1
taggit: created tag: v0.4.0-beta1
```

If the latest tag happens to be a pre-release version, the minor version is not
incremented. Instead, a tag with no pre-release extension is published at that
tag minor level.

```bash
$ taggit list
v0.0.0 ✓ | v0.0.0
v0.1.0 ✓ | v0.1.0
v0.2.0 ✓ | v0.2.0 v0.2.0-beta1
v0.3.0   | v0.3.0-rc-2 v0.3.0-rc-1
$ taggit minor
taggit: created tag: v0.3.0
```

#### major

Use `major` to increment the major level of the current latest version.

```bash
$ taggit major
taggit: created tag: v1.0.0
```

To publish a pre-release version, provide the pre-release version as an argument
to the `major` command.

```bash
$ taggit major rc1
taggit: created tag: v2.0.0-rc1
```

If the latest tag happens to be a pre-release version, the major version is not
incremented. Instead, a tag with no pre-release extension is published at that tag
major level.

```bash
$ taggit list
v1.0.0 ✓ | v1.0.0
v2.0.0   | v2.0.0-rc1 
$ taggit major
taggit: created tag: v2.0.0
```
# Asking Questions

For technical questions about `taggit`, just file an issue in the GitHub tracker.

For questions about Open Source in Indeed Engineering, send us an email at
opensource@indeed.com

# Contributing

We welcome contributions! Feel free to help make `taggit` better.

### Process

- Open an issue and describe the desired feature / bug fix before making
changes. It's useful to get a second pair of eyes before investing development
effort.
- Make the change. If adding a new feature, remember to provide tests that
demonstrate the new feature works, including any error paths. If contributing
a bug fix, add tests that demonstrate the erroneous behavior is fixed.
- Open a pull request. Automated CI tests will run. If the tests fail, please
make changes to fix the behavior, and repeat until the tests pass.
- Once everything looks good, one of the indeedeng members will review the
PR and provide feedback.

# Maintainers

The `oss.indeed.com/go/taggit` module is maintained by Indeed Engineering.

While we are always busy helping people get jobs, we will try to respond to
GitHub issues, pull requests, and questions within a couple of business days.

# Code of Conduct

`oss.indeed.com/go/taggit` is governed by the [Contributer Covenant v1.4.1](CODE_OF_CONDUCT.md)

For more information please contact opensource@indeed.com.

# License

The `oss.indeed.com/go/taggit` module is open source under the [BSD-3-Clause](LICENSE) license.
