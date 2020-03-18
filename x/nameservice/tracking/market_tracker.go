package tracking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market"
	ns "github.com/terra-project/core/x/nameservice"
)

// MarketHook - track sending amount
func MarketHook(k ns.Keeper) core.HookHandler {
	return func(ctx sdk.Context, msg sdk.Msg, res sdk.Result) {
		if swapMsg, ok := msg.(market.MsgSwap); ok {

			if swapMsg.OfferCoin.IsZero(){
				return
			}

			nameHash, registry, err := getRegistry(ctx, k, swapMsg.Trader)
			if err != nil {
				return
			}

			registry.LockedValue = registry.LockedValue.Sub(sdk.NewCoins(swapMsg.OfferCoin))

			event, err := findEvent(res.Events, market.EventSwap)
			if err != nil {
				return
			}

			attr, err := findAttribute(event, market.AttributeKeySwapCoin)
			if err != nil {
				return
			}

			swapCoin, err := sdk.ParseCoin(string(attr))
			if err != nil {
				return
			}

			attr, err = findAttribute(event, market.AttributeKeySwapFee)
			if err != nil {
				return
			}

			feeDecCoin, err := sdk.ParseDecCoin(string(attr))
			if err != nil {
				return
			}

			if !swapCoin.IsZero() {
				registry.LockedValue = registry.LockedValue.Add(sdk.Coins{swapCoin})
			}

			if !feeDecCoin.IsZero() {
				registry.SwapFeeContribution = registry.SwapFeeContribution.Add(sdk.DecCoins{feeDecCoin})
			}

			k.SetRegistry(ctx, nameHash, registry)
		}
	}
}
