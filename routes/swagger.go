package routes

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// SetupSwaggerEndpoints SwaggerUIとOpenAPI仕様書のエンドポイントを設定
func SetupSwaggerEndpoints(router *gin.Engine) {
	// OpenAPI仕様書を提供するエンドポイント
	router.GET("/openapi.yml", func(c *gin.Context) {
		// Lambda環境では実行ディレクトリが /var/task になるため、
		// 相対パスでファイルを探す
		openapiPath := "openapi.yml"

		// ファイルが存在するかチェック
		if _, err := os.Stat(openapiPath); os.IsNotExist(err) {
			// ファイルが存在しない場合は、実行ファイルと同じディレクトリを探す
			execPath, err := os.Executable()
			if err == nil {
				execDir := filepath.Dir(execPath)
				openapiPath = filepath.Join(execDir, "openapi.yml")
			}
		}

		// ファイルが存在するか最終チェック
		if _, err := os.Stat(openapiPath); os.IsNotExist(err) {
			// ファイルが見つからない場合は、埋め込まれたOpenAPI仕様書を提供
			openapiContent := `openapi: 3.0.0
info:
  title: 日記API
  description: ユーザーの日記を管理するためのRESTful API
  version: 1.0.0

servers:
  - url: https://api.feelog.takoscreamo.com/api
    description: 本番環境
  - url: http://localhost:8080/api
    description: ローカル環境

paths:
  /me/diaries:
    get:
      summary: 日記一覧取得
      description: 現在のユーザーの日記一覧を取得します
      responses:
        '200':
          description: 成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Diary'
        '500':
          description: サーバーエラー
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
    post:
      summary: 日記作成
      description: 現在のユーザーの新しい日記を作成します
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateDiaryDTO'
      responses:
        '201':
          description: 作成成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    $ref: '#/components/schemas/Diary'
        '400':
          description: リクエストが不正
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '409':
          description: 複合ユニークキー制約違反（同じユーザーの同じ日付の日記が既に存在）
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '500':
          description: サーバーエラー
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /me/diaries/range:
    get:
      summary: 期間指定日記取得
      description: 現在のユーザーの指定された期間の日記を取得します
      parameters:
        - name: start_date
          in: query
          required: true
          schema:
            type: string
            format: date
          description: 開始日（YYYY-MM-DD形式）
        - name: end_date
          in: query
          required: true
          schema:
            type: string
            format: date
          description: 終了日（YYYY-MM-DD形式）
      responses:
        '200':
          description: 取得成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Diary'
        '400':
          description: リクエストが不正（パラメータ不足）
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '500':
          description: サーバーエラー
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /me/diaries/{date}:
    get:
      summary: 日記取得
      description: 現在のユーザーの指定された日付の日記を取得します
      parameters:
        - name: date
          in: path
          required: true
          schema:
            type: string
            format: date
          description: 日付（YYYY-MM-DD形式）
      responses:
        '200':
          description: 取得成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    $ref: '#/components/schemas/Diary'
        '400':
          description: リクエストが不正（無効な日付）
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '404':
          description: 指定された日記が見つかりません
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '500':
          description: サーバーエラー
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
    put:
      summary: 日記更新
      description: 現在のユーザーの指定された日付の日記を更新します
      parameters:
        - name: date
          in: path
          required: true
          schema:
            type: string
            format: date
          description: 日付（YYYY-MM-DD形式）
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdateDiaryDTO'
      responses:
        '200':
          description: 更新成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    $ref: '#/components/schemas/Diary'
        '400':
          description: リクエストが不正
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '404':
          description: 指定された日記が見つかりません
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '500':
          description: サーバーエラー
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
    delete:
      summary: 日記削除
      description: 現在のユーザーの指定された日付の日記を削除します
      parameters:
        - name: date
          in: path
          required: true
          schema:
            type: string
            format: date
          description: 日付（YYYY-MM-DD形式）
      responses:
        '200':
          description: 削除成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: "日記が正常に削除されました"
        '400':
          description: リクエストが不正（無効な日付）
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '404':
          description: 指定された日記が見つかりません
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '500':
          description: サーバーエラー
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /me/analyze-diaries:
    get:
      summary: 日記分析
      description: 現在のユーザーの日記を分析してメンタルスコアの傾向を取得します
      responses:
        '200':
          description: 分析成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    $ref: '#/components/schemas/AnalysisResult'
        '500':
          description: サーバーエラー
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

components:
  schemas:
    Diary:
      type: object
      properties:
        id:
          type: integer
          example: 1
        user_id:
          type: integer
          example: 1
        date:
          type: string
          format: date
          example: "2025-06-25"
        mental:
          type: integer
          minimum: 0
          maximum: 10
          example: 7
        diary:
          type: string
          example: "今日は良い一日でした"
        created_at:
          type: string
          format: date-time
          example: "2025-06-25T10:00:00Z"
        updated_at:
          type: string
          format: date-time
          example: "2025-06-25T10:00:00Z"
      required:
        - user_id
        - date
        - mental
        - diary

    CreateDiaryDTO:
      type: object
      properties:
        date:
          type: string
          format: date
          example: "2025-06-25"
        mental:
          type: integer
          minimum: 0
          maximum: 10
          example: 7
        diary:
          type: string
          example: "今日は良い一日でした"
      required:
        - date
        - mental
        - diary

    UpdateDiaryDTO:
      type: object
      properties:
        mental:
          type: integer
          minimum: 0
          maximum: 10
          example: 8
        diary:
          type: string
          example: "更新された日記の内容"
      required:
        - mental
        - diary

    AnalysisResult:
      type: object
      properties:
        average_mental:
          type: number
          format: float
          example: 7.5
        total_entries:
          type: integer
          example: 30
        trend:
          type: string
          example: "上昇傾向"
        recommendations:
          type: array
          items:
            type: string
          example: ["定期的な運動を心がけましょう", "十分な睡眠を取るようにしましょう"]

    Error:
      type: object
      properties:
        error:
          type: string
          example: "エラーメッセージ"
      required:
        - error
`

			c.Header("Content-Type", "application/x-yaml")
			c.String(200, openapiContent)
			return
		}

		// ファイルを提供
		c.File(openapiPath)
	})

	// SwaggerUIを提供するエンドポイント
	router.GET("/swagger", func(c *gin.Context) {
		swaggerHTML := `<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>日記API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin:0;
            background: #fafafa;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/openapi.yml',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`
		c.Header("Content-Type", "text/html")
		c.String(200, swaggerHTML)
	})
}
