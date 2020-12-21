package tracing

import (
	"net/http"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/zalando/skipper/filters/filtertest"
	"github.com/zalando/skipper/tracing/tracingtest"
)

func TestTracingTag(t *testing.T) {
	for _, ti := range []struct {
		name     string
		value    string
		context  *filtertest.Context
		expected interface{}
	}{{
		"plain key value",
		"test_value",
		&filtertest.Context{
			FRequest: &http.Request{},
		},
		"test_value",
	}, {
		"tag from header",
		"${request.header.X-Flow-Id}",
		&filtertest.Context{
			FRequest: &http.Request{
				Header: http.Header{
					"X-Flow-Id": []string{"foo"},
				},
			},
		},
		"foo",
	}, {
		"tag from missing header",
		"${request.header.missing}",
		&filtertest.Context{
			FRequest: &http.Request{},
		},
		nil,
	},
	} {
		t.Run(ti.name, func(t *testing.T) {
			span := tracingtest.NewSpan("proxy")
			ti.context.FRequest = ti.context.FRequest.WithContext(opentracing.ContextWithSpan(ti.context.FRequest.Context(), span))

			s := NewTag()
			f, err := s.CreateFilter([]interface{}{"test_tag", ti.value})
			if err != nil {
				t.Fatal(err)
			}

			f.Request(ti.context)

			got := span.Tags["test_tag"]
			if got != ti.expected {
				t.Errorf("unexpected tag value '%v' != '%v'", got, ti.expected)
			}
		})
	}
}
