package main

import (
	"context"
	"emotra-backend/api/controllers"
	"emotra-backend/infra"
	"emotra-backend/repositories"
	"emotra-backend/routes"
	"emotra-backend/usecases"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var ginLambda *ginadapter.GinLambda
var isInitialized bool
var initOnce sync.Once

func init() {
	_ = godotenv.Load()
}

// initializeApp アプリケーションの初期化を実行
func initializeApp() {
	initOnce.Do(func() {
		log.Println("[DEBUG] Lambda initializeApp: 開始")

		// 初期化タイムアウトの設定
		done := make(chan bool, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[ERROR] Lambda initializeApp: Panic occurred: %v", r)
					isInitialized = false
				}
				done <- true
			}()

			infra.Initialize()
			db := infra.SetupDB()

			diaryRepository := repositories.NewDiaryRepository(db)
			diaryUsecase := usecases.NewDiaryUsecase(diaryRepository)
			diaryController := controllers.NewDiaryController(diaryUsecase)

			diaryAnalysisUsecase := usecases.NewDiaryAnalysisUsecase(diaryRepository)
			diaryAnalysisController := controllers.NewDiaryAnalysisController(diaryAnalysisUsecase)

			router := gin.Default()

			// カスタムログミドルウェアを追加
			router.Use(func(c *gin.Context) {
				log.Printf("[DEBUG] Gin: %s %s", c.Request.Method, c.Request.URL.Path)
				c.Next()
				log.Printf("[DEBUG] Gin: Response Status: %d", c.Writer.Status())
				if c.Writer.Status() >= 400 {
					log.Printf("[DEBUG] Gin: Error occurred for %s %s", c.Request.Method, c.Request.URL.Path)
				}
			})

			routes.SetupCORS(router)
			routes.SetupSwaggerEndpoints(router)
			routes.SetupAPIEndpoints(router, diaryController, diaryAnalysisController)

			// ヘルスチェックエンドポイント
			router.GET("/health", func(c *gin.Context) {
				c.JSON(200, gin.H{"status": "healthy"})
			})

			ginLambda = ginadapter.New(router)
			isInitialized = true
			log.Println("[DEBUG] Lambda initializeApp: 完了")
		}()

		// 30秒でタイムアウト
		select {
		case <-done:
			log.Println("[DEBUG] Lambda initializeApp: 正常完了")
		case <-time.After(30 * time.Second):
			log.Printf("[ERROR] Lambda initializeApp: タイムアウト")
			isInitialized = false
		}
	})
}

// Handler AWS Lambdaのハンドラー関数
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("[DEBUG] Lambda Handler: Received request: %s %s", req.HTTPMethod, req.Path)

	// 遅延初期化を実行
	initializeApp()

	// 初期化が完了していない場合はエラーレスポンスを返す
	if !isInitialized {
		log.Printf("[ERROR] Lambda Handler: Initialization not completed yet")
		return events.APIGatewayProxyResponse{
			StatusCode: 503,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error": "Service is initializing, please try again"}`,
		}, nil
	}

	// Ginアダプターを使用してリクエストを処理
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
