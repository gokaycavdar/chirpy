package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCleanedBody(t *testing.T) {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Temiz metin",
			input:    "I love Go",
			expected: "I love Go",
		},
		{
			name:     "Küfürlü metin",
			input:    "This is a kerfuffle",
			expected: "This is a ****",
		},
		{
			name:     "Büyük harfli küfür",
			input:    "I hate SHARBERT",
			expected: "I hate ****",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getCleanedBody(tt.input, badWords)
			if actual != tt.expected {
				t.Errorf("%s: beklenen %q, gelen %q", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestHandlerChirpsValidate(t *testing.T) {
	// Test için örnek bir JSON isteği
	body, _ := json.Marshal(map[string]string{
		"body": "I had a kerfuffle today",
	})

	// Gerçek bir request simülasyonu
	req := httptest.NewRequest("POST", "/api/validate_chirp", bytes.NewReader(body))
	// Yanıtı kaydedecek yapı (Response Recorder)
	rr := httptest.NewRecorder()

	// Handler'ı çalıştır
	handlerChirpsValidate(rr, req)

	// Durum kodunu kontrol et
	if rr.Code != http.StatusOK {
		t.Errorf("Yanlış durum kodu: beklenen 200, gelen %d", rr.Code)
	}

	// Yanıt içeriğini kontrol et
	var resBody struct {
		CleanedBody string `json:"cleaned_body"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resBody); err != nil {
		t.Fatalf("Yanıt decode edilemedi: %v", err)
	}

	expected := "I had a **** today"
	if resBody.CleanedBody != expected {
		t.Errorf("Yanlış temizlenmiş içerik: beklenen %q, gelen %q", expected, resBody.CleanedBody)
	}
}