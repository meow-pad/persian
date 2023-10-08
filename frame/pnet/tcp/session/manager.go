package session

import (
	"github.com/meow-pad/persian/errdef"
	"github.com/meow-pad/persian/frame/plog"
	"github.com/meow-pad/persian/frame/plog/cfield"
	"github.com/meow-pad/persian/frame/pnet"
	"reflect"
	"sync"
	"time"
)

func NewManager(name string, unregisterSessionLife int64) (*Manager, error) {
	if unregisterSessionLife <= 0 {
		return nil, errdef.ErrInvalidParams
	}
	return &Manager{name: name, unregisterSessionLife: unregisterSessionLife}, nil
}

type Manager struct {
	name                  string
	unregisterSessionLife int64

	// 未注册会话集合
	unregisterSessions sync.Map
	// 已注册会话集合
	registerSessions sync.Map
}

// AddSession
//
//	@Description: 添加未注册的会话
//	@receiver manager
//	@param Session 会话
//	@return error
func (manager *Manager) AddSession(svrSess Session) error {
	if svrSess.Context() != nil {
		plog.Warn("a new session has one context???")
		sessionId := svrSess.Id()
		if sessionId == InvalidSessionId {
			plog.Warn("id of session with context is invalid")
			// 没有Id的context直接清理掉
			svrSess.setContext(nil)
		} else {
			if _, ok := manager.registerSessions.Load(sessionId); !ok {
				plog.Warn("session with context is not in registerSessions")
				// 不在注册集合中，清理原注册信息
				svrSess.setContext(nil)
			} else {
				// 都已经完全注册了
				plog.Warn("session with context is already in registerSessions")
				return nil
			}
		}
	}
	manager.unregisterSessions.Store(svrSess, time.Now().Unix()+manager.unregisterSessionLife)
	return nil
}

// RemoveSession
//
//	@Description: 移除会话
//	@receiver manager
//	@param svrSess
func (manager *Manager) RemoveSession(svrSess Session) {
	if svrSess.Context() == nil {
		manager.unregisterSessions.Delete(svrSess)
	} else {
		sessionId := svrSess.Id()
		if sessionId != InvalidSessionId {
			manager.registerSessions.Delete(sessionId)
		}
	}
}

// RegisterSession
//
//	@Description: 注册会话
//	@receiver manager
//	@param Session 会话
//	@return error
func (manager *Manager) RegisterSession(svrSess Session, context Context) error {
	if context == nil {
		return errdef.ErrInvalidParams
	}
	if svrSess.Context() != nil {
		return pnet.ErrRegisteredSession
	}
	svrSess.setContext(context)
	if _, loaded := manager.unregisterSessions.LoadAndDelete(svrSess); !loaded {
		// ？？不在集合里
		plog.Error("cant find Session in unregisterSessions")
	}
	sessionId := svrSess.Id()
	if sessionId == InvalidSessionId {
		return pnet.ErrInvalidSessionId
	}
	if value, _ := manager.registerSessions.Load(sessionId); value != nil {
		if oldSession, ok := value.(Session); !ok {
			plog.Error("there is invalid value in registerSessions",
				cfield.String("valueType", reflect.TypeOf(value).String()))
		} else {
			// 关闭旧会话
			if err := oldSession.Close(); err != nil {
				plog.Error("close session error:", cfield.Error(err))
			}
		} // end of else
	}
	manager.registerSessions.Store(sessionId, svrSess)
	return nil
}

// GetSession
//
//	@Description: 获取会话
//	@receiver manager
//	@param sessionId
//	@return Session
func (manager *Manager) GetSession(sessionId uint64) Session {
	value, _ := manager.registerSessions.Load(sessionId)
	if value == nil {
		return nil
	}
	if sess, ok := value.(Session); !ok {
		manager.registerSessions.Delete(sessionId)
		plog.Error("invalid session type", cfield.String("type", reflect.TypeOf(value).String()))
		return nil
	} else {
		return sess
	}
}

// CheckSessions
//
//	@Description: 检查会话的有效性
//	@receiver manager
func (manager *Manager) CheckSessions() {
	manager.checkUnregisteredSessions()
	manager.checkRegisteredSessions()
}

// checkUnregisteredSessions
//
//	@Description: 检查未注册会话的是否过期
//	@receiver manager
func (manager *Manager) checkUnregisteredSessions() {
	now := time.Now().Unix()
	manager.unregisterSessions.Range(func(key, value any) bool {
		var (
			overdueTime int64
			svrSess     Session
			ok          bool
		)
		if overdueTime, ok = value.(int64); !ok {
			manager.unregisterSessions.Delete(key)
			plog.Error("there is invalid type value in unregisterSessions")
			return true
		}
		if overdueTime <= now {
			manager.unregisterSessions.Delete(key)
			if svrSess, ok = key.(Session); ok {
				plog.Debug("unregistered session has expired:",
					cfield.String("manager", manager.name),
					cfield.Uint64("conn", svrSess.Connection().Hash()))
				// 过期关闭
				if err := svrSess.Close(); err != nil {
					plog.Error("close session error:", cfield.Error(err))
				}
			} else {
				plog.Error("there is invalid type key in unregisterSessions",
					cfield.String("keyType", reflect.TypeOf(key).String()))
			}
			return true
		}
		return true
	})
}

// checkRegisteredSessions
//
//	@Description: 检查已注册会话是否已失活
//	@receiver manager *TcpServer
func (manager *Manager) checkRegisteredSessions() {
	now := time.Now().Unix()
	manager.registerSessions.Range(func(key any, value any) bool {
		if svrSess, ok := value.(Session); ok {
			ctx := svrSess.Context()
			if ctx == nil {
				plog.Error("nil context in registered sessions")
			}
			if ctx == nil || now >= ctx.Deadline() {
				plog.Debug("registered session has expired:",
					cfield.String("manager", manager.name),
					cfield.Uint64("id", svrSess.Id()))
				if svrSess.IsClosed() {
					// 已关闭则直接移除
					manager.registerSessions.Delete(key)
				} else {
					// 失活关闭
					if err := svrSess.Close(); err != nil {
						plog.Error("close deadline session error:", cfield.Error(err))
					}
					// 等到连接关闭时再移除该session
				}
			}
		} else {
			manager.registerSessions.Delete(key)
			plog.Error("there is invalid type value in unregisterSessions",
				cfield.String("valueType", reflect.TypeOf(value).String()))
		}
		return true
	})
}
