#!/usr/bin/env bash

set -euo pipefail

if [[ "${1} ${2} ${3} ${4}" != "list -u -m all" ]]; then
    echo 'expected "list -u -m all"'
    exit 1
fi

cat <<'EOF'
github.com/modprox/taggit
github.com/boltdb/bolt v1.3.1
github.com/cactus/go-statsd-client v3.1.1+incompatible
github.com/davecgh/go-spew v1.1.1
github.com/go-sql-driver/mysql v1.4.0 [v1.4.1]
github.com/gorilla/context v1.1.1
github.com/gorilla/csrf v1.5.1
github.com/gorilla/mux v1.6.2
github.com/gorilla/securecookie v1.1.1
github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
github.com/lib/pq v1.0.0
github.com/modprox/mp v0.0.5 [v0.0.10]
github.com/pkg/errors v0.8.0
github.com/pmezard/go-difflib v1.0.0
github.com/shoenig/atomicfs v0.1.1
github.com/shoenig/petrify/v4 v4.0.2
github.com/shoenig/regexplus v0.0.0
github.com/shoenig/toolkit v1.0.0
github.com/stretchr/objx v0.1.1
github.com/stretchr/testify v1.2.2
golang.org/x/sys v0.0.0-20180909124046-d0be0721c37e [v0.0.0-20181217223516-dcdaa6325bcb]
google.golang.org/appengine v1.1.0 [v1.3.0]
EOF

exit 0