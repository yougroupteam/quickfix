// Copyright (c) quickfixengine.org  All rights reserved.
//
// This file may be distributed under the terms of the quickfixengine.org
// license as defined by quickfixengine.org and appearing in the file
// LICENSE included in the packaging of this file.
//
// This file is provided AS IS with NO WARRANTY OF ANY KIND, INCLUDING
// THE WARRANTY OF DESIGN, MERCHANTABILITY AND FITNESS FOR A
// PARTICULAR PURPOSE.
//
// See http://www.quickfixengine.org/LICENSE for licensing information.
//
// Contact ask@quickfixengine.org if any conditions of this licensing
// are not clear to you.

package quickfix

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/proxy"

	"github.com/yougroupteam/quickfix/config"
)

func loadDialerConfig(settings *SessionSettings) (dialer proxy.ContextDialer, err error) {
	stdDialer := &net.Dialer{}
	if settings.HasSetting(config.SocketTimeout) {
		timeout, err := settings.DurationSetting(config.SocketTimeout)
		if err != nil {
			timeoutInt, err := settings.IntSetting(config.SocketTimeout)
			if err != nil {
				return stdDialer, err
			}

			stdDialer.Timeout = time.Duration(timeoutInt) * time.Second
		} else {
			stdDialer.Timeout = timeout
		}
	}
	dialer = stdDialer

	if !settings.HasSetting(config.ProxyType) {
		return
	}

	var proxyType string
	if proxyType, err = settings.Setting(config.ProxyType); err != nil {
		return
	}

	switch proxyType {
	case "socks":
		var proxyHost string
		var proxyPort int
		if proxyHost, err = settings.Setting(config.ProxyHost); err != nil {
			return
		} else if proxyPort, err = settings.IntSetting(config.ProxyPort); err != nil {
			return
		}

		proxyAuth := new(proxy.Auth)
		if settings.HasSetting(config.ProxyUser) {
			if proxyAuth.User, err = settings.Setting(config.ProxyUser); err != nil {
				return
			}
		}
		if settings.HasSetting(config.ProxyPassword) {
			if proxyAuth.Password, err = settings.Setting(config.ProxyPassword); err != nil {
				return
			}
		}

		var proxyDialer proxy.Dialer

		proxyDialer, err = proxy.SOCKS5("tcp", fmt.Sprintf("%s:%d", proxyHost, proxyPort), proxyAuth, stdDialer)
		if err != nil {
			return
		}

		if contextDialer, ok := proxyDialer.(proxy.ContextDialer); ok {
			dialer = contextDialer
		} else {
			err = fmt.Errorf("proxy does not support context dialer")
			return
		}

	default:
		err = fmt.Errorf("unsupported proxy type %s", proxyType)
	}

	return
}
