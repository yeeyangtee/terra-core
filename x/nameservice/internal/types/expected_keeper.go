package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// MarketKeeper expected market keeper
type MarketKeeper interface {
	ComputeInternalSwap(ctx sdk.Context, offerCoin sdk.DecCoin, askDenom string) (sdk.DecCoin, sdk.Error)
}

// SupplyKeeper defines expected supply keeper
type SupplyKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) supplyexported.ModuleAccountI
	GetSupply(ctx sdk.Context) (supply supplyexported.SupplyI)
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error

	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error
}

// TreasuryKeeper defines expected treasury keeper
type TreasuryKeeper interface {
	GetTaxRate(ctx sdk.Context) (taxRate sdk.Dec)
	GetTaxCap(ctx sdk.Context, denom string) (taxCap sdk.Int)
}

// AccountKeeper defines expected account keeper
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
}
