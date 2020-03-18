package tracking

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	core "github.com/terra-project/core/types"
	"github.com/terra-project/core/x/auth"
	ns "github.com/terra-project/core/x/nameservice"
)

// NameserviceHook - track register unregister
func NameserviceHook(k ns.Keeper, accKeeper auth.AccountKeeper) core.HookHandler {
	return func(ctx sdk.Context, msg sdk.Msg, res sdk.Result) {
		if registerMsg, ok := msg.(ns.MsgRegisterSubName); ok {
			nameHash, registry, err := getRegistry(ctx, k, registerMsg.Address)
			if err != nil {
				return
			}

			acc := accKeeper.GetAccount(ctx, registerMsg.Address)
			if acc.GetCoins().IsZero() {
				return
			}

			registry.LockedValue = registry.LockedValue.Add(acc.GetCoins())
			k.SetRegistry(ctx, nameHash, registry)
		} else if _, ok := msg.(ns.MsgUnregisterSubName); ok {
			var addr sdk.AccAddress
			event, err := findEvent(res.Events, ns.EventTypeUnregister)
			if err != nil {
				return
			}

			attr, err := findAttribute(event, ns.AttributeKeyAddress)
			if err != nil {
				return
			}

			addr, err = sdk.AccAddressFromBech32(string(attr))
			if err != nil {
				return
			}

			nameHash, registry, err := getRegistry(ctx, k, addr)
			if err != nil {
				return
			}

			acc := accKeeper.GetAccount(ctx, addr)
			if acc.GetCoins().IsZero() {
				return
			}

			registry.LockedValue = registry.LockedValue.Sub(acc.GetCoins())
			k.SetRegistry(ctx, nameHash, registry)
		}
	}
}
func findAttribute(event sdk.Event, attrName string) ([]byte, error) {
	for _, attr := range event.Attributes {
		if string(attr.Key) == attrName {
			return attr.Value, nil
		}
	}
	return nil, errors.New("attribute is not found")
}

func findEvent(events sdk.Events, eventType string) (sdk.Event, error) {
	for _, event := range events {
		if event.Type == eventType {
			return event, nil
		}
	}

	return sdk.Event{}, errors.New("event is not found")
}

func getRegistry(ctx sdk.Context, k ns.Keeper, accAddr sdk.AccAddress) (nameHash ns.NameHash, registry ns.Registry, err error) {
	nameHash, err = k.GetReverseResolve(ctx, accAddr)
	if err != nil {
		return nil, ns.Registry{}, err
	}

	registry, err = k.GetRegistry(ctx, nameHash)

	return
}
