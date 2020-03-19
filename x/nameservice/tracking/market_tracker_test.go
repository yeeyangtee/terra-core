package tracking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/market"
	ns "github.com/terra-project/core/x/nameservice"
	"github.com/terra-project/core/x/nameservice/internal/keeper"
	"testing"
	"time"
)

func TestMarketHook(t *testing.T) {
	input := keeper.CreateTestInput(t)

	name := ns.Name("valid.terra")
	nameHash, childHash := name.NameHash()
	registry := ns.NewRegistry(name, keeper.Addrs[0], time.Now().UTC())
	registry.LockedValue = keeper.InitCoins
	input.NameserviceKeeper.SetRegistry(input.Ctx, nameHash, registry)
	input.NameserviceKeeper.SetResolve(input.Ctx, nameHash, childHash, keeper.Addrs[0])
	input.NameserviceKeeper.SetReverseResolve(input.Ctx, keeper.Addrs[0], nameHash)

	offerCoin := sdk.NewInt64Coin(core.MicroLunaDenom, core.MicroUnit)
	retCoin := sdk.NewInt64Coin(core.MicroKRWDenom, 1000*core.MicroUnit)
	swapMsg := market.MsgSwap{Trader: keeper.Addrs[0], OfferCoin: offerCoin, AskDenom: core.MicroKRWDenom}
	expectedSwapFee := sdk.NewInt64DecCoin(core.MicroKRWDenom, 1000)

	hook := MarketHook(input.NameserviceKeeper)
	input.Ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			market.EventSwap,
			sdk.NewAttribute(market.AttributeKeySwapCoin, retCoin.String()),
			sdk.NewAttribute(market.AttributeKeySwapFee, expectedSwapFee.String()),
		),
	})

	hook(input.Ctx, swapMsg, sdk.Result{Events: input.Ctx.EventManager().Events()})
	registry, err := input.NameserviceKeeper.GetRegistry(input.Ctx, nameHash)
	require.NoError(t, err)
	require.True(t, registry.SwapFeeContribution.IsEqual(sdk.DecCoins{expectedSwapFee}))
	require.True(t, registry.LockedValue.IsEqual(keeper.InitCoins.Sub(sdk.Coins{offerCoin}).Add(sdk.Coins{retCoin})))

	// state must be same without events
	hook(input.Ctx, swapMsg, sdk.Result{})
	registry, err = input.NameserviceKeeper.GetRegistry(input.Ctx, nameHash)
	require.NoError(t, err)
	require.True(t, registry.SwapFeeContribution.IsEqual(sdk.DecCoins{expectedSwapFee}))
	require.True(t, registry.LockedValue.IsEqual(keeper.InitCoins.Sub(sdk.Coins{offerCoin}).Add(sdk.Coins{retCoin})))

	// state must be same without attributes
	input.Ctx = input.Ctx.WithEventManager(sdk.NewEventManager())
	input.Ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			market.EventSwap,
		),
	})

	hook(input.Ctx, swapMsg, sdk.Result{Events: input.Ctx.EventManager().Events()})
	registry, err = input.NameserviceKeeper.GetRegistry(input.Ctx, nameHash)
	require.NoError(t, err)
	require.True(t, registry.SwapFeeContribution.IsEqual(sdk.DecCoins{expectedSwapFee}))
	require.True(t, registry.LockedValue.IsEqual(keeper.InitCoins.Sub(sdk.Coins{offerCoin}).Add(sdk.Coins{retCoin})))

	// state must be same without any single attribute
	input.Ctx = input.Ctx.WithEventManager(sdk.NewEventManager())
	input.Ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			market.EventSwap,
			sdk.NewAttribute(market.AttributeKeySwapCoin, retCoin.String()),
		),
	})

	hook(input.Ctx, swapMsg, sdk.Result{Events: input.Ctx.EventManager().Events()})
	registry, err = input.NameserviceKeeper.GetRegistry(input.Ctx, nameHash)
	require.NoError(t, err)
	require.True(t, registry.SwapFeeContribution.IsEqual(sdk.DecCoins{expectedSwapFee}))
	require.True(t, registry.LockedValue.IsEqual(keeper.InitCoins.Sub(sdk.Coins{offerCoin}).Add(sdk.Coins{retCoin})))

	// state must be same with invalid attributes
	input.Ctx = input.Ctx.WithEventManager(sdk.NewEventManager())
	input.Ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			market.EventSwap,
			sdk.NewAttribute(market.AttributeKeySwapCoin, "invalid coin"),
		),
	})

	hook(input.Ctx, swapMsg, sdk.Result{Events: input.Ctx.EventManager().Events()})
	registry, err = input.NameserviceKeeper.GetRegistry(input.Ctx, nameHash)
	require.NoError(t, err)
	require.True(t, registry.SwapFeeContribution.IsEqual(sdk.DecCoins{expectedSwapFee}))
	require.True(t, registry.LockedValue.IsEqual(keeper.InitCoins.Sub(sdk.Coins{offerCoin}).Add(sdk.Coins{retCoin})))


	// state must be same with invalid attributes 2
	input.Ctx = input.Ctx.WithEventManager(sdk.NewEventManager())
	input.Ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			market.EventSwap,
			sdk.NewAttribute(market.AttributeKeySwapCoin, retCoin.String()),
			sdk.NewAttribute(market.AttributeKeySwapFee, "invalid coin"),
		),
	})

	hook(input.Ctx, swapMsg, sdk.Result{Events: input.Ctx.EventManager().Events()})
	registry, err = input.NameserviceKeeper.GetRegistry(input.Ctx, nameHash)
	require.NoError(t, err)
	require.True(t, registry.SwapFeeContribution.IsEqual(sdk.DecCoins{expectedSwapFee}))
	require.True(t, registry.LockedValue.IsEqual(keeper.InitCoins.Sub(sdk.Coins{offerCoin}).Add(sdk.Coins{retCoin})))

	// swap from the account which has no registry
	swapMsg = market.MsgSwap{Trader: keeper.Addrs[1], OfferCoin: offerCoin, AskDenom: core.MicroKRWDenom}
	hook(input.Ctx, swapMsg, sdk.Result{})
}
