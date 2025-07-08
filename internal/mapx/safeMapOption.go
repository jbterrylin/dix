package mapx

type safeMapOption struct {
	maxDeleteCount int64
}

type SafeMapOption func(*safeMapOption)

func WithMaxDeleteCount(maxDeleteCount int64) SafeMapOption {
	return func(o *safeMapOption) {
		o.maxDeleteCount = maxDeleteCount
	}
}
