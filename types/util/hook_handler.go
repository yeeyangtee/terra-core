package util

import sdk "github.com/cosmos/cosmos-sdk/types"

//___________________________
// HookHandler is interface to handle various purpose after actions
type HookHandler func(ctx sdk.Context, msg sdk.Msg, res sdk.Result)
