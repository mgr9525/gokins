package mgr

import "context"

var (
	mgrCtx context.Context
	Cancel context.CancelFunc
)

func init() {
	mgrCtx, Cancel = context.WithCancel(context.Background())
}
