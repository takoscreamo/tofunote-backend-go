# TOFU NOTE Backend - オニオンアーキテクチャ（Mermaid図）

## 概要
TOFU NOTE Backendは、豆腐メンタルの可視化・記録アプリケーションのバックエンドAPIです。軽量DDD/オニオンアーキテクチャを採用し、保守性・拡張性を重視した設計となっています。

## オニオンアーキテクチャ構造図

### 1. 全体構造図

```mermaid
graph TB
    subgraph "🚀 エントリーポイント層"
        A1[main.go<br/>Lambda用エントリーポイント]
        A2[cmd/local/main.go<br/>ローカル開発用エントリーポイント]
    end
    
    subgraph "🌐 プレゼンテーション層"
        B1[api/controllers/<br/>HTTPハンドラ]
        B2[routes/<br/>ルーティング・ミドルウェア]
        B3[middleware/<br/>JWT認証等]
    end
    
    subgraph "🎯 ユースケース層"
        C1[usecases/diary_usecase.go<br/>日記ビジネスロジック]
        C2[usecases/user_withdraw_usecase.go<br/>ユーザー退会ロジック]
        C3[usecases/diary_analysis_usecase.go<br/>日記分析ロジック]
    end
    
    subgraph "🏛️ ドメイン層"
        D1[domain/diary/<br/>日記エンティティ・値オブジェクト]
        D2[domain/user/<br/>ユーザーエンティティ・値オブジェクト]
        D3[domain/*/repository.go<br/>リポジトリインターフェース]
    end
    
    subgraph "🔄 リポジトリ実装層"
        E1[repositories/diary_repositories.go<br/>GORM実装]
        E2[repositories/user_repositories.go<br/>GORM実装]
    end
    
    subgraph "🏗️ インフラ層"
        F1[infra/db/<br/>DBモデル・接続]
        F2[infra/migrations/<br/>DBマイグレーション]
        F3[infra/jwt.go<br/>JWT処理]
        F4[infra/initializer.go<br/>依存注入・初期化]
    end
    
    subgraph "📁 補助ディレクトリ"
        G1[scripts/<br/>データ移行スクリプト]
        G2[openapi.yml<br/>API仕様書]
        G3[Makefile<br/>ビルドコマンド]
    end
    
    A1 --> B1
    A2 --> B1
    B1 --> C1
    B1 --> C2
    B1 --> C3
    C1 --> D3
    C2 --> D3
    C3 --> D3
    D3 --> E1
    D3 --> E2
    E1 --> F1
    E2 --> F1
    F1 --> F2
    F1 --> F3
    F1 --> F4
    
    style A1 fill:#e1f5fe
    style A2 fill:#e1f5fe
    style B1 fill:#f3e5f5
    style B2 fill:#f3e5f5
    style B3 fill:#f3e5f5
    style C1 fill:#e8f5e8
    style C2 fill:#e8f5e8
    style C3 fill:#e8f5e8
    style D1 fill:#fff3e0
    style D2 fill:#fff3e0
    style D3 fill:#fff3e0
    style E1 fill:#fce4ec
    style E2 fill:#fce4ec
    style F1 fill:#f1f8e9
    style F2 fill:#f1f8e9
    style F3 fill:#f1f8e9
    style F4 fill:#f1f8e9
    style G1 fill:#fafafa
    style G2 fill:#fafafa
    style G3 fill:#fafafa
```

### 2. 依存関係の流れ図

```mermaid
flowchart TD
    A[外部リクエスト<br/>HTTP/HTTPS] --> B[🌐 プレゼンテーション層<br/>api/controllers/]
    B --> C[🎯 ユースケース層<br/>usecases/]
    C --> D[🏛️ ドメイン層<br/>domain/]
    D --> E[🔄 リポジトリ実装層<br/>repositories/]
    E --> F[🏗️ インフラ層<br/>infra/]
    F --> G[外部システム<br/>PostgreSQL, JWT, etc.]
    
    style A fill:#e3f2fd
    style B fill:#f3e5f5
    style C fill:#e8f5e8
    style D fill:#fff3e0
    style E fill:#fce4ec
    style F fill:#f1f8e9
    style G fill:#ffebee
```

### 3. 日記機能の詳細フロー

```mermaid
sequenceDiagram
    participant Client as クライアント
    participant Controller as DiaryController
    participant Usecase as DiaryUsecase
    participant Domain as Diary Entity
    participant Repository as DiaryRepository
    participant DB as PostgreSQL
    
    Client->>Controller: POST /api/diaries
    Controller->>Controller: JWT認証・バリデーション
    Controller->>Usecase: Create(diary)
    Usecase->>Domain: ドメインロジック実行
    Usecase->>Repository: Create(diary)
    Repository->>DB: INSERT INTO diaries
    DB-->>Repository: 成功レスポンス
    Repository-->>Usecase: 成功
    Usecase-->>Controller: 成功
    Controller-->>Client: 201 Created
    
    Note over Controller,DB: 依存関係注入により<br/>テスト時にモックに置き換え可能
```

