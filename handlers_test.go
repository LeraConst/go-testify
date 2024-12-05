package main

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
    req, err := http.NewRequest("GET", "/cafe?count=10&city=moscow", nil)
    assert.NoError(t, err)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)

    assert.Equal(t, http.StatusOK, responseRecorder.Code)
    assert.Equal(t, "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент", responseRecorder.Body.String())
}

// Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое
func TestMainHandlerWhenOk(t *testing.T) {
    req, err := http.NewRequest("GET", "/cafe?count=2&city=moscow", nil)
    assert.NoError(t, err)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)

    assert.Equal(t, http.StatusOK, responseRecorder.Code)
    assert.NotEmpty(t, responseRecorder.Body.String())
    assert.Equal(t, "Мир кофе,Сладкоежка", responseRecorder.Body.String())
}

// Город, который передаётся в параметре city, не поддерживается
func TestMainHandlerWhenWrongCityValue(t *testing.T) {
    req, err := http.NewRequest("GET", "/cafe?count=2&city=tula", nil)
    assert.NoError(t, err)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)

    assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
    assert.Equal(t, "wrong city value", responseRecorder.Body.String())
}
