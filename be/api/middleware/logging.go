package middleware

import (
	"fmt"
	"net/http"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func makeLogEntry(c echo.Context) *logrus.Entry {
	requestId := uuid.New().String()

	c.Request().Header.Set("X-request-id", requestId)

	return logrus.WithFields(logrus.Fields{
		"request_id": requestId,
		"method": c.Request().Method,
		"uri": c.Request().URL.String(),
		"path": c.Request().URL.Path,
		"query": c.Request().URL.RawQuery,
		"remote_addr": c.Request().RemoteAddr,
	})
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		makeLogEntry(c).Info("Incoming HTTP request")
		return next(c)
	}
}

func ErrorHandler(err error, c echo.Context) {
	report, ok := err.(*echo.HTTPError)
	if ok {
        report.Message = fmt.Sprintf("http error %d - %v", report.Code, report.Message)
    } else {
        report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    makeLogEntry(c).Error(report.Message)
    c.HTML(report.Code, report.Message.(string))
}