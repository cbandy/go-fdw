[![Build Status](https://travis-ci.org/cbandy/go-fdw.svg?branch=master)](https://travis-ci.org/cbandy/go-fdw)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fcbandy%2Fgo-fdw.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fcbandy%2Fgo-fdw?ref=badge_shield)

# Go-FDW

Go-FDW is a package for implementing [PostgreSQL][] [foreign data wrappers][] in [Go][] and [cgo][].

The goal is to let you write your wrapper almost entirely in Go with this package as an import.

The API is young and very likely to change; please open [issues][] with any ideas, feature requests, or bugs you encounter!

[cgo]: https://golang.org/cmd/cgo/
[foreign data wrappers]: https://www.postgresql.org/docs/current/static/ddl-foreign-data.html
[issues]: https://github.com/cbandy/go-fdw/issues
[Go]: https://golang.org/
[PostgreSQL]: https://www.postgresql.org/


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fcbandy%2Fgo-fdw.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fcbandy%2Fgo-fdw?ref=badge_large)