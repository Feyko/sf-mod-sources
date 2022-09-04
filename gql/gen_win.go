//go:build gen
// +build gen

package gql

//go:generate powershell "go run git.sr.ht/~emersion/gqlclient/cmd/gqlintrospect https://api.ficsit.app/v2/query > schema.graphqls ; (Get-Content schema.graphqls) -replace '`n', '`r`n' | Set-Content schema.graphqls"
//go:generate go run git.sr.ht/~emersion/gqlclient/cmd/gqlclientgen -s schema.graphqls -o generated.go -q query.graphql
