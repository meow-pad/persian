package pboot

import (
	"context"
	"github.com/go-spring/spring-core/gs"
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/frame/plog/pfield"
	"github.com/meow-pad/persian/utils/coding"
	"go.uber.org/zap"
)

type LifeCycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	CName() string
}

// lifeCycleWrapper
type lifeCycleWrapper struct {
	lc LifeCycle
}

func (lf *lifeCycleWrapper) OnAppStart(_ gs.Context) {
	// nothing
}

func (lf *lifeCycleWrapper) OnAppStop(_ context.Context) {
	// nothing
}

func (lf *lifeCycleWrapper) Start(ctx context.Context) error {
	return lf.lc.Start(ctx)
}

func (lf *lifeCycleWrapper) Stop(ctx context.Context) error {
	return lf.lc.Stop(ctx)
}

func (lf *lifeCycleWrapper) CName() string {
	return lf.lc.CName()
}

func initLifeCycleMgr() {
	Object(new(lifeCycleManager)).Order(OrderMax).Export((*gs.AppEvent)(nil))
}

// lifeCycleManager
//
//	@Description: 生命周期管理器
type lifeCycleManager struct {
}

func (mgr *lifeCycleManager) OnAppStart(gsCtx gs.Context) {
	ctx := gsCtx.Context()
	for _, event := range app.Events {
		lc, _ := event.(LifeCycle)
		if lc != nil {
			err := coding.SafeRunWithContext(lc.Start, ctx)
			if err != nil {
				plog.Panic("starting error:", pfield.String("module", lc.CName()), zap.Error(err))
			} else {
				plog.Info("starting success:", pfield.String("module", lc.CName()))
			}
		}
	} // end of for
}

func (mgr *lifeCycleManager) OnAppStop(ctx context.Context) {
	for i := len(app.Events) - 1; i >= 0; i-- {
		event := app.Events[i]
		lc, _ := event.(LifeCycle)
		if lc != nil {
			err := coding.SafeRunWithContext(lc.Stop, ctx)
			if err != nil {
				plog.Error("stopping error:", pfield.String("module", lc.CName()), zap.Error(err))
			} else {
				plog.Info("stopping success:", pfield.String("module", lc.CName()))
			}
		}
	} // end of for
}
