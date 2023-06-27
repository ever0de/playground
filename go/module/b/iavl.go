package b

import (
	"github.com/cosmos/iavl"
	dbm "github.com/tendermint/tm-db"
)

func NewMutableTree(db dbm.DB, cacheSize int, skipFastStorageUpgrade bool) (*iavl.MutableTree, error) {
	return iavl.NewMutableTree(db, cacheSize, skipFastStorageUpgrade)
}
