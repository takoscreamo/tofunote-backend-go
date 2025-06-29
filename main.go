package main

import (
	"context"
	"feelog-backend/api/controllers"
	"feelog-backend/infra"
	"feelog-backend/repositories"
	"feelog-backend/routes"
	"feelog-backend/usecases"
	"log"
	"os"
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

			log.Println("[DEBUG] Lambda initializeApp: ginadapter.New(router) 開始")
			ginLambda = ginadapter.New(router)
			log.Println("[DEBUG] Lambda initializeApp: ginadapter.New(router) 完了")

			isInitialized = true
			log.Println("[DEBUG] Lambda initializeApp: 完了")
		}()

		// 60秒でタイムアウト（本番環境では長めに設定）
		select {
		case <-done:
			log.Println("[DEBUG] Lambda initializeApp: 正常完了")
		case <-time.After(60 * time.Second):
			log.Printf("[ERROR] Lambda initializeApp: タイムアウト")
			isInitialized = false
		}
	})
}

// Handler AWS Lambdaのハンドラー関数
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("[DEBUG] Lambda Handler: Received request: %s %s", req.HTTPMethod, req.Path)

	// pingエンドポイントの特別なデバッグログ
	if req.Path == "/ping" {
		log.Printf("[DEBUG] Lambda Handler: Ping endpoint detected - Path: %s, Method: %s", req.Path, req.HTTPMethod)
	}

	// 環境変数の確認（パスワードは隠す）
	log.Printf("[DEBUG] Lambda Handler: ENV=%s, DB_HOST=%s, DB_USER=%s, DB_NAME=%s, DB_PORT=%s",
		os.Getenv("ENV"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

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
			Body: `{"error": "Service is initializing, please try again", "path": "` + req.Path + `"}`,
		}, nil
	}

	// Ginアダプターを使用してリクエストを処理
	response, err := ginLambda.ProxyWithContext(ctx, req)
	if err != nil {
		log.Printf("[ERROR] Lambda Handler: Proxy error for %s %s: %v", req.HTTPMethod, req.Path, err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error": "Internal server error", "path": "` + req.Path + `", "details": "` + err.Error() + `"}`,
		}, nil
	}

	// pingエンドポイントのレスポンスデバッグログ
	if req.Path == "/ping" {
		log.Printf("[DEBUG] Lambda Handler: Ping response - Status: %d, Body: %s", response.StatusCode, response.Body)
	}

	log.Printf("[DEBUG] Lambda Handler: Response status: %d for %s %s", response.StatusCode, req.HTTPMethod, req.Path)
	return response, nil
}

func main() {
	log.Println("[DEBUG] Lambda main: lambda.Start 開始")
	lambda.Start(Handler)
}
