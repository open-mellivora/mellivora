//go:generate go run -mod=mod github.com/golang/mock/mockgen -package roundtripper -destination=testing/roundtripper/roundtripper.go net/http RoundTripper
//go:generate go run -mod=mod github.com/golang/mock/mockgen -package core -destination=testing/core/spider.go -source=core/spider.go Spider

package mellivora
