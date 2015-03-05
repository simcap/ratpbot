package services

type Service interface {
	Get() (string, error)
}

type GlobalTraffic struct{}

func (s GlobalTraffic) Get() (string, error) {
	return "Le service est globalement bon!", nil
}
