package types

import (
	"fmt"

	wasmvm "github.com/CosmWasm/wasmvm"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
)

// Model is a struct that holds a KV pair
type Model struct {
	Key   core.Base64Bytes `json:"key"`
	Value core.Base64Bytes `json:"value"`
}

// String implements fmt.Stringer interface
func (m Model) String() string {
	return fmt.Sprintf(`Model
	Key:   %s,
	Value: %s`,
		m.Key, m.Value)
}

// CodeInfo is data for the uploaded contract WASM code
type CodeInfo struct {
	CodeID    uint64           `json:"code_id"`
	CodeHash  core.Base64Bytes `json:"code_hash"`
	Creator   sdk.AccAddress   `json:"creator"`
	VMVersion wasmvm.VMVersion `json:"vm_version"`
}

// String implements fmt.Stringer interface
func (ci CodeInfo) String() string {
	return fmt.Sprintf(`CodeInfo
	CodeID:      %d,
	CodeHash:    %s, 
	Creator:     %s,
	VMVersion:   %d`,
		ci.CodeID, ci.CodeHash, ci.Creator, ci.VMVersion)
}

// NewCodeInfo fills a new Contract struct
func NewCodeInfo(codeID uint64, codeHash []byte, creator sdk.AccAddress, vmVersion wasmvm.VMVersion) CodeInfo {
	return CodeInfo{
		CodeID:    codeID,
		CodeHash:  codeHash,
		Creator:   creator,
		VMVersion: vmVersion,
	}
}

// ContractInfo stores a WASM contract instance
type ContractInfo struct {
	Address    sdk.AccAddress   `json:"address"`
	Owner      sdk.AccAddress   `json:"owner"`
	CodeID     uint64           `json:"code_id"`
	InitMsg    core.Base64Bytes `json:"init_msg"`
	Migratable bool             `json:"migratable"`
}

// NewContractInfo creates a new instance of a given WASM contract info
func NewContractInfo(codeID uint64, address, owner sdk.AccAddress, initMsg []byte, migratable bool) ContractInfo {
	return ContractInfo{
		Address:    address,
		CodeID:     codeID,
		Owner:      owner,
		InitMsg:    initMsg,
		Migratable: migratable,
	}
}

// String implements fmt.Stringer interface
func (ci ContractInfo) String() string {
	return fmt.Sprintf(`ContractInfo
	Address:    %s,
	CodeID:     %d, 
	Owner:      %s,
	InitMsg:    %s,
	Migratable  %v,
	`,
		ci.Address, ci.CodeID, ci.Owner, ci.InitMsg, ci.Migratable)
}

// NewWasmEnv initializes params for a contract instance
func NewWasmEnv(ctx sdk.Context, contractAddr sdk.AccAddress) wasmvmtypes.Env {
	return wasmvmtypes.Env{
		Block: wasmvmtypes.BlockInfo{
			Height:  uint64(ctx.BlockHeight()),
			Time:    uint64(ctx.BlockTime().Unix()),
			ChainID: ctx.ChainID(),
		},
		Contract: wasmvmtypes.ContractInfo{
			Address: contractAddr.String(),
		},
	}
}

// NewWasmInfo initializes the MessageInfo for a contract instance
func NewWasmInfo(creator sdk.AccAddress, deposit sdk.Coins) wasmvmtypes.MessageInfo {
	return wasmvmtypes.MessageInfo{
		Sender:    creator.String(),
		SentFunds: NewWasmCoins(deposit),
	}
}

// NewWasmCoins translates between Cosmos SDK coins and Wasm coins
func NewWasmCoins(cosmosCoins sdk.Coins) (wasmCoins []wasmvmtypes.Coin) {
	for _, coin := range cosmosCoins {
		wasmCoin := wasmvmtypes.Coin{
			Denom:  coin.Denom,
			Amount: coin.Amount.String(),
		}
		wasmCoins = append(wasmCoins, wasmCoin)
	}
	return wasmCoins
}
