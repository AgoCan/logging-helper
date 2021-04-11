package service

type Health struct{}

func (h Health) Status() string {
	return "Working!"
}
