### Stock Tracker API - OpenAPI Specification

```yaml
openapi: 3.0.3
info:
  title: Stock Tracker API
  version: 1.0.0
  description: |
    API para autenticación, acciones de brokers, recomendaciones, suscripciones y datos de mercado.
    
    Notas:
    - Autenticación via Bearer JWT en `Authorization: Bearer <token>`.
    - Varias rutas aceptan acceso invitado (guest) pero devuelven menos datos.
    - Límite de peticiones por tier; respuestas 429 cuando se excede.
servers:
  - url: http://localhost:8080
    description: Desarrollo local
tags:
  - name: Auth
  - name: Stocks
  - name: Brokers
  - name: Recommendations
  - name: Subscriptions
  - name: Market Data
  - name: Health
paths:
  /health:
    get:
      tags: [Health]
      summary: Health check
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/HealthResponse'
        '503':
          description: Servicio no disponible
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/HealthResponse'

  /api/v1/auth/register:
    post:
      tags: [Auth]
      summary: Registro de usuario
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/RegisterRequest'
      responses:
        '201':
          description: Usuario creado
          content:
            application/json:
              schema:
                type: object
                properties:
                  user:
                    $ref: '#/components/schemas/User'
                  tokens:
                    $ref: '#/components/schemas/TokenPair'
        '400': { $ref: '#/components/responses/BadRequest' }

  /api/v1/auth/login:
    post:
      tags: [Auth]
      summary: Inicio de sesión
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/LoginRequest'
      responses:
        '200':
          description: Autenticado
          content:
            application/json:
              schema:
                type: object
                properties:
                  user:
                    $ref: '#/components/schemas/User'
                  tokens:
                    $ref: '#/components/schemas/TokenPair'
        '400': { $ref: '#/components/responses/BadRequest' }
        '401': { $ref: '#/components/responses/Unauthorized' }

  /api/v1/auth/refresh:
    post:
      tags: [Auth]
      summary: Refrescar token
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/RefreshTokenRequest'
      responses:
        '200':
          description: Token renovado
          content:
            application/json:
              schema:
                type: object
                properties:
                  tokens:
                    $ref: '#/components/schemas/TokenPair'
        '400': { $ref: '#/components/responses/BadRequest' }
        '401': { $ref: '#/components/responses/Unauthorized' }

  /api/v1/stocks:
    get:
      tags: [Stocks]
      summary: Listado de eventos/acciones
      parameters:
        - in: query
          name: ticker
          schema: { type: string }
        - in: query
          name: company
          schema: { type: string }
        - in: query
          name: brokerage
          schema: { type: string }
        - in: query
          name: action
          schema: { type: string }
        - in: query
          name: rating_from
          schema: { type: string }
        - in: query
          name: rating_to
          schema: { type: string }
        - in: query
          name: date_from
          schema: { type: string, format: date-time }
        - in: query
          name: date_to
          schema: { type: string, format: date-time }
        - in: query
          name: sort_by
          schema: { type: string }
        - in: query
          name: sort_order
          schema: { type: string, enum: [asc, desc] }
        - in: query
          name: limit
          schema: { type: integer, minimum: 1, maximum: 1000 }
        - in: query
          name: offset
          schema: { type: integer, minimum: 0 }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/StockListResponse'
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/stocks/enhanced:
    get:
      tags: [Stocks]
      summary: Listado de acciones con filtros avanzados
      parameters:
        - in: query
          name: tickers
          description: Lista de tickers (soporta `tickers` o `tickers[]`)
          schema:
            type: array
            items: { type: string }
          style: form
          explode: true
        - in: query
          name: companies
          schema: { type: array, items: { type: string } }
          style: form
          explode: true
        - in: query
          name: brokerages
          schema: { type: array, items: { type: string } }
          style: form
          explode: true
        - in: query
          name: actions
          schema: { type: array, items: { type: string } }
          style: form
          explode: true
        - in: query
          name: rating_from
          schema: { type: string }
        - in: query
          name: rating_to
          schema: { type: string }
        - in: query
          name: date_from
          schema: { type: string, format: date-time }
        - in: query
          name: date_to
          schema: { type: string, format: date-time }
        - in: query
          name: date_ranges
          description: Rango(s) de fechas como pares "from,to" separados por `|`
          schema: { type: string }
        - in: query
          name: last_hours
          schema: { type: integer, minimum: 1 }
        - in: query
          name: last_days
          schema: { type: integer, minimum: 1 }
        - in: query
          name: last_weeks
          schema: { type: integer, minimum: 1 }
        - in: query
          name: last_months
          schema: { type: integer, minimum: 1 }
        - in: query
          name: target_from
          schema: { type: number }
        - in: query
          name: target_to
          schema: { type: number }
        - in: query
          name: min_target_change
          schema: { type: number }
        - in: query
          name: max_target_change
          schema: { type: number }
        - in: query
          name: has_target_price
          schema: { type: boolean }
        - in: query
          name: has_rating
          schema: { type: boolean }
        - in: query
          name: min_broker_score
          schema: { type: number }
        - in: query
          name: max_broker_score
          schema: { type: number }
        - in: query
          name: sort_by
          schema: { type: string }
        - in: query
          name: sort_order
          schema: { type: string, enum: [asc, desc] }
        - in: query
          name: limit
          schema: { type: integer, minimum: 1, maximum: 1000 }
        - in: query
          name: offset
          schema: { type: integer, minimum: 0 }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/StockListResponse'
        '400': { $ref: '#/components/responses/BadRequest' }
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/stocks/tickers:
    get:
      tags: [Stocks]
      summary: Tickers únicos
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items: { type: string }
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/stocks/companies:
    get:
      tags: [Stocks]
      summary: Empresas únicas
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items: { type: string }
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/stocks/stats:
    get:
      tags: [Stocks]
      summary: Estadísticas básicas
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: object
                    properties:
                      total_stocks: { type: integer }
                      last_updated: { type: string, format: date-time }
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/stocks/{ticker}:
    get:
      tags: [Stocks]
      summary: Eventos por ticker
      parameters:
        - in: path
          name: ticker
          required: true
          schema: { type: string }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Stock'
        '400': { $ref: '#/components/responses/BadRequest' }
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/stocks/{id}:
    get:
      tags: [Stocks]
      summary: Obtener por ID (no implementado)
      parameters:
        - in: path
          name: id
          required: true
          schema: { type: string, format: uuid }
      responses:
        '501': { $ref: '#/components/responses/NotImplemented' }
    put:
      tags: [Stocks]
      summary: Actualizar (no implementado)
      security: [{ bearerAuth: [] }]
      parameters:
        - in: path
          name: id
          required: true
          schema: { type: string, format: uuid }
      responses:
        '501': { $ref: '#/components/responses/NotImplemented' }
        '401': { $ref: '#/components/responses/Unauthorized' }
    delete:
      tags: [Stocks]
      summary: Eliminar (no implementado)
      security: [{ bearerAuth: [] }]
      parameters:
        - in: path
          name: id
          required: true
          schema: { type: string, format: uuid }
      responses:
        '501': { $ref: '#/components/responses/NotImplemented' }
        '401': { $ref: '#/components/responses/Unauthorized' }

  /api/v1/stocks:
    post:
      tags: [Stocks]
      summary: Crear (no implementado)
      security: [{ bearerAuth: [] }]
      responses:
        '501': { $ref: '#/components/responses/NotImplemented' }
        '401': { $ref: '#/components/responses/Unauthorized' }

  /api/v1/brokers/scores:
    get:
      tags: [Brokers]
      summary: Brokers con puntajes
      parameters:
        - in: query
          name: limit
          schema: { type: integer, minimum: 1 }
        - in: query
          name: order
          description: Orden de puntaje
          schema: { type: string, enum: [asc, desc] }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items: { $ref: '#/components/schemas/BrokerWithScore' }
                  meta:
                    type: object
                    properties:
                      count: { type: integer }
                      user_tier: { $ref: '#/components/schemas/UserTier' }
                      features:
                        type: array
                        items: { type: string }
                      rate_limit_remaining:
                        type: integer
                        description: Puede estar presente
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/recommendations:
    get:
      tags: [Recommendations]
      summary: Recomendaciones
      parameters:
        - in: query
          name: limit
          schema: { type: integer, default: 10, minimum: 1 }
        - in: query
          name: min_score
          schema: { type: number, minimum: 0, maximum: 1 }
        - in: query
          name: type
          schema:
            $ref: '#/components/schemas/RecommendationType'
        - in: query
          name: exclude
          description: Lista separada por coma de tickers a excluir
          schema: { type: string }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema: { $ref: '#/components/schemas/RecommendationResponse' }
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/recommendations/{ticker}:
    get:
      tags: [Recommendations]
      summary: Recomendación por ticker
      parameters:
        - in: path
          name: ticker
          required: true
          schema: { type: string }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  data: { $ref: '#/components/schemas/AggregatedRecommendation' }
                  meta:
                    type: object
                    properties:
                      user_tier: { $ref: '#/components/schemas/UserTier' }
                      features:
                        type: array
                        items: { type: string }
        '400': { $ref: '#/components/responses/BadRequest' }
        '404': { $ref: '#/components/responses/NotFound' }
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/recommendations/preview/{ticker}:
    get:
      tags: [Recommendations]
      summary: Vista previa de beneficios de upgrade
      security: [{ bearerAuth: [] }]
      parameters:
        - in: path
          name: ticker
          required: true
          schema: { type: string }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  current_tier:
                    type: object
                    properties:
                      tier: { $ref: '#/components/schemas/UserTier' }
                      data: { $ref: '#/components/schemas/AggregatedRecommendation' }
                  premium_preview:
                    nullable: true
                    allOf:
                      - $ref: '#/components/schemas/AggregatedRecommendation'
                  upgrade_benefits:
                    type: array
                    items: { type: string }
        '401': { $ref: '#/components/responses/Unauthorized' }
        '400': { $ref: '#/components/responses/BadRequest' }
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/subscriptions:
    post:
      tags: [Subscriptions]
      summary: Crear suscripción (simulación de pago)
      security: [{ bearerAuth: [] }]
      requestBody:
        required: true
        content:
          application/json:
            schema: { $ref: '#/components/schemas/PaymentSimulationRequest' }
      responses:
        '201':
          description: Creado
          content:
            application/json:
              schema: { $ref: '#/components/schemas/Subscription' }
        '400': { $ref: '#/components/responses/BadRequest' }
        '401': { $ref: '#/components/responses/Unauthorized' }

  /api/v1/subscriptions/{id}/payment:
    post:
      tags: [Subscriptions]
      summary: Simular pago de suscripción
      security: [{ bearerAuth: [] }]
      parameters:
        - in: path
          name: id
          required: true
          schema: { type: string, format: uuid }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  message: { type: string }
        '400': { $ref: '#/components/responses/BadRequest' }
        '401': { $ref: '#/components/responses/Unauthorized' }

  /api/v1/subscriptions/{id}:
    get:
      tags: [Subscriptions]
      summary: Obtener suscripción por ID
      security: [{ bearerAuth: [] }]
      parameters:
        - in: path
          name: id
          required: true
          schema: { type: string, format: uuid }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema: { $ref: '#/components/schemas/Subscription' }
        '401': { $ref: '#/components/responses/Unauthorized' }
        '404': { $ref: '#/components/responses/NotFound' }

  /api/v1/subscriptions/active:
    get:
      tags: [Subscriptions]
      summary: Obtener suscripción activa del usuario
      security: [{ bearerAuth: [] }]
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema: { $ref: '#/components/schemas/Subscription' }
        '401': { $ref: '#/components/responses/Unauthorized' }
        '404': { $ref: '#/components/responses/NotFound' }

  /api/v1/subscriptions/plans:
    get:
      tags: [Subscriptions]
      summary: Listar planes disponibles
      security: [{ bearerAuth: [] }]
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: array
                items: { $ref: '#/components/schemas/SubscriptionPlan' }

  /api/v1/market-data/analysis/{ticker}:
    get:
      tags: [Market Data]
      summary: Análisis integral por ticker
      parameters:
        - in: path
          name: ticker
          required: true
          schema: { type: string }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/MarketDataResponseAnalysis'
        '400': { $ref: '#/components/responses/BadRequest' }
        '404': { $ref: '#/components/responses/NotFound' }
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/market-data/trend/{ticker}:
    get:
      tags: [Market Data]
      summary: Tendencia por ticker
      parameters:
        - in: path
          name: ticker
          required: true
          schema: { type: string }
        - in: query
          name: period
          schema: { type: string, enum: [1d, 1w, 1m, 3m] }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/MarketDataResponseTrend'
        '400': { $ref: '#/components/responses/BadRequest' }
        '404': { $ref: '#/components/responses/NotFound' }
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/market-data/summary:
    get:
      tags: [Market Data]
      summary: Resumen de mercado
      parameters:
        - in: query
          name: period
          schema: { type: string, enum: [1d, 1w, 1m, 3m] }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/MarketDataResponseSummary'
        '400': { $ref: '#/components/responses/BadRequest' }
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/market-data/top-performers:
    get:
      tags: [Market Data]
      summary: Mejores desempeños
      parameters:
        - in: query
          name: period
          schema: { type: string, enum: [1d, 1w, 1m, 3m] }
        - in: query
          name: limit
          schema: { type: integer, minimum: 1, maximum: 100 }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/MarketDataResponseAnalysisList'
        '400': { $ref: '#/components/responses/BadRequest' }
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/market-data/worst-performers:
    get:
      tags: [Market Data]
      summary: Peores desempeños
      parameters:
        - in: query
          name: period
          schema: { type: string, enum: [1d, 1w, 1m, 3m] }
        - in: query
          name: limit
          schema: { type: integer, minimum: 1, maximum: 100 }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/MarketDataResponseAnalysisList'
        '400': { $ref: '#/components/responses/BadRequest' }
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/market-data/latest/{ticker}:
    get:
      tags: [Market Data]
      summary: Último dato de mercado por ticker
      parameters:
        - in: path
          name: ticker
          required: true
          schema: { type: string }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/MarketDataResponseLatest'
        '400': { $ref: '#/components/responses/BadRequest' }
        '404': { $ref: '#/components/responses/NotFound' }
        '500': { $ref: '#/components/responses/InternalError' }

  /api/v1/market-data/{ticker}:
    get:
      tags: [Market Data]
      summary: Datos de mercado por ticker (con filtros)
      parameters:
        - in: path
          name: ticker
          required: true
          schema: { type: string }
        - in: query
          name: data_source
          schema: { type: string, enum: [yahoo_finance, alpha_vantage, manual] }
        - in: query
          name: data_quality
          schema: { type: string, enum: [excellent, good, fair, poor] }
        - in: query
          name: start_date
          schema: { type: string, format: date }
        - in: query
          name: end_date
          schema: { type: string, format: date }
        - in: query
          name: min_price
          schema: { type: number }
        - in: query
          name: max_price
          schema: { type: number }
        - in: query
          name: min_change
          schema: { type: number }
        - in: query
          name: max_change
          schema: { type: number }
        - in: query
          name: min_volume
          schema: { type: integer }
        - in: query
          name: max_volume
          schema: { type: integer }
        - in: query
          name: trend_direction
          schema: { type: string, enum: [bullish, bearish, neutral] }
        - in: query
          name: risk_level
          schema: { type: string, enum: [low, medium, high] }
        - in: query
          name: sort_by
          schema: { type: string, default: data_timestamp }
        - in: query
          name: sort_order
          schema: { type: string, enum: [asc, desc], default: desc }
        - in: query
          name: limit
          schema: { type: integer, minimum: 1, default: 50 }
        - in: query
          name: offset
          schema: { type: integer, minimum: 0, default: 0 }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/MarketDataResponseList'
        '400': { $ref: '#/components/responses/BadRequest' }
        '404': { $ref: '#/components/responses/NotFound' }
        '500': { $ref: '#/components/responses/InternalError' }

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  responses:
    BadRequest:
      description: Petición inválida
      content:
        application/json:
          schema: { $ref: '#/components/schemas/ErrorResponse' }
    Unauthorized:
      description: No autorizado
      content:
        application/json:
          schema: { $ref: '#/components/schemas/ErrorResponse' }
    Forbidden:
      description: Prohibido
      content:
        application/json:
          schema: { $ref: '#/components/schemas/ErrorResponse' }
    NotFound:
      description: No encontrado
      content:
        application/json:
          schema: { $ref: '#/components/schemas/ErrorResponse' }
    TooManyRequests:
      description: Límite de peticiones excedido
      content:
        application/json:
          schema: { $ref: '#/components/schemas/ErrorResponse' }
    InternalError:
      description: Error interno del servidor
      content:
        application/json:
          schema: { $ref: '#/components/schemas/ErrorResponse' }
    NotImplemented:
      description: No implementado
      content:
        application/json:
          schema: { $ref: '#/components/schemas/ErrorResponse' }

  schemas:
    ErrorResponse:
      type: object
      properties:
        error: { type: string }

    HealthResponse:
      type: object
      properties:
        status: { type: string }
        database: { type: string }
        timestamp: { type: string }

    UserTier:
      type: string
      enum: [guest, basic, premium]

    TokenPair:
      type: object
      properties:
        access_token: { type: string }
        refresh_token: { type: string }
        expires_in: { type: integer }

    User:
      type: object
      properties:
        id: { type: string, format: uuid }
        email: { type: string, format: email }
        first_name: { type: string }
        last_name: { type: string }
        tier: { $ref: '#/components/schemas/UserTier' }
        is_verified: { type: boolean }
        last_login: { type: string, format: date-time, nullable: true }
        created_at: { type: string, format: date-time }
        updated_at: { type: string, format: date-time }

    RegisterRequest:
      type: object
      required: [email, password, first_name, last_name]
      properties:
        email: { type: string, format: email }
        password: { type: string, minLength: 8 }
        first_name: { type: string }
        last_name: { type: string }

    LoginRequest:
      type: object
      required: [email, password]
      properties:
        email: { type: string, format: email }
        password: { type: string }

    RefreshTokenRequest:
      type: object
      required: [refresh_token]
      properties:
        refresh_token: { type: string }

    Stock:
      type: object
      properties:
        id: { type: string, format: uuid }
        ticker: { type: string }
        company: { type: string }
        brokerage: { type: string }
        broker_id: { type: string, format: uuid }
        action: { type: string }
        rating_from: { type: string }
        rating_to: { type: string }
        target_from: { type: number }
        target_to: { type: number }
        event_time: { type: string, format: date-time }
        created_at: { type: string, format: date-time }
        updated_at: { type: string, format: date-time }

    StockPagination:
      type: object
      properties:
        page: { type: integer }
        limit: { type: integer }
        total_pages: { type: integer }
        total_items: { type: integer }
        has_next: { type: boolean }
        has_prev: { type: boolean }

    StockListResponse:
      type: object
      properties:
        data:
          type: array
          items: { $ref: '#/components/schemas/Stock' }
        pagination: { $ref: '#/components/schemas/StockPagination' }

    BrokerWithScore:
      type: object
      properties:
        id: { type: string, format: uuid }
        name: { type: string }
        credibility_score: { type: number }
        report_count: { type: integer }
        calculated_score: { type: number }
        created_at: { type: string }
        updated_at: { type: string }

    RecommendationType:
      type: string
      enum: [strong_buy, buy, hold, sell, strong_sell]

    ScoringFactor:
      type: object
      properties:
        name: { type: string }
        score: { type: number }
        weight: { type: number }
        explanation: { type: string }
        tier: { type: string }

    ExternalStockData:
      type: object
      properties:
        current_price: { type: number }
        day_change: { type: number }
        day_change_percent: { type: number }
        volume: { type: integer }
        market_cap: { type: integer }
        pe_ratio: { type: number, nullable: true }
        dividend_yield: { type: number, nullable: true }
        week_52_high: { type: number, nullable: true }
        week_52_low: { type: number, nullable: true }
        avg_volume: { type: integer, nullable: true }
        last_updated: { type: string, format: date-time }

    AIGeneratedInsights:
      type: object
      properties:
        market_sentiment: { type: string }
        sentiment_score: { type: number }
        risk_assessment: { type: string }
        key_drivers:
          type: array
          items: { type: string }
        competitor_analysis:
          type: array
          items: { type: string }
        news_impact: { type: number, nullable: true }
        technical_indicators:
          type: array
          items: { type: string }
        generated_at: { type: string, format: date-time }

    AggregatedRecommendation:
      type: object
      properties:
        id: { type: string, format: uuid }
        ticker: { type: string }
        company_name: { type: string }
        total_events: { type: integer }
        positive_events: { type: integer }
        negative_events: { type: integer }
        avg_target_change: { type: number }
        latest_target_price: { type: number }
        broker_consensus: { type: number }
        basic_score: { type: number }
        confidence: { type: number }
        recommendation_type: { $ref: '#/components/schemas/RecommendationType' }
        scoring_factors:
          type: array
          items: { $ref: '#/components/schemas/ScoringFactor' }
        tier: { type: string }
        external_data: { $ref: '#/components/schemas/ExternalStockData' }
        ai_insights: { $ref: '#/components/schemas/AIGeneratedInsights' }
        last_event_time: { type: string, format: date-time }
        created_at: { type: string, format: date-time }
        expires_at: { type: string, format: date-time }

    RecommendationMeta:
      type: object
      properties:
        count: { type: integer }
        user_tier: { $ref: '#/components/schemas/UserTier' }
        features:
          type: array
          items: { type: string }
        cache_hit: { type: boolean }
        generation_time: { type: string }
        rate_limit_remaining: { type: integer, nullable: true }

    RecommendationResponse:
      type: object
      properties:
        data:
          type: array
          items: { $ref: '#/components/schemas/AggregatedRecommendation' }
        meta: { $ref: '#/components/schemas/RecommendationMeta' }

    SubscriptionPlan:
      type: object
      properties:
        plan: { type: string, enum: [monthly, yearly] }
        name: { type: string }
        price: { type: number }
        currency: { type: string }
        duration: { type: string }
        features:
          type: array
          items: { type: string }

    PaymentSimulationRequest:
      type: object
      required: [plan]
      properties:
        plan: { type: string, enum: [monthly, yearly] }

    Subscription:
      type: object
      properties:
        id: { type: string, format: uuid }
        user_id: { type: string, format: uuid }
        plan: { type: string, enum: [monthly, yearly] }
        status: { type: string, enum: [pending, active, cancelled, expired] }
        price: { type: number }
        currency: { type: string }
        start_date: { type: string, format: date-time }
        end_date: { type: string, format: date-time }
        payment_reference: { type: string }
        created_at: { type: string, format: date-time }
        updated_at: { type: string, format: date-time }

    MarketData:
      type: object
      properties:
        id: { type: string, format: uuid }
        ticker: { type: string }
        data_source: { type: string, enum: [yahoo_finance, alpha_vantage, manual] }
        data_quality: { type: string, enum: [excellent, good, fair, poor] }
        current_price: { type: number }
        day_change: { type: number }
        day_change_percent: { type: number }
        volume: { type: integer }
        market_cap: { type: integer, nullable: true }
        pe_ratio: { type: number, nullable: true }
        dividend_yield: { type: number, nullable: true }
        week_52_high: { type: number, nullable: true }
        week_52_low: { type: number, nullable: true }
        avg_volume: { type: integer, nullable: true }
        collected_at: { type: string, format: date-time }
        data_timestamp: { type: string, format: date-time }
        created_at: { type: string, format: date-time }
        updated_at: { type: string, format: date-time }

    MarketDataAnalysis:
      type: object
      properties:
        id: { type: string, format: uuid }
        ticker: { type: string }
        current_price: { type: number }
        day_change: { type: number }
        day_change_percent: { type: number }
        volume: { type: integer }
        week_52_high: { type: number, nullable: true }
        week_52_low: { type: number, nullable: true }
        data_timestamp: { type: string, format: date-time }
        collected_at: { type: string, format: date-time }
        data_quality: { type: string }
        data_source: { type: string }
        price_position: { type: number }
        volume_activity: { type: number }
        volatility_score: { type: number }
        trend_direction: { type: string, enum: [bullish, bearish, neutral] }
        risk_level: { type: string, enum: [low, medium, high] }

    MarketDataTrend:
      type: object
      properties:
        ticker: { type: string }
        period: { type: string, enum: [1d, 1w, 1m, 3m] }
        start_price: { type: number }
        end_price: { type: number }
        total_change: { type: number }
        total_change_percent: { type: number }
        max_price: { type: number }
        min_price: { type: number }
        avg_volume: { type: integer }
        volatility: { type: number }
        trend_strength: { type: string, enum: [strong, moderate, weak] }
        direction: { type: string, enum: [up, down, sideways] }
        data_points: { type: integer }

    MarketDataSummary:
      type: object
      properties:
        total_records: { type: integer }
        unique_tickers: { type: integer }
        avg_price: { type: number }
        avg_day_change: { type: number }
        avg_day_change_percent: { type: number }
        total_volume: { type: integer }
        avg_volume: { type: integer }
        most_active_ticker: { type: string }
        best_performer: { type: string }
        worst_performer: { type: string }
        bullish_count: { type: integer }
        bearish_count: { type: integer }
        neutral_count: { type: integer }
        period: { type: string }
        last_updated: { type: string, format: date-time }

    MarketDataResponse:
      type: object
      properties:
        data: { }
        pagination:
          type: object
          properties:
            total: { type: integer }
            page: { type: integer }
            per_page: { type: integer }
            total_pages: { type: integer }
            has_next: { type: boolean }
            has_previous: { type: boolean }
        message: { type: string }
        metadata:
          type: object
          properties:
            generated_at: { type: string, format: date-time }
            data_points: { type: integer }
            period: { type: string }
            filters:
              type: object

    MarketDataResponseAnalysis:
      allOf:
        - $ref: '#/components/schemas/MarketDataResponse'
        - type: object
          properties:
            data: { $ref: '#/components/schemas/MarketDataAnalysis' }

    MarketDataResponseTrend:
      allOf:
        - $ref: '#/components/schemas/MarketDataResponse'
        - type: object
          properties:
            data: { $ref: '#/components/schemas/MarketDataTrend' }

    MarketDataResponseSummary:
      allOf:
        - $ref: '#/components/schemas/MarketDataResponse'
        - type: object
          properties:
            data: { $ref: '#/components/schemas/MarketDataSummary' }

    MarketDataResponseAnalysisList:
      allOf:
        - $ref: '#/components/schemas/MarketDataResponse'
        - type: object
          properties:
            data:
              type: array
              items: { $ref: '#/components/schemas/MarketDataAnalysis' }

    MarketDataResponseLatest:
      allOf:
        - $ref: '#/components/schemas/MarketDataResponse'
        - type: object
          properties:
            data: { $ref: '#/components/schemas/MarketData' }

    MarketDataResponseList:
      allOf:
        - $ref: '#/components/schemas/MarketDataResponse'
        - type: object
          properties:
            data:
              type: array
              items: { $ref: '#/components/schemas/MarketData' }
```

