package a

import (
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/iavl"
)

func NewMutableTree(db dbm.DB, cacheSize int, skipFastStorageUpgrade bool) (*iavl.MutableTree, error) {
	return iavl.NewMutableTree(db, cacheSize, skipFastStorageUpgrade)
}
