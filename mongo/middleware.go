package mongo

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
)

// Middleware -
type Middleware interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type middleware struct {
	session  *mgo.Session
	mongoURL *url.URL
}

// NewMiddleware -
func NewMiddleware(mongoURL string, useSSL bool) (Middleware, error) {
	url, err := url.Parse(mongoURL)
	if err != nil {
		return nil, err
	}

	fmt.Println("URL: " + mongoURL)

	s := &mgo.Session{}
	if useSSL {
		var err error
		s, err = createSessionWithSSL(url)
		if err != nil {
			return nil, err
		}
	} else {
		s, err = mgo.Dial(url.String())
		if err != nil {
			return nil, err
		}
	}

	return &middleware{s, url}, nil
}

func (mw *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	s := mw.session.Copy()
	context := SetInCtx(r.Context(), s.DB(strings.TrimPrefix(mw.mongoURL.Path, "/")))
	next(w, r.WithContext(context))
	defer s.Close()
}

func createSessionWithSSL(url *url.URL) (*mgo.Session, error) {
	info := &mgo.DialInfo{
		Addrs:    strings.Split(url.Host, ","),
		Database: strings.TrimPrefix(url.Path, "/"),
		Timeout:  10 * time.Second,
	}

	if url.User != nil {
		info.Username = url.User.Username()
		info.Password, _ = url.User.Password()
	}

	info.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), &tls.Config{})
	}

	s, err := mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
	}

	return s, nil
}
