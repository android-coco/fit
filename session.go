package fit

import (
	jkl_fmt "fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	// private session provider storage, which element implemented SessionProvider interface methods
	g_Providers = make(map[string]SessionProvider)
	// Private global session manager
	g_GlobalSessionManager *SessionManager
)

type SessionProvider interface {
	// Initialize session with sid, then return the session
	SessionInit() (Session, error)
	// Get the session for sid, call SessionInit(sid string) if the session not exist
	SessionRead(sid string) (Session, error)
	// Destroy the session with sid
	SessionDestroy(sid string) error
	// Delete expired session - garbage collection
	RemoveAllSession()
}

// Session storage
type Session interface {
	Set(key, value interface{})      //set session value
	Get(key interface{}) interface{} //get session value
	Delete(key interface{}) error    //delete session value
	SessionID() string               //back current sessionID
}

type SessionManager struct {
	cookieName string          // private
	lock       sync.Mutex      // protects session while writing
	provider   SessionProvider // practical session storage
}

// Registering provider with providerName, can refer to init() in memory.go file
// Provider must implement SessionProvider interface methods
func Register(name string, provider SessionProvider) {
	if provider == nil {
		panic("Session Panic: Register provider is nil")
	}
	// Repeated registration
	if _, dup := g_Providers[name]; dup {
		panic("Session Panic: Register called twice for provider " + name)
	}
	g_Providers[name] = provider
}

// Get the global session manager, which is singleton
func GlobalManager() *SessionManager {
	if g_GlobalSessionManager == nil {
		g_GlobalSessionManager, _ = NewSessionManager("memory", Config().SessionKey, Config().SessionTimeout)
	}
	return g_GlobalSessionManager
}

// New session manager, but the provider must be registered with providername
func NewSessionManager(providerName, cookieName string, maxlifetime time.Duration) (*SessionManager, error) {
	provider, ok := g_Providers[providerName]
	if !ok {
		return nil, jkl_fmt.Errorf("Session Error: unknown session provider %q (forgotten register ?)", providerName)
	} else if maxlifetime <= 0 {
		return nil, jkl_fmt.Errorf("Session Error: the maxlifetime is invalid, value is %d", maxlifetime)
	} else {
		provider.(*Provider).maxlifetime = maxlifetime
		return &SessionManager{
			cookieName: cookieName,
			provider:   provider,
		}, nil
	}
}

// initialize session with http.reqeust , usually used in login API
func (manager *SessionManager) SessionStart(w *Response, r *Request) (Session, error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	// generate sessionID if can't find the cookie
	ck, err := r.Cookie(manager.cookieName)
	//jkl_fmt.Println("JK Cookie:  ", ck)
	if err != nil || ck.Value == "" {
		// initialize session and then set cookie
		maxlifetime := int(manager.provider.(*Provider).maxlifetime)
		session, err_dup := manager.provider.SessionInit()
		cookie_dup := http.Cookie{
			Name:     manager.cookieName,
			Value:    url.QueryEscape(session.SessionID()),
			Path:     "/",
			HttpOnly: true, // prevent the session hijack
			MaxAge:   maxlifetime,
		}
		http.SetCookie(w.Writer(), &cookie_dup)
		return session, err_dup
	}

	sid, _ := url.QueryUnescape(ck.Value)
	session, err_dup := manager.provider.SessionRead(sid)
	// sessionId is invalid（out of date）
	if err_dup != nil || session == nil {
		expiration := time.Now().AddDate(-1, 0, 0)
		cookie_dup := http.Cookie{
			Name:     manager.cookieName,
			HttpOnly: true,
			Path:     "/",
			Expires:  expiration, //
			MaxAge:   -1,         // client will clear local cookie_dup at once
		}
		http.SetCookie(w.Writer(), &cookie_dup)
	}

	return session, err_dup
}

// destroy sessionid, usually used in logout API.
func (manager *SessionManager) SessionDestroy(w *Response, r *Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(cookie.Value)
		expiration := time.Now().AddDate(-1, 0, 0)
		cookie := http.Cookie{
			Name:     manager.cookieName,
			Path:     "/",
			HttpOnly: true,
			Expires:  expiration, //
			MaxAge:   -1,         // client will clear local cookie at once
		}
		http.SetCookie(w.Writer(), &cookie)
	}
}

// Remove All Session， Prudent Operation!
func (manager *SessionManager) RemoveAllSession() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.RemoveAllSession()
}
