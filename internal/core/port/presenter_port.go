package port

type Presenter interface {
	Present(any) ([]byte, error)
}