package helpers

import "fmt"

type WebResponse struct {
	Code int
	Status string
	Data interface{}
}

func (wr WebResponse) Error() string {
	return fmt.Sprintf("(%v) %v, %v ", wr.Code, wr.Status, wr.Data)
}