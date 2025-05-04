package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func safeString(val interface{}) string {
	if val == nil {
	 return ""
	}
	return fmt.Sprintf("%v", val)
   }

type LogHandler struct {}

func (l *LogHandler) AddLog(c *gin.Context) {
	shopID := c.Param("shopID")

	if shopID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не указан идентификатор магазина"})
		return
	}

	filePath := filepath.Join(FILE_PATH, shopID+".log")

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось прочитать тело запроса"})
	return
	}

	logEntry := map[string]interface{}{
		"time":    time.Now().Format(time.RFC3339),
		"method":  c.Request.Method,
		"path":    c.Request.URL.Path,
		"headers": c.Request.Header,
		"body":    string(bodyBytes),
	}

	logBytes, err := json.Marshal(logEntry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сериализовать лог"})
		return
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось открыть файл"})
		return
	}
	defer f.Close()

	if _, err := f.Write(append(logBytes, '\n')); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось записать в файл"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name": "Тестовая интеграция",
		"time": time.Now().UTC().Format(time.RFC3339),
		"version": "12312",
	})
}

func (l *LogHandler) RemoveLogs(c *gin.Context) {
	shopID := c.Param("shopID")
	filePath := filepath.Join(FILE_PATH, shopID+".log")

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "лог уже отсутствует"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось удалить файл"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "лог удалён"})
}

func (l *LogHandler) GetLogs(c *gin.Context) {
	shopID := c.Param("shopID")
	filePath := filepath.Join(FILE_PATH, shopID+".log")

	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.String(http.StatusNotFound, "Лог не найден")
		} else {
			c.String(http.StatusInternalServerError, "Не удалось прочитать файл")
		}
		return
	}

	// Разбиваем файл по строкам
	lines := bytes.Split(content, []byte("\n"))
	var logs []map[string]interface{}

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var logEntry map[string]interface{}
		if err := json.Unmarshal(line, &logEntry); err != nil {
			c.String(http.StatusInternalServerError, "Не удалось распарсить лог")
			return
		}
		logs = append(logs, logEntry)
	}

	// Рендерим простую HTML-страницу
	html := "<html><head><title>Логи магазина " + shopID + "</title>"
	html += `<style>
	body { font-family: Arial, sans-serif; padding: 20px; }
	table { width: 100%; border-collapse: collapse; }
	th, td { border: 1px solid #ccc; padding: 8px; text-align: left; }
	th { background-color: #f5f5f5; }
	</style>`
	html += "</head><body>"
	html += "<h1>Логи магазина " + shopID + "</h1>"
	html += "<table><thead><tr><th>Время</th><th>Метод</th><th>Путь</th><th>Тело запроса</th></tr></thead><tbody>"

	for _, log := range logs {
		timeVal := safeString(log["time"])
		methodVal := safeString(log["method"])
		pathVal := safeString(log["path"])
		bodyVal := safeString(log["body"])

		html += "<tr>"
		html += "<td>" + timeVal + "</td>"
		html += "<td>" + methodVal + "</td>"
		html += "<td>" + pathVal + "</td>"
		html += "<td><pre style='white-space: pre-wrap;'>" + bodyVal + "</pre></td>"
		html += "</tr>"
	}

	html += "</tbody></table>"
	html += "</body></html>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func (l *LogHandler) GetExternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"status": "500", "message": "Внутренняя ошибка сервера"})
}

func (l *LogHandler) SlowLogs(c *gin.Context) {
	time.Sleep(12 * time.Second)
		
	c.JSON(http.StatusOK, gin.H{
		"name": "Тестовая интеграция",
		"time": time.Now().UTC().Format(time.RFC3339),
		"version": "12312",
		"text": "Этот ответ занял 12 секунд",
	})
}