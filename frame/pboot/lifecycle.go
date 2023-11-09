package pboot

import (
	"context"
	"github.com/go-spring/spring-core/gs"
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/frame/plog/pfield"
	"github.com/meow-pad/persian/utils/coding"
	"go.uber.org/zap"
	"sort"
	"strings"
)

type LifeCycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	CName() string
}

func initLifeCycleMgr() {
	Object(new(lifeCycleManager)).Init(func(lcMgr *lifeCycleManager) error {
		return lcMgr.init()
	}).Order(OrderMax).Export((*gs.AppEvent)(nil))
}

var (
	lifeCycleMap = make(map[LifeCycle]*Bean)
)

func addLifeCycle(lc LifeCycle, bean *Bean) {
	lifeCycleMap[lc] = bean
}

// lifeCycleManager
//
//	@Description: 生命周期管理器
type lifeCycleManager struct {
	lcList []LifeCycle
}

func (mgr *lifeCycleManager) init() error {
	lcList := make([]LifeCycle, 0, len(lifeCycleMap))
	for lc := range lifeCycleMap {
		lcList = append(lcList, lc)
	}
	sort.SliceStable(lcList, func(i, j int) bool {
		beanI := lifeCycleMap[lcList[i]]
		beanJ := lifeCycleMap[lcList[j]]
		if beanI.bOrder != beanJ.bOrder {
			return beanI.bOrder < beanJ.bOrder
		}
		return strings.Compare(lcList[i].CName(), lcList[j].CName()) == -1
	})
	mgr.lcList = lcList
	return nil
}

func (mgr *lifeCycleManager) OnAppStart(gsCtx gs.Context) {
	ctx := gsCtx.Context()
	for _, lc := range mgr.lcList {
		plog.Debug("start lifecycle:" + lc.CName())
		err := coding.SafeRunWithContext(lc.Start, ctx)
		if err != nil {
			plog.Panic("starting error:", pfield.String("module", lc.CName()), zap.Error(err))
		} else {
			plog.Info("starting success:", pfield.String("module", lc.CName()))
		}
	} // end of for
}

func (mgr *lifeCycleManager) OnAppStop(ctx context.Context) {
	for i := len(mgr.lcList) - 1; i >= 0; i-- {
		lc := mgr.lcList[i]
		plog.Debug("stop lifecycle:" + lc.CName())
		err := coding.SafeRunWithContext(lc.Stop, ctx)
		if err != nil {
			plog.Error("stopping error:", pfield.String("module", lc.CName()), zap.Error(err))
		} else {
			plog.Info("stopping success:", pfield.String("module", lc.CName()))
		}
	} // end of for
}
