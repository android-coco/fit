package session

import (
	jkl_fmt "fmt"
	"sync"
	"time"
)

var (
	// private
	g_Provider  *Provider = &Provider{}
	g_timewheel *TimeWheel
)

type SessionStorage struct {
	sid   string                      // sessionID
	value map[interface{}]interface{} // store
}

// implement Session interface
func (st *SessionStorage) Set(key, value interface{}) {
	// st.sessionUpdate()
	st.value[key] = value
}

func (st *SessionStorage) Get(key interface{}) interface{} {
	// st.sessionUpdate()
	if v, ok := st.value[key]; ok {
		return v
	}
	return nil
}

func (st *SessionStorage) Delete(key interface{}) error {
	// st.sessionUpdate()
	if _, ok := st.value[key]; ok {
		delete(st.value, key)
		return nil
	}
	return jkl_fmt.Errorf("Session Delete Storage Fail! didn't exsit value for key")
}

func (st *SessionStorage) SessionID() string {
	return st.sid
}

// Updating session accesse time
// func (st *SessionStorage) sessionUpdate() {
// 	g_timewheel.UpdateTicker(st.SessionID())
// }

type Provider struct {
	lock        sync.Mutex                 // lock
	sessions    map[string]*SessionStorage // session storage
	maxlifetime time.Duration              // the session will be destroyed if out of date
}

// handle timewheel timeout action, remove session
func (pd *Provider) TimeWheelExpireFunc(sid string) {
	pd.lock.Lock()
	delete(pd.sessions, sid)
	// stop timewheel if sessions map is empty
	if length := len(pd.sessions); length == 0 {
		g_timewheel.Stop()
	}
	pd.lock.Unlock()
}

// Implementing SessionProvider interface methods
// Initializing a new session stroage which implemented Session interface
func (pd *Provider) SessionInit() (Session, error) {
	pd.lock.Lock()
	defer pd.lock.Unlock()

	// add a new task to the timewheel, the sessionId will be returned
	sid, err := g_timewheel.AddTask(pd, pd.maxlifetime)
	if err != nil {
		return nil, jkl_fmt.Errorf("Session initialize Fail! The provider is invalid that has not been registered")
	}

	value := make(map[interface{}]interface{}, 0)
	storage := &SessionStorage{
		sid: sid,
		// timeAccessed: time.Now(),
		value: value,
	}
	pd.sessions[sid] = storage

	// start timewheel sticker
	if g_timewheel.isOn == false {
		g_timewheel.Start()
	}
	return storage, nil
}

func (pd *Provider) SessionRead(sid string) (Session, error) {
	if storage := pd.sessions[sid]; storage != nil {
		// update access time
		g_timewheel.UpdateTicker(sid)
		return storage, nil
	}
	return nil, jkl_fmt.Errorf("SessionId is expired!")
}

func (pd *Provider) SessionDestroy(sid string) error {
	pd.lock.Lock()
	defer pd.lock.Unlock()
	if v := pd.sessions[sid]; v != nil {
		delete(pd.sessions, sid)
		g_timewheel.RemoveTask(sid)

		// stop timewheel while sessions map is empty
		if length := len(pd.sessions); length == 0 {
			g_timewheel.Stop()
		}
		return nil
	}
	return jkl_fmt.Errorf("ession Memory Destroy Error: can't find the session by sid - %s", sid)
}

// Launch garbage collection if out of date
func (pd *Provider) RemoveAllSession() {
	pd.lock.Lock()
	defer pd.lock.Unlock()
	for k, _ := range pd.sessions {
		delete(pd.sessions, k)
		g_timewheel.RemoveTask(k)
	}
	// stop timewheel
	g_timewheel.Stop()
}

// register session provider
func init() {
	g_Provider.sessions = make(map[string]*SessionStorage)
	g_timewheel = NewTimeWheel()
	Register("memory", g_Provider)
}
