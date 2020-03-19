package tracking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	core "github.com/terra-project/core/types"
	ns "github.com/terra-project/core/x/nameservice"
	"github.com/terra-project/core/x/nameservice/internal/types"
)

// BankHook - track sending amount
func BankHook(k ns.Keeper, tk types.TreasuryKeeper) core.HookHandler {
	return func(ctx sdk.Context, msg sdk.Msg, _ sdk.Result) {
		if sendMsg, ok := msg.(bank.MsgSend); ok {
			updateLockedValueWithTaxContribution(ctx, k, tk, sendMsg.FromAddress, sendMsg.Amount, true)
			updateLockedValueWithTaxContribution(ctx, k, tk, sendMsg.ToAddress, sendMsg.Amount, false)
		} else if multiMsg, ok := msg.(bank.MsgMultiSend); ok {

			for _, input := range multiMsg.Inputs {
				updateLockedValueWithTaxContribution(ctx, k, tk, input.Address, input.Coins, true)
			}

			for _, output := range multiMsg.Outputs {
				updateLockedValueWithTaxContribution(ctx, k, tk, output.Address, output.Coins, false)
			}
		}
	}
}

func updateLockedValueWithTaxContribution(ctx sdk.Context, k ns.Keeper, tk types.TreasuryKeeper, accAddr sdk.AccAddress, amount sdk.Coins, isSender bool) {
	if amount.IsZero() {
		return
	}

	nameHash, registry, err := getRegistry(ctx, k, accAddr)
	if err != nil {
		return
	}

	if isSender {
		taxes := computeTax(ctx, tk, amount)

		registry.LockedValue = registry.LockedValue.Sub(amount)
		registry.TaxContribution = registry.TaxContribution.Add(taxes)
	} else {
		registry.LockedValue = registry.LockedValue.Add(amount)
	}

	k.SetRegistry(ctx, nameHash, registry)
}

// computes the stability tax according to tax-rate and tax-cap
func computeTax(ctx sdk.Context, tk types.TreasuryKeeper, principal sdk.Coins) (taxes sdk.Coins) {
	taxRate := tk.GetTaxRate(ctx)
	if taxRate.Equal(sdk.ZeroDec()) {
		return
	}

	for _, coin := range principal {
		if coin.Denom == core.MicroLunaDenom {
			continue
		}

		taxDue := sdk.NewDecFromInt(coin.Amount).Mul(taxRate).TruncateInt()

		// If tax due is greater than the tax cap, cap!
		taxCap := tk.GetTaxCap(ctx, coin.Denom)
		if taxDue.GT(taxCap) {
			taxDue = taxCap
		}

		if taxDue.IsZero() {
			continue
		}

		taxes = taxes.Add(sdk.NewCoins(sdk.NewCoin(coin.Denom, taxDue)))
	}

	return
}
