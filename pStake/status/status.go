package status

import "encoding/json"

const (
	COSMOS   = "COSMOS"
	ETHEREUM = "ETHEREUM"
)

type Status struct {
	Name            string
	LastCheckHeight int64
}

func GetCosmosStatus() (Status, error) {
	var status Status
	data, err := db.Get([]byte(COSMOS), nil)
	if err != nil {
		return status, err
	}
	err = json.Unmarshal(data, &status)
	if err != nil {
		return status, err
	}
	return status, nil
}

func SetCosmosStatus(height int64) error {
	status := Status{
		Name:            COSMOS,
		LastCheckHeight: height,
	}
	b, err := json.Marshal(status)
	if err != nil {
		return err
	}
	err = db.Put([]byte(status.Name), b, nil)
	if err != nil {
		return err
	}
	return nil
}

func GetEthereumStatus() (Status, error) {
	var status Status
	data, err := db.Get([]byte(ETHEREUM), nil)
	if err != nil {
		return status, err
	}
	err = json.Unmarshal(data, &status)
	if err != nil {
		return status, err
	}
	return status, nil
}

func SetEthereumStatus(height int64) error {
	status := Status{
		Name:            ETHEREUM,
		LastCheckHeight: height,
	}
	b, err := json.Marshal(status)
	if err != nil {
		return err
	}
	err = db.Put([]byte(status.Name), b, nil)
	if err != nil {
		return err
	}
	return nil
}