### 4. レイヤー間のインターフェース関係

```mermaid
classDiagram
    class DiaryController {
        +usecase IDiaryUsecase
        +FindAll()
        +Create()
        +Update()
        +Delete()
    }
    
    class IDiaryUsecase {
        <<interface>>
        +FindAll()
        +Create()
        +Update()
        +Delete()
    }
    
    class DiaryUsecase {
        +repository DiaryRepository
        +FindAll()
        +Create()
        +Update()
        +Delete()
    }
    
    class DiaryRepository {
        <<interface>>
        +FindAll()
        +Create()
        +Update()
        +Delete()
    }
    
    class DiaryRepositoryImpl {
        +db *gorm.DB
        +FindAll()
        +Create()
        +Update()
        +Delete()
    }
    
    class Diary {
        +ID string
        +UserID string
        +Date string
        +Mental Mental
        +Diary string
    }
    
    DiaryController --> IDiaryUsecase : 依存
    DiaryUsecase ..|> IDiaryUsecase : 実装
    DiaryUsecase --> DiaryRepository : 依存
    DiaryRepositoryImpl ..|> DiaryRepository : 実装
    DiaryUsecase --> Diary : 使用
```

### 5. データ変換フロー

```mermaid
flowchart LR
    A[HTTP Request<br/>JSON] --> B[Controller<br/>DTO変換]
    B --> C[Usecase<br/>ビジネスロジック]
    C --> D[Domain<br/>エンティティ]
    D --> E[Repository<br/>DBモデル変換]
    E --> F[GORM<br/>SQL実行]
    F --> G[PostgreSQL]
    
    G --> H[GORM<br/>結果取得]
    H --> I[Repository<br/>ドメイン変換]
    I --> J[Usecase<br/>結果返却]
    J --> K[Controller<br/>レスポンスDTO変換]
    K --> L[HTTP Response<br/>JSON]
    
    style A fill:#e3f2fd
    style B fill:#f3e5f5
    style C fill:#e8f5e8
    style D fill:#fff3e0
    style E fill:#fce4ec
    style F fill:#f1f8e9
    style G fill:#ffebee
    style H fill:#f1f8e9
    style I fill:#fce4ec
    style J fill:#e8f5e8
    style K fill:#f3e5f5
    style L fill:#e3f2fd
```

### 6. テスト構造図

```mermaid
graph TB
    subgraph "🧪 テスト層"
        T1[api/controllers/*_test.go<br/>コントローラーテスト]
        T2[usecases/*_test.go<br/>ユースケーステスト]
        T3[repositories/*_test.go<br/>リポジトリテスト]
        T4[domain/*_test.go<br/>ドメインテスト]
        T5[infra/*_test.go<br/>インフラテスト]
    end
    
    subgraph "🔧 テストツール"
        M1[go-sqlmock<br/>DBモック]
        M2[testify<br/>アサーション]
        M3[標準testing<br/>テストフレームワーク]
    end
    
    T1 --> M2
    T1 --> M3
    T2 --> M2
    T2 --> M3
    T3 --> M1
    T3 --> M2
    T3 --> M3
    T4 --> M2
    T4 --> M3
    T5 --> M2
    T5 --> M3
    
    style T1 fill:#fff8e1
    style T2 fill:#fff8e1
    style T3 fill:#fff8e1
    style T4 fill:#fff8e1
    style T5 fill:#fff8e1
    style M1 fill:#e0f2f1
    style M2 fill:#e0f2f1
    style M3 fill:#e0f2f1
```

## オニオンアーキテクチャの特徴

### 🎯 **依存関係の方向**
- 内側の層（ドメイン層）は外側の層に依存しない
- 外側の層は内側の層のインターフェースに依存
- 技術的な変更がビジネスロジックに影響しない

### 🧪 **テスタビリティ**
- 各層を独立してテスト可能
- モックを使用した単体テストが容易
- テーブル駆動テストで品質担保

### 🔧 **保守性・拡張性**
- 新機能追加時の影響範囲が限定的
- 技術スタック変更時の影響を最小化
- コードの責務が明確に分離

## 技術スタック

- **言語**: Go 1.24.1
- **Webフレームワーク**: Gin
- **ORM**: GORM
- **データベース**: PostgreSQL
- **認証**: JWT
- **API仕様**: OpenAPI/Swagger
- **テスト**: 標準testing + testify
- **マイグレーション**: golang-migrate

このアーキテクチャにより、TOFU NOTE Backendは保守性・拡張性に優れた、高品質なAPIを提供しています。 