package mellivora_test

import (
	"testing"

	"github.com/open-mellivora/mellivora"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"

	core_test "github.com/open-mellivora/mellivora/testing/core"
)

func TestEngine_Run(t *testing.T) {
	var result []string
	c := mellivora.NewEngine(32)
	ctl := gomock.NewController(t)
	ms := core_test.NewMockSpider(ctl)
	task, _ := mellivora.Get("https://baidu.com", func(c *mellivora.Context) mellivora.Task {
		result = append(result, "4")
		return nil
	})
	ms.EXPECT().StartRequests().Return(task)

	c.Use(
		mellivora.NewMiddleware(func(handleFunc mellivora.MiddlewareFunc) mellivora.MiddlewareFunc {
			return func(c *mellivora.Context) error {
				result = append(result, "1")
				err := handleFunc(c)
				result = append(result, "2")
				c.SetResponse(mellivora.NewResponse(nil))
				return err
			}
		}),
		mellivora.NewMiddleware(func(handleFunc mellivora.MiddlewareFunc) mellivora.MiddlewareFunc {
			return func(c *mellivora.Context) error {
				result = append(result, "3")
				return nil
			}
		}),
	)

	c.Run(ms)
	assert.Equal(t, result, []string{"1", "3", "2", "4"})
}
