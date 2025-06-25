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
