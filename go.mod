module github.com/omniboost/go-ikentoo

go 1.23.0

toolchain go1.24.4

require (
	github.com/gorilla/schema v0.0.0-20171211162101-9fa3b6af65dc
	github.com/hashicorp/go-multierror v1.1.1
	github.com/pkg/errors v0.9.1
	golang.org/x/oauth2 v0.30.0
	gopkg.in/guregu/null.v3 v3.5.0
)

require github.com/hashicorp/errwrap v1.1.0 // indirect

replace github.com/gorilla/schema => github.com/omniboost/schema v1.1.1-0.20211111150515-2e872025e306
