package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	models "project/Model"
	service "project/Service"

	"github.com/iris-contrib/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

type GroqRequest struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type ChatHistory struct {
	Messages []Message `json:"messages"`
}

type ChatHistoryManager struct {
	histories *ChatHistory
}

func NewChatHistoryManager() *ChatHistoryManager {
	return &ChatHistoryManager{
		histories: &ChatHistory{},
	}
}

func (m *ChatHistoryManager) GetHistory() []Message {
	if m.histories != nil {
		return m.histories.Messages
	}
	return []Message{}
}

func (m *ChatHistoryManager) SaveHistory(messages []Message) {
	const maxHistorySize = 50
	if len(messages) > maxHistorySize {
		messages = messages[len(messages)-maxHistorySize:]
	}
	m.histories = &ChatHistory{Messages: messages}
}

func main() {
	app := iris.New()
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           86400,
	})
	app.UseRouter(crs)

	historyManager := NewChatHistoryManager()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Không thể tải file .env: ", err)
	}
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Println("Warning: You must set the GROQ_API_KEY environment variable to use this server.")
	}

	app.Post("/api/save", func(ctx iris.Context) {
		var requestData struct {
			Dialog models.Dialog `json:"dialog"`
			Words  []models.Word `json:"words"`
		}
		if err := ctx.ReadJSON(&requestData); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(map[string]string{"error": "Dữ liệu không hợp lệ"})
			return
		}

		if err := service.CreateData(requestData.Words, requestData.Dialog); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(map[string]string{"error": "Lỗi khi lưu dữ liệu"})
			return
		}

		ctx.JSON(map[string]string{"message": "Lưu thành công !!"})
	})

	app.Post("/ask", func(ctx iris.Context) {
		var requestData struct {
			Prompt string `json:"prompt"`
		}
		if err := ctx.ReadJSON(&requestData); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"error": "Invalid JSON: " + err.Error()})
			return
		}
		if requestData.Prompt == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"error": "Prompt không được để trống"})
			return
		}

		messages := historyManager.GetHistory()
		messages = append(messages, Message{Role: "user", Content: requestData.Prompt})

		groqReq := GroqRequest{
			Messages: messages,
			Model:    "mixtral-8x7b-32768",
			// Model:       "qwen-2.5-coder-32b",
			Temperature: 0,
		}
		jsonData, err := json.Marshal(groqReq)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Lỗi khi chuẩn bị dữ liệu request"})
			return
		}

		req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Lỗi khi tạo request đến Groq API"})
			return
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Printf("Error calling Groq API: %v", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Không thể kết nối tới Groq API"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response: %v", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Lỗi khi đọc phản hồi từ API"})
			return
		}

		var groqResp GroqResponse
		if err := json.Unmarshal(body, &groqResp); err != nil {
			log.Printf("Error unmarshaling response: %v", err)
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Lỗi khi phân tích phản hồi từ API"})
			return
		}

		if groqResp.Error.Message != "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"error": groqResp.Error.Message})
			return
		}
		if len(groqResp.Choices) == 0 {
			ctx.JSON(iris.Map{"error": "Không nhận được phản hồi từ API"})
			return
		}

		markdown := groqResp.Choices[0].Message.Content
		messages = append(messages, Message{Role: "assistant", Content: markdown})
		historyManager.SaveHistory(messages)

		unsafeHTML := blackfriday.Run([]byte(markdown))
		safeHTML := bluemonday.UGCPolicy().SanitizeBytes(unsafeHTML)

		ctx.JSON(iris.Map{
			"markdown": markdown,
			"html":     string(safeHTML),
		})
	})

	app.HandleDir("/", iris.Dir("./Frontend"))
	app.Listen(":8000")
}
