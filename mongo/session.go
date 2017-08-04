package mongo

import (
	"crypto/tls"
	"net"
	"net/url"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
)

// CreateSession returns a new mongo session and its DB name
func CreateSession(mongoURL string, useSSL bool) (*mgo.Session, string, error) {
	url, err := url.Parse(mongoURL)
	if err != nil {
		return nil, "", err
	}

	s := &mgo.Session{}
	if useSSL {
		var err error
		s, err = createSessionWithSSL(url)
		if err != nil {
			return nil, "", err
		}
	} else {
		s, err = mgo.Dial(url.String())
		if err != nil {
			return nil, "", err
		}
	}

	return s, strings.TrimPrefix(url.Path, "/"), nil
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
