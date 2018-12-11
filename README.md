## taggit

`taggit` is a git repository workflow tool for publishing semver releases

documentation @ [modprox.org](https://modprox.org)

[![Go Report Card](https://goreportcard.com/badge/github.com/modprox/taggit)](https://goreportcard.com/report/github.com/modprox/taggit) 
[![Build Status](https://travis-ci.org/modprox/taggit.svg?branch=master)](https://travis-ci.org/modprox/taggit) 
[![GoDoc](https://godoc.org/github.com/modprox/taggit?status.svg)](https://godoc.org/github.com/modprox/taggit) 
[![License](https://img.shields.io/github/license/modprox/taggit.svg?style=flat-square)](LICENSE)

### Project Management
Issue [tracker](https://github.com/modprox/taggit/issues)

### Instruction Manual

#### help
Use `help` to have `taggit` display proper usage arguments.
```
$ taggit help
usage: taggit [list, zero, patch, minor, major]
```

#### list
Use `list` to have `taggit` list all of the semver compatible tags associated with the repository.
```
$ taggit list
v0.0.1 ✓ | v0.0.1 v0.0.1-alpha2 v0.0.1-alpha1 v0.0.1-alpha
v0.0.2   | v0.0.2-rc2 v0.0.2-rc1
v0.0.3 ✓ | v0.0.3 v0.0.3-rc3
v0.0.4 ✓ | v0.0.4 v0.0.4-alpha2 v0.0.4-alpha1
v0.0.5 ✓ | v0.0.5
v0.0.6 ✓ | v0.0.6
v0.0.7   | v0.0.7-rc2 v0.0.7-rc1
```
The output is organized into two sections. On the left are the base versions of the set of tags. On the right side is a list of tags associated with the base tag version. The order of the list is decreasing based on semver semantics. If the base version is a published tag, a checkmark appears next to it on the left side.

#### zero
Use `zero` to create the very first `v0.0.0` tag associated with a repository. The version increment commands will not work until the very first tag already exists.
```
$ taggit zero
taggit: creating tag: v0.0.0
taggit: created tag: v0.0.0
taggit: published tag: v0.0.0
```

#### patch
Use `patch` to increment the patch level of the current latest version.
```
$ taggit list
v0.0.0 ✓ | v0.0.0
$ taggit patch
taggit: creating tag: v0.0.1
taggit: created tag: v0.0.1
taggit: published tag: v0.0.1
```
If the latest tag happens to be a pre-release version, the patch level is not incremented. Instead, a tag with no pre-release extension is published at that tag patch level.
```
$ taggit list
v0.0.0 ✓ | v0.0.0
v0.0.1 ✓ | v0.0.1
v0.0.2   | v0.0.2-alpha1
$ taggit patch
taggit: creating tag: v0.0.2
taggit: created tag: v0.0.2
taggit: published tag: v0.0.2
```
To publish a pre-release version, provide the pre-release version as an argument to the `patch` command.
```
$ taggit patch beta1
taggit: creating tag: v0.0.3-beta1
taggit: created tag: v0.0.3-beta1
taggit: published tag: v0.0.3-beta1
```

#### minor
Use `minor` to increment the minor level of the current latest version.
```
$ taggit list
v0.0.0 ✓ | v0.0.0
v0.0.1 ✓ | v0.0.1
v0.0.2 ✓ | v0.0.2 v0.0.2-alpha1
v0.0.3 ✓ | v0.0.3 v0.0.3-beta1 
$ taggit minor
taggit: creating tag: v0.1.0
taggit: created tag: v0.1.0
taggit: published tag: v0.1.0
```
If the latest tag happens to be a pre-release version, the minor version is not incremented. Instead, a tag with no pre-release extension is published at that tag minor level.
```
[p1 bar (master)] $ taggit list
v0.0.0 ✓ | v0.0.0
v0.1.0 ✓ | v0.1.0
v0.2.0 ✓ | v0.2.0 v0.2.0-beta1
v0.3.0   | v0.3.0-rc-2 v0.3.0-rc-1
[p1 bar (master)] $ taggit minor
taggit: creating tag: v0.3.0
taggit: created tag: v0.3.0
taggit: published tag: v0.3.0
```
To publish a pre-release version, provide the pre-release version as an argument to the `minor` command.
```
$ taggit minor beta1
taggit: creating tag: v0.4.0-beta1
taggit: created tag: v0.4.0-beta1
taggit: published tag: v0.4.0-beta1
```

#### major
Use `major` to increment the major level of the current latest version.
```
$ taggit major
taggit: creating tag: v1.0.0
taggit: created tag: v1.0.0
taggit: published tag: v1.0.0
```
If the latest tag happens to be a pre-release version, the major version is not incremented. Instead, a tag with no pre-release extension is published at that tag major level.
```
$ taggit list
v1.0.0 ✓ | v1.0.0
v2.0.0   | v2.0.0-rc1 
$ taggit major
taggit: creating tag: v2.0.0
taggit: created tag: v2.0.0
taggit: published tag: v2.0.0
```
To publish a pre-release version, provide the pre-release version as an argument to the `major` command.
```
$ taggit major rc1
taggit: creating tag: v2.0.0-rc1
taggit: created tag: v2.0.0-rc1
taggit: published tag: v2.0.0-rc1
```
