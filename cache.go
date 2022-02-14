package httpservercache

import (
	"http"
	"net"
	"sync"

	"github.com/gorilla/mux"
)

type Mux = mux.Router // with Go 1.18 this will be a generic parameter

type Cache struct {
	lock       sync.Mutex
	cache      map[string]*Wrapper
	getHandler func(addr string) (mux.Router, error)
}

type Wrapper struct {
	handler  Mux
	users    int
	lock     sync.Mutex
	addr     string
	cache    *Cache
	listener net.Listener
}

func (c *Cache) Get(addr string) (Wrapper, error) {
	c.lock.Lock()
	if wrapper, ok := c.cache[addr]; ok {
		c.lock.Unlock()
		wrapper.lock.Lock()
		defer wrapper.lock.Unlock()
		if wrapper.users <= 0 {
			// This should only happen if users of the wrapper drops to zero
			// after c.lock is released but before wrapper.lock is obtained.
			return c.Get(addr)
		}
		wrapper.users++
		return wrapper
	}
	defer c.lock.Unlock()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "listen")
	}
	r := mux.NewRouter()
	err := http.Serve(listener, r)
	if err != nil {
		return nil, errors.Wrap(err, "http.Serve")
	}
	wrapper := &Wrapper{
		handler:  r,
		users:    1,
		addr:     addr,
		cache:    c,
		listener: listener,
	}
	c.cache[addr] = wrapper
	return wrapper, nil
}

// Close is used to indicate that wrapper returned by Cache.Get() is no
// longer needed by the caller of Cache.Get().  Since Wrappers are cached
// there may be other users.  The underlying listener is stopped when the
// number of active users reaches zero.
func (w *Wrapper) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()
	var err error
	if users == 1 {
		cache.lock.Lock()
		defer cache.lock.Unlock()
		delete(cache.cache, w.addr)
		err = w.listener.Close()
	}
	w.users--
	return err
}

// Unwrap obtains an exclusive lock on the returned handler.  Call
// the unlock function to release the lock.
func (w *Wrapper) Unwrap() (handler Mux, unlock func(), err error) {
	w.lock.Lock()
	if w.users <= 0 {
		return nil, nil, error.Errorf("Handler for %s already closed", w.addr)
	}
	return w.handler, func() {
		w.lock.Unlock()
	}, nil
}

func New(getHandler func(addr string) (http.Handler, error)) *Cache {
	return &Cache{
		getHanlder: getHandler,
		cache:      make(map[string]http.Handler),
	}
}
