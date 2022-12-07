package constant

const (
	OrderStatusPending  string = "P"
	OrderStatusAccepted string = "A"
	OrderStatusRejected string = "R"
)

var OrderStatuses = []string{OrderStatusAccepted, OrderStatusPending, OrderStatusRejected}
