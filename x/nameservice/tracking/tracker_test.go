package tracking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	ns "github.com/terra-project/core/x/nameservice"
	"github.com/terra-project/core/x/nameservice/internal/keeper"
	"testing"
	"time"
)

func TestNameserviceHook(t *testing.T) {
	input := keeper.CreateTestInput(t)

	name := ns.Name("valid.terra")
	nameHash, childHash := name.NameHash()
	registry := ns.NewRegistry(name, keeper.Addrs[0], time.Now().UTC())
	input.NameserviceKeeper.SetRegistry(input.Ctx, nameHash, registry)
	input.NameserviceKeeper.SetResolve(input.Ctx, nameHash, childHash, keeper.Addrs[0])
	input.NameserviceKeeper.SetReverseResolve(input.Ctx, keeper.Addrs[0], nameHash)

	registerMsg := ns.MsgRegisterSubName{Name: name, Address: keeper.Addrs[0], Owner: keeper.Addrs[0]}
	unregisterMsg := ns.MsgUnregisterSubName{Name: name, Owner: keeper.Addrs[0]}

	hook := NameserviceHook(input.NameserviceKeeper, input.AccKeeper)
	hook(input.Ctx, registerMsg, sdk.Result{})

	registry, err := input.NameserviceKeeper.GetRegistry(input.Ctx, nameHash)
	require.NoError(t, err)
	require.True(t, registry.LockedValue.IsEqual(keeper.InitCoins))

	// without event tags, no state changes occur
	hook(input.Ctx, unregisterMsg, sdk.Result{})

	registry, err = input.NameserviceKeeper.GetRegistry(input.Ctx, nameHash)
	require.NoError(t, err)
	require.True(t, registry.LockedValue.IsEqual(keeper.InitCoins))

	// valid event tags
	input.Ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			ns.EventTypeUnregister,
			sdk.NewAttribute(ns.AttributeKeyAddress, keeper.Addrs[0].String()),
		),
	)

	hook(input.Ctx, unregisterMsg, sdk.Result{Events: input.Ctx.EventManager().Events()})

	registry, err = input.NameserviceKeeper.GetRegistry(input.Ctx, nameHash)
	require.NoError(t, err)
	require.True(t, registry.LockedValue.IsZero())
}
