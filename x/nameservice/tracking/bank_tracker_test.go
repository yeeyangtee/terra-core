package tracking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	core "github.com/terra-project/core/types"
	"testing"
	"time"

	"github.com/terra-project/core/x/bank"
	ns "github.com/terra-project/core/x/nameservice"
	"github.com/terra-project/core/x/nameservice/internal/keeper"
)

type mockTreasuryKeeper struct{}

var constTaxRate = sdk.NewDecWithPrec(675, 5)
var constTaxCap = sdk.OneInt().MulRaw(core.MicroUnit)

// GetTaxRate mock func
func (k mockTreasuryKeeper) GetTaxRate(_ sdk.Context) (taxRate sdk.Dec) {
	return constTaxRate
}

// GetTaxCap mock func
func (k mockTreasuryKeeper) GetTaxCap(_ sdk.Context, _ string) (taxCap sdk.Int) {
	return constTaxCap
}

func TestBankHook(t *testing.T) {
	input := keeper.CreateTestInput(t)

	name := ns.Name("valid.terra")
	nameHash, childHash := name.NameHash()
	registry := ns.NewRegistry(name, keeper.Addrs[0], time.Now().UTC())
	registry.LockedValue = keeper.InitCoins
	input.NameserviceKeeper.SetRegistry(input.Ctx, nameHash, registry)
	input.NameserviceKeeper.SetResolve(input.Ctx, nameHash, childHash, keeper.Addrs[0])
	input.NameserviceKeeper.SetReverseResolve(input.Ctx, keeper.Addrs[0], nameHash)

	expectedTaxes := computeTax(input.Ctx, mockTreasuryKeeper{}, keeper.InitCoins)

	sendMsg := bank.MsgSend{FromAddress: keeper.Addrs[0], ToAddress: keeper.Addrs[1], Amount: keeper.InitCoins}
	multiSendMsg := bank.MsgMultiSend{Inputs: []bank.Input{{Address: keeper.Addrs[1], Coins: keeper.InitCoins}}, Outputs: []bank.Output{{Address: keeper.Addrs[0], Coins: keeper.InitCoins}}}

	hook := BankHook(input.NameserviceKeeper, mockTreasuryKeeper{})
	hook(input.Ctx, sendMsg, sdk.Result{})

	// locked value should be zero
	registry, err := input.NameserviceKeeper.GetRegistry(input.Ctx, nameHash)
	require.NoError(t, err)
	require.True(t, registry.LockedValue.IsZero())
	require.Equal(t, expectedTaxes, registry.TaxContribution)

	hook(input.Ctx, multiSendMsg, sdk.Result{})

	// locked value should be InitCoins
	registry, err = input.NameserviceKeeper.GetRegistry(input.Ctx, nameHash)
	require.NoError(t, err)
	require.True(t, keeper.InitCoins.IsEqual(registry.LockedValue))
	require.Equal(t, expectedTaxes, registry.TaxContribution)
}

func TestComputeTax(t *testing.T) {
	input := keeper.CreateTestInput(t)

	// luna will be ignored & krw is under cap & usd over cap & sdr is too small
	coins := sdk.Coins{
		sdk.NewInt64Coin(core.MicroLunaDenom, 12343353),
		sdk.NewInt64Coin(core.MicroKRWDenom, 1000),
		sdk.NewInt64Coin(core.MicroUSDDenom, 10000000000000000),
		sdk.NewInt64Coin(core.MicroSDRDenom, 1),
	}

	expectedTaxes := sdk.Coins{
		sdk.NewCoin(core.MicroKRWDenom, constTaxRate.MulInt64(1000).TruncateInt()),
		sdk.NewCoin(core.MicroUSDDenom, constTaxCap),
	}

	taxes := computeTax(input.Ctx, mockTreasuryKeeper{}, coins)
	require.Equal(t, expectedTaxes, taxes)
}
