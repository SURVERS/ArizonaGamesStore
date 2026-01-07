package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

func VerifyRecaptcha(token string) (bool, float64, error) {
	secretKey := os.Getenv("RECAPTCHA_SECRET_KEY")
	if secretKey == "" {
		return true, 1.0, nil
	}

	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{
			"secret":   {secretKey},
			"response": {token},
		})

	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, 0, err
	}

	var result RecaptchaResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return false, 0, err
	}

	if !result.Success {
		return false, 0, fmt.Errorf("recaptcha verification failed: %v", result.ErrorCodes)
	}

	return true, result.Score, nil
}
