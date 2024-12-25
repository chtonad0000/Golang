//go:build !solution

package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func MakeHandler(service interface{}) http.Handler {
	mux := http.NewServeMux()
	serviceType := reflect.TypeOf(service)
	serviceValue := reflect.ValueOf(service)

	for i := 0; i < serviceType.NumMethod(); i++ {
		method := serviceType.Method(i)

		mux.HandleFunc("/"+method.Name, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			reqType := method.Type.In(2).Elem()
			reqValue := reflect.New(reqType)

			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(reqValue.Interface()); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			ctx := r.Context()
			result := method.Func.Call([]reflect.Value{
				serviceValue,
				reflect.ValueOf(ctx),
				reqValue,
			})

			if !result[1].IsNil() {
				http.Error(w, result[1].Interface().(error).Error(), http.StatusInternalServerError)
				return
			}

			respValue := result[0].Interface()
			w.Header().Set("Content-Type", "application/json")
			encoder := json.NewEncoder(w)
			if err := encoder.Encode(respValue); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		})
	}

	return mux
}

func Call(ctx context.Context, endpoint string, method string, req, rsp interface{}) error {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to serialize request: %w", err)
	}

	url := fmt.Sprintf("%s/%s", endpoint, method)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		errClose := Body.Close()
		if errClose != nil {
			return
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("server returned error: %s", string(body))
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(rsp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
