package mellivora

import (
	"testing"
)

func TestDownloader_Next(t *testing.T) {
	// req, _ := http.NewRequest(http.MethodGet, "https://baidu.com", nil)
	// resp := &http.Response{
	// 	StatusCode: http.StatusOK,
	// 	Body:       ioutil.NopCloser(bytes.NewBufferString("hello")),
	// }
	// c := NewContext(nil, req, nil)
	// downloader := NewDownloader()
	// ctl := gomock.NewController(t)
	// t.Run("下载成功", func(t *testing.T) {
	// 	mr := roundtripper.NewMockRoundTripper(ctl)
	// 	mr.EXPECT().RoundTrip(gomock.Any()).Return(resp, nil)
	// 	c.SetHTTPClient(&http.Client{Transport: mr})
	// 	err := downloader.Next(func(c *Context) error {
	// 		assert.NotEqual(t, c.Response, nil)
	// 		assert.Equal(t, c.Response.StatusCode, http.StatusOK)
	// 		str, err := c.String()
	// 		assert.Equal(t, str, "hello")
	// 		assert.Equal(t, err, nil)
	// 		return nil
	// 	})(c)
	// 	assert.Equal(t, err, nil)
	// })

	// t.Run("Response有数据", func(t *testing.T) {
	// 	c.SetResponse(NewResponse(nil))
	// 	mr := roundtripper.NewMockRoundTripper(ctl)
	// 	c.SetHTTPClient(&http.Client{Transport: mr})
	// 	err := downloader.Next(func(c *Context) error {
	// 		t.Errorf("Response有数据进入了后续流程")
	// 		return nil
	// 	})(c)
	// 	assert.Equal(t, err, nil)
	// })

	// t.Run("下载失败", func(t *testing.T) {
	// 	c.SetResponse(nil)
	// 	mr := roundtripper.NewMockRoundTripper(ctl)
	// 	mr.EXPECT().RoundTrip(gomock.Any()).Return(nil, net.ErrClosed)
	// 	c.SetHTTPClient(&http.Client{Transport: mr})
	// 	err := downloader.Next(func(c *Context) error {
	// 		t.Errorf("下载失败进入了后续流程")
	// 		return nil
	// 	})(c)
	// 	assert.NotEqual(t, err, nil)
	// })
}
