package rest

type RestError struct {
	Reason string
	Err    error
}

func (this *RestError) Error() string {
	return this.Reason
}
