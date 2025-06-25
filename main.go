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
	log.Println("[DEBUG] Lambda init: 開始")
	_ = godotenv.Load()
	log.Println("[DEBUG] Lambda init: infra.Initialize() 開始")
	infra.Initialize()
	log.Println("[DEBUG] Lambda init: infra.Initialize() 完了")

	log.Println("[DEBUG] Lambda init: infra.SetupDB() 開始")
	db := infra.SetupDB()
	log.Println("[DEBUG] Lambda init: infra.SetupDB() 完了")

	log.Println("[DEBUG] Lambda init: repositories.NewDiaryRepository 開始")
	diaryRepository := repositories.NewDiaryRepository(db)
	log.Println("[DEBUG] Lambda init: usecases.NewDiaryUsecase 開始")
	diaryUsecase := usecases.NewDiaryUsecase(diaryRepository)
	log.Println("[DEBUG] Lambda init: controllers.NewDiaryController 開始")
	diaryController := controllers.NewDiaryController(diaryUsecase)

	log.Println("[DEBUG] Lambda init: usecases.NewDiaryAnalysisUsecase 開始")
	diaryAnalysisUsecase := usecases.NewDiaryAnalysisUsecase(diaryRepository)
	log.Println("[DEBUG] Lambda init: controllers.NewDiaryAnalysisController 開始")
	diaryAnalysisController := controllers.NewDiaryAnalysisController(diaryAnalysisUsecase)

	log.Println("[DEBUG] Lambda init: gin.Default() 開始")
	router := gin.Default()
	log.Println("[DEBUG] Lambda init: gin.Default() 完了")

	log.Println("[DEBUG] Lambda init: routes.SetupCORS 開始")
	routes.SetupCORS(router)
	log.Println("[DEBUG] Lambda init: routes.SetupCORS 完了")

	log.Println("[DEBUG] Lambda init: routes.SetupSwaggerEndpoints 開始")
	routes.SetupSwaggerEndpoints(router)
	log.Println("[DEBUG] Lambda init: routes.SetupSwaggerEndpoints 完了")

	log.Println("[DEBUG] Lambda init: routes.SetupAPIEndpoints 開始")
	routes.SetupAPIEndpoints(router, diaryController, diaryAnalysisController)
	log.Println("[DEBUG] Lambda init: routes.SetupAPIEndpoints 完了")

	log.Println("[DEBUG] Lambda init: ginadapter.New(router) 開始")
	ginLambda = ginadapter.New(router)
	log.Println("[DEBUG] Lambda init: ginadapter.New(router) 完了")
	log.Println("[DEBUG] Lambda init: 完了")
}

// Handler AWS Lambdaのハンドラー関数
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("[DEBUG] Lambda Handler: Received request: %s %s", req.HTTPMethod, req.Path)

	// Ginアダプターを使用してリクエストを処理
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	log.Println("[DEBUG] Lambda main: lambda.Start 開始")
	lambda.Start(Handler)
}
