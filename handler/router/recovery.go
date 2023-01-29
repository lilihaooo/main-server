// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

func Recovery(c *gin.Context) {
	defer Recover("server", "Recovery", func(err interface{}) {
		// Check for a broken connection, as it is not really a
		// condition that warrants a panic stack trace.
		var brokenPipe bool
		if ne, ok := err.(*net.OpError); ok {
			if se, ok := ne.Err.(*os.SyscallError); ok {
				if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
					brokenPipe = true
				}
			}
		}
		log.Error().Interface("trace", string(debug.Stack())).Interface("err", err).Msg("recovery")
		if brokenPipe {
			// If the connection is dead, we can't write a status to it.
			c.Error(err.(error)) // nolint: errcheck
			c.Abort()
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	})
	c.Next()
}

type RecoverHandler func(p interface{})

// Recover will recovers a given method of app with panic counter metric saved
func Recover(app, method string, postHandlers ...RecoverHandler) {
	if p := recover(); p != nil {
		if len(postHandlers) > 0 {
			for _, postRecoverFn := range postHandlers {
				postRecoverFn(p)
			}
		}
	}
}
