package ionCannonClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hiring-seedtag/mihai-lupoiu-go-backend-test/internal/core/domain"
)

// IonCannonClient represents the Ion Cannon client.
// [Improvement] add timeout to http client.
type IonCannonClient struct {
	BaseURL string
}

// NewIonCannonClient creates a new instance of the IonCannonClient.
func NewIonCannonClient(baseURL string) *IonCannonClient {
	return &IonCannonClient{
		BaseURL: baseURL,
	}
}

// CheckStatus sends an HTTP GET request to check the Ion Cannon status.
func (c *IonCannonClient) CheckStatus() (*domain.IonCannon, error) {
	url := c.BaseURL + "/status"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to check Ion Cannon status: %s", resp.Status)
	}

	var status struct {
		Generation int  `json:"generation"`
		Available  bool `json:"available"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}

	res := &domain.IonCannon{
		Generation: status.Generation,
		Available:  status.Available,
	}

	return res, nil
}

// FireCommand sends an HTTP POST request to fire the Ion Cannon.
func (c *IonCannonClient) FireCommand(
	targetX int,
	targetY int,
	enemies int,
) (casualties int, generation int, err error) {
	url := c.BaseURL + "/fire"
	body := map[string]interface{}{
		"target": map[string]int{
			"x": targetX,
			"y": targetY,
		},
		"enemies": enemies,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return 0, 0, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("failed to fire Ion Cannon: %s", resp.Status)
	}

	var result struct {
		Casualties int `json:"casualties"`
		Generation int `json:"generation"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, err
	}

	return result.Casualties, result.Generation, nil
}
