package main

import (
	a "test_a"
	b "test_b"

	cosmosdb "github.com/cosmos/cosmos-db"
	tmdb "github.com/tendermint/tm-db"
)

func main() {
	cosmosDB := cosmosdb.NewMemDB()
	a.NewMutableTree(cosmosDB, 0, false)

	tmDB := tmdb.NewMemDB()
	b.NewMutableTree(tmDB, 0, false)
}
