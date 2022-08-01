package common

import (
	"fmt"
	pbAct "github.com/dukhyungkim/libharago/gen/go/proto/action"
)

var (
	ErrDuplicateCommand  = func(name string) error { return fmt.Errorf("duplicate command name: %s", name) }
	ErrConnEtcd          = func(err error) error { return fmt.Errorf("failed to connect etcd; %w", err) }
	ErrHarborResponse    = func(err error) error { return fmt.Errorf("HarborError!: %v", err) }
	ErrUnknownActionType = func(actionType pbAct.ActionType) error { return fmt.Errorf("unknown action type: %v", actionType) }
)
