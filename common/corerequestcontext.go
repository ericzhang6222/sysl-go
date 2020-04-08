package common

import (
	"context"
	"net/http"

	"github.com/anz-bank/sysl-go/common/internal"

	"github.com/anz-bank/pkg/log"
)

func NewLoggingRoundTripper(name string, base http.RoundTripper) http.RoundTripper {
	// temporary pass-through to get the real round tripper from the request context
	return &tempRoundtripper{name, base}
}

type coreRequestContext struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

type reqHeaderContext struct {
	header http.Header
}

type respHeaderAndStatusContext struct {
	header http.Header
	status int
}

// RequestHeaderToContext create a new context containing the request header
func RequestHeaderToContext(ctx context.Context, header http.Header) context.Context {
	return context.WithValue(ctx, reqHeaderContextKey{}, &reqHeaderContext{header})
}

// RequestHeaderFromContext retrieve the request header from the context
func RequestHeaderFromContext(ctx context.Context) http.Header {
	reqHeader := getReqHeaderContext(ctx)

	if reqHeader == nil {
		return nil
	}
	return reqHeader.header
}

// RespHeaderAndStatusToContext create a new context containing the response header and status
func RespHeaderAndStatusToContext(ctx context.Context, header http.Header, status int) context.Context {
	return context.WithValue(ctx, respHeaderAndStatusContextKey{}, &respHeaderAndStatusContext{header, status})
}

// RespHeaderAndStatusFromContext retrieve response header and status from the context
func RespHeaderAndStatusFromContext(ctx context.Context) (header http.Header, status int) {
	respHeaderAndStatus := getRespHeaderAndStatusContext(ctx)

	if respHeaderAndStatus == nil {
		return nil, http.StatusOK
	}

	header = respHeaderAndStatus.header
	status = respHeaderAndStatus.status
	return
}

func UpdateResponseStatus(ctx context.Context, status int) error {
	respHeaderAndStatus := getRespHeaderAndStatusContext(ctx)

	if respHeaderAndStatus == nil {
		return CreateError(ctx, InternalError, "response status not in context", nil)
	}
	respHeaderAndStatus.status = status
	return nil
}

func CoreRequestContextMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := log.With(traceIDLogField, GetTraceIDFromContext(r.Context())).Onto(r.Context())

			ctx = internal.AddResponseBodyMonitorToContext(ctx)
			defer internal.CheckForUnclosedResponses(ctx)

			reqLogger, entry := internal.NewRequestLogger(entry, r)
			w = reqLogger.ResponseWriter(w)
			defer reqLogger.FlushLog()
			ctx = LoggerToContext(ctx, logger, entry)

			r = r.WithContext(ctx)

			tl := internal.NewRequestTimer(w, r)
			w = tl.RespWrapper
			defer tl.Log(entry)

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

type coreRequestContextKey struct{}

func getCoreContext(ctx context.Context) *coreRequestContext {
	return ctx.Value(coreRequestContextKey{}).(*coreRequestContext)
}

type reqHeaderContextKey struct{}
type respHeaderAndStatusContextKey struct{}

func getReqHeaderContext(ctx context.Context) *reqHeaderContext {
	if ctx.Value(reqHeaderContextKey{}) == nil {
		return nil
	}
	return ctx.Value(reqHeaderContextKey{}).(*reqHeaderContext)
}

func getRespHeaderAndStatusContext(ctx context.Context) *respHeaderAndStatusContext {
	if ctx.Value(respHeaderAndStatusContextKey{}) == nil {
		return nil
	}
	return ctx.Value(respHeaderAndStatusContextKey{}).(*respHeaderAndStatusContext)
}

type tempRoundtripper struct {
	name string
	base http.RoundTripper
}

func (t *tempRoundtripper) RoundTrip(r *http.Request) (*http.Response, error) {
	logentry := GetLogEntryFromContext(r.Context()).WithField("Downsteam", t.name)
	return internal.NewLoggingRoundTripper(logentry, t.base).RoundTrip(r)
}
