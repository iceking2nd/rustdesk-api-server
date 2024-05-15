package Debugger

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/global"
	"io"
)

func RequestLogger() gin.HandlerFunc {
	const MESSAGE = "debug request logger"
	log := global.Log.WithField("function", "app.Middlewares.Debugger.RequestLogger")
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf)
		log.WithField("request_uri", c.Request.RequestURI).Debugln(MESSAGE)
		log.WithField("headers", c.Request.Header).Debugln(MESSAGE)
		log.WithField("method", c.Request.Method).Debugln(MESSAGE)
		log.WithField("body", string(body)).Debugln(MESSAGE)
		c.Next()
	}
}

func requestLoggerReadBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(reader)

	s := buf.String()
	return s
}
