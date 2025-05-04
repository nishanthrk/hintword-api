package sms_provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"hintword.com/api/app/configs"
)

type Msg91 struct {
	singletonInstance *Msg91

	BaseUrl string

	AuthKey string
	// Mutex for thread safety
	mu sync.Mutex
}

func (m *Msg91) GetInstance() *Msg91 {

	// Check if the instance already exists
	if m.singletonInstance == nil {
		// Use a mutex to ensure thread safety during instance creation
		m.mu.Lock()
		defer m.mu.Unlock()

		// Check again inside the critical section to prevent race conditions
		if m.singletonInstance == nil {
			// Create the singleton instance
			m.singletonInstance = m
		}
	}

	return m.singletonInstance
}

type msg91Response struct {
	RequestId string `json:"request_id"`
	Type      string `json:"type"`
	Message   string `json:"message"`
}

func (m *Msg91) SendOtp(templateId string, mobile string, payload interface{}, spoof bool) (err error) {
	if !configs.GetConfig().IsProd() || spoof {
		fmt.Printf("message might be spoofed for %v\n", mobile)
		return nil
	}
	response := msg91Response{}
	url := fmt.Sprintf("%s/otp?template_id=%s&mobile=%s", m.BaseUrl, templateId, mobile)
	client := &http.Client{}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("authkey", m.AuthKey)
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Type == "error" {
		return errors.New(response.Message)
	}
	return nil
}

func (m *Msg91) VerifyOtp(otp string, mobile string, spoof bool) (err error) {
	if !configs.GetConfig().IsProd() || spoof {
		fmt.Printf("message might be spoofed for %v\n", mobile)
		if time.Now().Format("2006") == otp {
			return nil
		}

		return errors.New(fmt.Sprintf("invalid otp for %s",
			configs.GetConfig().Env))
	}
	response := msg91Response{}
	url := fmt.Sprintf("%s/otp/verify?otp=%s&mobile=%s", m.BaseUrl, otp, mobile)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("authkey", m.AuthKey)
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Type == "error" {
		return errors.New(response.Message)
	}
	return nil
}

func (m *Msg91) ResendOtp(retryType string, mobile string, payload interface{}, spoof bool) (err error) {
	if !configs.GetConfig().IsProd() || spoof {
		fmt.Printf("message might be spoofed for %v\n", mobile)
		return nil
	}
	response := msg91Response{}
	url := fmt.Sprintf("%s/otp/retry?retrytype=%s&mobile=%s", m.BaseUrl, retryType, mobile)
	client := &http.Client{}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("authkey", m.AuthKey)
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Type == "error" {
		return errors.New(response.Message)
	}
	return nil
}
