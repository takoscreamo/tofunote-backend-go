package main

import (
	"context"
	"emotra-backend/api/controllers"
	"emotra-backend/infra"
	"emotra-backend/repositories"
	"emotra-backend/routes"
	"emotra-backend/usecases"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var ginLambda *ginadapter.GinLambda

func init() {
	// .envファイルがなくてもエラーを無視
	_ = godotenv.Load()

	// データベース初期化
	infra.Initialize()
	db := infra.SetupDB()

	// 依存関係の設定
	diaryRepository := repositories.NewDiaryRepository(db)
	diaryUsecase := usecases.NewDiaryUsecase(diaryRepository)
	diaryController := controllers.NewDiaryController(diaryUsecase)

	diaryAnalysisUsecase := usecases.NewDiaryAnalysisUsecase(diaryRepository)
	diaryAnalysisController := controllers.NewDiaryAnalysisController(diaryAnalysisUsecase)

	// Ginルーターの設定
	router := gin.Default()

	// CORS設定を追加
	routes.SetupCORS(router)

	// SwaggerUIとOpenAPI仕様書のエンドポイントを設定
	routes.SetupSwaggerEndpoints(router)

	// APIエンドポイントを設定
	routes.SetupAPIEndpoints(router, diaryController, diaryAnalysisController)

	// Lambda用のアダプターを作成
	ginLambda = ginadapter.New(router)
}

// Handler AWS Lambdaのハンドラー関数
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received request: %s %s", req.HTTPMethod, req.Path)

	// Ginアダプターを使用してリクエストを処理
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
