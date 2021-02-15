package wasm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/wasm/internal/keeper"
)

// BeginBlocker handles softfork over param changes
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Uncomment to simulate state corruption at height 10
	//
	// if ctx.BlockHeader().Height == 10 {
	// 	fmt.Println("corrupting state")
	// 	params := k.GetParams(ctx)
	// 	params.MaxContractMsgSize = 4097
	// 	k.SetParams(ctx, params)
	// }
	if core.IsSoftforkHeight(ctx, 1) {
		params := k.GetParams(ctx)
		params.MaxContractMsgSize = 4096
		k.SetParams(ctx, params)
	}
}
