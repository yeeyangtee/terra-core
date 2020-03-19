package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Registry - data struct to hold name information
type Registry struct {
	Name                Name           `json:"name" yaml:"name"`
	Owner               sdk.AccAddress `json:"owner" yaml:"owner"`
	EndTime             time.Time      `json:"end_time" yaml:"end_time"`
	LockedValue         sdk.Coins      `json:"locked_value" yaml:"locked_value"`
	TaxContribution     sdk.Coins      `json:"tax_contribution" yaml:"tax_contribution"`
	SwapFeeContribution sdk.DecCoins   `json:"swap_fee_contribution" yaml:"swap_fee_contribution"`
}

// NewRegistry returns Registry instance
func NewRegistry(name Name, owner sdk.AccAddress, endTime time.Time) Registry {
	return Registry{
		Name:                name,
		Owner:               owner,
		EndTime:             endTime,
	}
}

// String implements fmt.Stringer interface
func (r Registry) String() string {
	return fmt.Sprintf(`Registry
Name:                %s
Owner:               %s
EndTime:             %s 
LockValue:           %s
TaxContribution:     %s
SwapFeeContribution: %s
`, r.Name, r.Owner.String(), r.EndTime, r.LockedValue, r.TaxContribution, r.SwapFeeContribution)
}
