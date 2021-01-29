package keeper

import (
	"fmt"

	wasmvm "github.com/CosmWasm/wasmvm"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-project/core/x/wasm/internal/types"
)

func humanAddress(canon []byte) (string, uint64, error) {
	if len(canon) != sdk.AddrLen {
		return "", types.HumanizeCost, fmt.Errorf("Expected %d byte address", sdk.AddrLen)
	}
	return sdk.AccAddress(canon).String(), types.HumanizeCost, nil
}

func canonicalAddress(human string) ([]byte, uint64, error) {
	bz, err := sdk.AccAddressFromBech32(human)
	return bz, types.CanonicalizeCost, err
}

var cosmwasmAPI = wasmvm.GoAPI{
	HumanAddress:     humanAddress,
	CanonicalAddress: canonicalAddress,
}

// wasmGasMeter wraps the GasMeter from context and multiplies all reads by out defined multiplier
type wasmGasMeter struct {
	originalMeter sdk.GasMeter
	gasMultiplier uint64
}

var _ wasmvm.GasMeter = wasmGasMeter{}

func (m wasmGasMeter) GasConsumed() sdk.Gas {
	return m.originalMeter.GasConsumed() * m.gasMultiplier
}

// return gas meter interface for wasm gas meter
func (k Keeper) getGasMeter(ctx sdk.Context) wasmGasMeter {
	return wasmGasMeter{
		originalMeter: ctx.GasMeter(),
		gasMultiplier: types.GasMultiplier,
	}
}

// return remaining gas in wasm gas unit
func (k Keeper) getGasRemaining(ctx sdk.Context) uint64 {
	meter := ctx.GasMeter()

	// avoid integer overflow
	if meter.IsOutOfGas() {
		return 0
	}

	remaining := (meter.Limit() - meter.GasConsumed())
	if maxGas := k.MaxContractGas(ctx); remaining > maxGas {
		remaining = maxGas
	}
	return remaining * types.GasMultiplier
}

// converts contract gas usage to sdk gas and consumes it
func (k Keeper) consumeGas(ctx sdk.Context, gas uint64, descriptor string) {
	consumed := gas / types.GasMultiplier
	ctx.GasMeter().ConsumeGas(consumed, descriptor)
}
