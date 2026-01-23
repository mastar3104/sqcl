module github.com/mastar3104/sqcl

go 1.25.0

require (
	github.com/chzyer/readline v1.5.1
	github.com/go-sql-driver/mysql v1.9.3
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	golang.org/x/sys v0.0.0-20220310020820-b874c991c1a5 // indirect
)

replace github.com/chzyer/readline => ./vendor/github.com/chzyer/readline
