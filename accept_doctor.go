package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	expectedToken = "7a4f891b2e613dca"
	updateURL     = "http://127.0.0.1:8000/trauma/update_async/"
)

type TraumaResult struct {
	ID     string `json:"id"`
	Result string `json:"result"`
	Token  string `json:"token"`
}

func main() {
	http.HandleFunc("/trauma", handleProcess)
	fmt.Println("Server running at port :8088")
	http.ListenAndServe(":8088", nil)
}

func handleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при декодировании JSON: %s", err), http.StatusBadRequest)
		return
	}

	fmt.Println(requestBody.ID)

	w.WriteHeader(http.StatusOK)

	go func() {
		delay := 5
		time.Sleep(time.Duration(delay) * time.Second)

		result := "Проверка врача прошла успешно 1"
		if rand.Intn(2) == 0 {
			result = "Проверка врача прошла неуспешно 0"
		}

		// Отправка результата на другой сервер
		traumaRes := TraumaResult{
			ID:     strconv.Itoa(requestBody.ID),
			Result: result,
			Token:  expectedToken,
		}

		fmt.Println("json", traumaRes)
		jsonValue, err := json.Marshal(traumaRes)
		if err != nil {
			fmt.Println("Ошибка при маршализации JSON:", err)
			return
		}

		req, err := http.NewRequest(http.MethodPut, updateURL, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println("Ошибка при создании запроса на обновление:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса на обновление:", err)
			return
		}
		defer resp.Body.Close()

		fmt.Println("Ответ от сервера обновления:", resp.Status)
	}()
}
