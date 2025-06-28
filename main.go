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
var initMutex sync.Mutex

func init() {
	log.Println("[DEBUG] Lambda init: 開始")
	_ = godotenv.Load()
	log.Println("[DEBUG] Lambda init: 完了")
}

// resetInitialization 初期化状態をリセット
func resetInitialization() {
	initMutex.Lock()
	defer initMutex.Unlock()

	isInitialized = false
	ginLambda = nil
	// sync.Onceをリセットするために新しいインスタンスを作成
	initOnce = sync.Once{}
}

// initializeApp アプリケーションの初期化を実行
func initializeApp() {
	initMutex.Lock()
	defer initMutex.Unlock()

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

			log.Println("[DEBUG] Lambda initializeApp: infra.Initialize() 開始")
			infra.Initialize()
			log.Println("[DEBUG] Lambda initializeApp: infra.Initialize() 完了")

			log.Println("[DEBUG] Lambda initializeApp: infra.SetupDB() 開始")
			db := infra.SetupDB()
			log.Println("[DEBUG] Lambda initializeApp: infra.SetupDB() 完了")

			log.Println("[DEBUG] Lambda initializeApp: repositories.NewDiaryRepository 開始")
			diaryRepository := repositories.NewDiaryRepository(db)
			log.Println("[DEBUG] Lambda initializeApp: repositories.NewDiaryRepository 完了")

			log.Println("[DEBUG] Lambda initializeApp: usecases.NewDiaryUsecase 開始")
			diaryUsecase := usecases.NewDiaryUsecase(diaryRepository)
			log.Println("[DEBUG] Lambda initializeApp: usecases.NewDiaryUsecase 完了")

			log.Println("[DEBUG] Lambda initializeApp: controllers.NewDiaryController 開始")
			diaryController := controllers.NewDiaryController(diaryUsecase)
			log.Println("[DEBUG] Lambda initializeApp: controllers.NewDiaryController 完了")

			log.Println("[DEBUG] Lambda initializeApp: usecases.NewDiaryAnalysisUsecase 開始")
			diaryAnalysisUsecase := usecases.NewDiaryAnalysisUsecase(diaryRepository)
			log.Println("[DEBUG] Lambda initializeApp: usecases.NewDiaryAnalysisUsecase 完了")

			log.Println("[DEBUG] Lambda initializeApp: controllers.NewDiaryAnalysisController 開始")
			diaryAnalysisController := controllers.NewDiaryAnalysisController(diaryAnalysisUsecase)
			log.Println("[DEBUG] Lambda initializeApp: controllers.NewDiaryAnalysisController 完了")

			log.Println("[DEBUG] Lambda initializeApp: gin.Default() 開始")
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

			log.Println("[DEBUG] Lambda initializeApp: gin.Default() 完了")

			log.Println("[DEBUG] Lambda initializeApp: routes.SetupCORS 開始")
			routes.SetupCORS(router)
			log.Println("[DEBUG] Lambda initializeApp: routes.SetupCORS 完了")

			log.Println("[DEBUG] Lambda initializeApp: routes.SetupSwaggerEndpoints 開始")
			routes.SetupSwaggerEndpoints(router)
			log.Println("[DEBUG] Lambda initializeApp: routes.SetupSwaggerEndpoints 完了")

			log.Println("[DEBUG] Lambda initializeApp: routes.SetupAPIEndpoints 開始")
			routes.SetupAPIEndpoints(router, diaryController, diaryAnalysisController)
			log.Println("[DEBUG] Lambda initializeApp: routes.SetupAPIEndpoints 完了")

			// ヘルスチェックとデバッグ情報を提供するエンドポイント
			router.GET("/health", func(c *gin.Context) {
				log.Printf("[DEBUG] Health check endpoint called")
				c.JSON(200, gin.H{"status": "healthy"})
			})

			log.Println("[DEBUG] Lambda initializeApp: ginadapter.New(router) 開始")
			ginLambda = ginadapter.New(router)
			log.Println("[DEBUG] Lambda initializeApp: ginadapter.New(router) 完了")

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
	log.Println("[DEBUG] Lambda main: lambda.Start 開始")
	lambda.Start(Handler)
}
