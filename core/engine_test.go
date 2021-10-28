package core_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"

	"icode.baidu.com/baidu/goodcoder/wangyufeng04/core"
	core_test "icode.baidu.com/baidu/goodcoder/wangyufeng04/record/core"
)

func TestEngine_Run(t *testing.T) {
	c := core.NewEngine()
	ctl := gomock.NewController(t)
	ms := core_test.NewMockSpider(ctl)
	ms.EXPECT().StartRequests(gomock.Any())
	var result []string
	c.Use(
		core.NewMiddleware(func(handleFunc core.HandlerFunc) core.HandlerFunc {
			return func(c *core.Context) error {
				result = append(result, "1")
				err := handleFunc(c)
				result = append(result, "2")
				c.SetResponse(core.NewResponse(nil))
				return err
			}
		}),
		core.NewMiddleware(func(handleFunc core.HandlerFunc) core.HandlerFunc {
			return func(c *core.Context) error {
				result = append(result, "3")
				return nil
			}
		}),
	)

	_ = c.Get("https://baidu.com", func(c *core.Context) error {
		result = append(result, "4")
		return nil
	}, core.DontFilter())
	c.Run(ms)
	assert.Equal(t, result, []string{"1", "3", "2", "4"})
}
