package rest

type RestError struct {
	Reason string
	Err    error
}

func (this *RestError) Error() string {
	return this.Reason
}

//type ParamError struct {
//	Name   string
//	Reason string
//}

//func (this *ParamError) Error() string {
//	return this.Name + ":" + this.Reason
//}
