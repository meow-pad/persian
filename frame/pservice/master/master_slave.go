package master

import (
	"context"
	"errors"
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/frame/plog/pfield"
	"github.com/meow-pad/persian/frame/pservice/cache"
	"github.com/meow-pad/persian/utils/coding"
	"github.com/meow-pad/persian/utils/timewheel"
	"sync/atomic"
	"time"
)

type MSHandler interface {
	OnBeMainService() error
	OnBeMainServiceFailed()
	OnLeaveMainService()
	OnTaskData(dataType int32, data any)
	OnKeepTick()
}

const (
	taskTypeKeep   = 1
	taskTypeCustom = 2
)

type MSTask struct {
	taskType int32
	dataType int32
	data     any
}

func NewMSService(opts ...Option) (*MSService, error) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	if err := options.check(); err != nil {
		return nil, err
	}
	return &MSService{Options: options}, nil
}

type MSService struct {
	*Options

	checkTask     *timewheel.Task
	mainSrv       atomic.Bool // 是否是维护数据的主服务
	mainSrvInstId atomic.Pointer[string]
	taskChan      chan *MSTask
	closed        atomic.Bool
}

func (srv *MSService) Start(ctx context.Context) error {
	srv.taskChan = make(chan *MSTask, 1)
	go srv.running()
	srv.checkTask = srv.TWTimer.AddCron(time.Duration(srv.TickIntervalSec)*time.Second, srv.timeTick)
	// 尝试成为主服务
	if err := srv.tryToBeMainService(true); err != nil {
		plog.Info("try to be main service failed",
			pfield.String("SrvName", srv.SrvName),
			pfield.Error(err))
	}
	return nil
}

func (srv *MSService) Stop(ctx context.Context) error {
	if !srv.closed.CompareAndSwap(false, true) {
		return nil
	}
	if err := srv.TWTimer.Remove(srv.checkTask); err != nil {
		plog.Error("remove check first island task failed", pfield.Error(err))
	}
	close(srv.taskChan)
	// 关闭时删除缓存
	if srv.mainSrv.Load() {
		srv.deleteDCache()
	}
	return nil
}

// tryToBeMainService
//
//	@Description: io阻塞执行成为主服务
//	@receiver srv
//	@return error
func (srv *MSService) tryToBeMainService(deleteCacheOnErr bool) error {
	_, err := srv.Cache.AddOrUpdate(srv.DistributionCacheKey, srv.ServiceId, srv.DistributionCacheSignature,
		srv.DistributionCacheExpireSec, 0, 0, 0, nil)
	if err != nil {
		return err
	}
	// 执行成为主服务的准备操作
	if err = srv.Handler.OnBeMainService(); err != nil {
		srv.Handler.OnBeMainServiceFailed()
		if deleteCacheOnErr {
			// 失败了则删除缓存返回
			srv.deleteDCache()
		}
		return err
	}
	srv.mainSrv.Store(true)
	return nil
}

// leaveMainService
//
//	@Description: io阻塞执行离开主服务
//	@receiver srv
//	@param deleteDCache
func (srv *MSService) leaveMainService(deleteDCache bool) {
	if deleteDCache {
		// 删除分布式缓存
		srv.deleteDCache()
	}
	srv.Handler.OnLeaveMainService()
	srv.mainSrv.Store(false)
}

