package dal

import "context"

type HealthSeed struct {
	Name   string
	Status string
}

type SystemDAL struct{}

func NewSystemDAL() *SystemDAL {
	return &SystemDAL{}
}

func (d *SystemDAL) GetHealth(_ context.Context) (*HealthSeed, error) {
	return &HealthSeed{
		Name:   "clear-bill-server",
		Status: "ok",
	}, nil
}
