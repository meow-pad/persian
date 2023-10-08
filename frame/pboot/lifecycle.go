package pboot

import (
	"context"
	"github.com/go-spring/spring-core/gs"
	"go.uber.org/zap"
	"persian/errdef"
	"persian/frame/plog"
	"persian/frame/plog/cfield"
	"persian/utils/coding"
)

type LifeCycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	CName() string
}

type LifeCycleBase struct {
}

func (lf *LifeCycleBase) OnAppStart(_ gs.Context) {
	// dummy
}

func (lf *LifeCycleBase) OnAppStop(_ context.Context) {
	// dummy
}

func (lf *LifeCycleBase) Start(context.Context) error {
	return errdef.ErrNotImplemented
}

func (lf *LifeCycleBase) Stop(context.Context) error {
	return errdef.ErrNotImplemented
}

func (lf *LifeCycleBase) CName() string {
	return ""
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
				plog.Panic("starting error:", cfield.String("module", lc.CName()), zap.Error(err))
			} else {
				plog.Info("starting success:", cfield.String("module", lc.CName()))
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
				plog.Panic("stopping error:", cfield.String("module", lc.CName()), zap.Error(err))
			} else {
				plog.Info("stopping success:", cfield.String("module", lc.CName()))
			}
		}
	} // end of for
}
