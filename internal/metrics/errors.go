package metrics

import "fmt"

type metricNotFoundError struct {
	msg string
}

func MetricNotFoundError(name string) *metricNotFoundError {
	return &metricNotFoundError{
		msg: fmt.Sprintf("metric %q not found", name),
	}
}

func (e *metricNotFoundError) Error() string {
	return e.msg
}

type metricFailedToRegister struct {
	msg string
}

func MetricFailedToRegisterError(name string) *metricFailedToRegister {
	return &metricFailedToRegister{
		msg: fmt.Sprintf("failed to register metric %q", name),
	}
}

func (e *metricFailedToRegister) Error() string {
	return e.msg
}