// keepMainService
//
//	@Description: 尝试保持主服务状态
//	@receiver srv
//	@return error
func (srv *MSService) keepMainService() error {
	defer srv.Handler.OnKeepTick()
	if srv.mainSrv.Load() {
		srv.mainSrvInstId.Store(nil)
		// 如果是主服务则续写缓存
		_, err := srv.Cache.AddOrUpdate(srv.DistributionCacheKey, srv.ServiceId, srv.DistributionCacheSignature,
			srv.DistributionCacheExpireSec, 0, 0, 0, nil)
		if err != nil {
			if errors.Is(err, cache.ErrCacheExist) {
				// 缓存中存在其他服务，且自身状态是主服务，则清理状态
				srv.leaveMainService(false)
				return nil
			} else {
				// 直接退出，不退出主服务，否则下次选举期间没有可用服务
				return err
			}
		}
		return nil
	} else {
		// 判定当前是否有主服务，没有则尝试成为主服务
		cServiceId, exist, err := srv.Cache.Get(srv.DistributionCacheKey)
		if err != nil {
			return err
		}
		if exist {
			srv.mainSrvInstId.Store(&cServiceId)
			if cServiceId == srv.ServiceId {
				// 缓存中存在当前服务，却不是主服务，可能是重启了或者上次续约失败了，尝试成为主服务
				if err = srv.tryToBeMainService(false); err != nil {
					// 失败了则删除缓存返回,等待新的主服务产生
					srv.deleteDCache()
					return err
				}
			} else {
				// 缓存中存在其他服务，且自身不是主服务，则不处理
			}
			return nil
		} else {
			srv.mainSrvInstId.Store(nil)
			// 缓存中不存在主服务，尝试成为主服务
			if err = srv.tryToBeMainService(true); err != nil {
				return err
			}
			return nil
		}
	} // end of else
}

// deleteDCache
//
//	@Description: 删除分布式缓存
//	@receiver srv
func (srv *MSService) deleteDCache() {
	cErr := srv.Cache.Delete(srv.DistributionCacheKey, srv.DistributionCacheSignature)
	if cErr != nil {
		plog.Error("delete first_island_rank Cache failed", pfield.Error(cErr))
	}
}

func (srv *MSService) AddTask(dataType int32, data any, waitOrNot bool) bool {
	return srv.addTask(&MSTask{
		taskType: taskTypeCustom,
		dataType: dataType,
		data:     data,
	}, waitOrNot)
}

func (srv *MSService) addTask(task *MSTask, waitOrNot bool) bool {
	if srv.closed.Load() {
		return false
	}
	if waitOrNot {
		srv.taskChan <- task
		return true
	} else {
		select {
		case srv.taskChan <- task:
			return true
		default:
			return false
		}
	}
}

func (srv *MSService) timeTick() {
	if srv.closed.Load() {
		return
	}
	// 定时执行保活
	retryTimes := 2 // 失败重试次数
	for i := 0; i < retryTimes; i++ {
		if srv.addTask(&MSTask{
			taskType: taskTypeKeep,
		}, false) {
			return
		}
		if i != retryTimes-1 {
			time.Sleep(time.Second)
		}
	}
}

func (srv *MSService) running() {
	defer coding.CatchPanicError("first island creation running error:", func() {
		if srv.closed.Load() {
			return
		}
		go srv.running()
	})
	if srv.closed.Load() {
		return
	}
	for {
		select {
		case task, ok := <-srv.taskChan:
			if task == nil || !ok {
				return
			}
			switch task.taskType {
			case taskTypeKeep:
				if err := srv.keepMainService(); err != nil {
					plog.Error("keep first island service failed", pfield.Error(err))
				}
			case taskTypeCustom:
				srv.Handler.OnTaskData(task.dataType, task.data)
			default:
				plog.Error("unknown master-slave service task type",
					pfield.String("SrvName", srv.SrvName),
					pfield.Int32("taskType", task.taskType))
			}
		}
	} // end of for
}

func (srv *MSService) IsMainService() bool {
	return srv.mainSrv.Load()
}

func (srv *MSService) GetMainServiceId() string {
	if srv.IsMainService() {
		return srv.ServiceId
	} else {
		mainSrvId := srv.mainSrvInstId.Load()
		if mainSrvId != nil {
			return *mainSrvId
		} else {
			return srv.GetLatestMainServiceId()
		}
	}
}

func (srv *MSService) GetLatestMainServiceId() string {
	if srv.IsMainService() {
		srv.mainSrvInstId.Store(nil)
		return srv.ServiceId
	} else {
		cServiceId, exist, err := srv.Cache.Get(srv.DistributionCacheKey)
		if err != nil {
			plog.Error("get first island_rank Cache failed", pfield.Error(err))
			return ""
		}
		if exist {
			srv.mainSrvInstId.Store(&cServiceId)
			return cServiceId
		} else {
			srv.mainSrvInstId.Store(nil)
			return ""
		}
	}
}
