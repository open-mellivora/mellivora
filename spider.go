package mellivora

type Spider interface {
	// StartRequests generate first Task
	StartRequests() Task
}
