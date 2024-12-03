package main

import (
    "net/http"
    "net/http/httptest"
    "strconv"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

var cafeList = map[string][]string{
    "moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
    countStr := req.URL.Query().Get("count")
    if countStr == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("count missing"))
        return
    }

    count, err := strconv.Atoi(countStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("wrong count value"))
        return
    }

    city := req.URL.Query().Get("city")

    cafe, ok := cafeList[city]
    if !ok {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("wrong city value"))
        return
    }

    if count > len(cafe) {
        count = len(cafe)
    }

    answer := strings.Join(cafe[:count], ",")

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(answer))
}

// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
    // Создаем запрос, где count больше, чем доступное количество кафе
    req, err := http.NewRequest("GET", "/cafe?count=10&city=moscow", nil)
    assert.NoError(t, err)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)

    // Проверяем код ответа и содержание тела
    assert.Equal(t, http.StatusOK, responseRecorder.Code)
    assert.Equal(t, "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент", responseRecorder.Body.String())
}

// Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое
func TestMainHandlerWhenOk(t *testing.T) {
    // Создаем запрос с корректными параметрами
    req, err := http.NewRequest("GET", "/cafe?count=2&city=moscow", nil)
    assert.NoError(t, err)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)

    // Проверяем код ответа и содержание тела
    assert.Equal(t, http.StatusOK, responseRecorder.Code)
    assert.NotEmpty(t, responseRecorder.Body.String())
    assert.Equal(t, "Мир кофе,Сладкоежка", responseRecorder.Body.String())
}

// Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа
func TestMainHandlerWhenWrongCityValue(t *testing.T) {
    // Создаем запрос с неподдерживаемым городом
    req, err := http.NewRequest("GET", "/cafe?count=2&city=tula", nil)
    assert.NoError(t, err)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)

     // Проверяем код ответа и сообщение об ошибке
     assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
     assert.Equal(t, "wrong city value", responseRecorder.Body.String())
} 