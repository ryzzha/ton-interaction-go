package app

import (
	"ton-utils-go/internal/structures"

	"github.com/xssnick/tonutils-go/tlb"
)

func InitTlb() {
	tlb.Register(structures.DedustAssetNative{})
	tlb.Register(structures.DedustAssetJetton{})
}