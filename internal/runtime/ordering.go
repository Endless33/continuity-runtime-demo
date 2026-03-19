package runtime

type OrderingMode string

const (
	OrderingStrict  OrderingMode = "strict"
	OrderingPartial OrderingMode = "partial"
)

type OrderingPolicy struct {
	Mode       OrderingMode
	WindowSize int
}

func NewOrderingPolicy() *OrderingPolicy {
	return &OrderingPolicy{
		Mode:       OrderingPartial,
		WindowSize: 3,
	}
}

func (o *OrderingPolicy) AllowOutOfOrder(expected, got int) bool {
	if o.Mode == OrderingStrict {
		return got == expected
	}

	if got >= expected && got <= expected+o.WindowSize {
		return true
	}

	return false
}