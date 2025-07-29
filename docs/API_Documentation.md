# API de Stock Tracker - Documentación Completa

## Tabla de Contenidos
1. [Información General](#información-general)
2. [Sistema de Autenticación](#sistema-de-autenticación)
3. [Sistema de Niveles de Usuario](#sistema-de-niveles-de-usuario)
4. [Endpoints de Autenticación](#endpoints-de-autenticación)
5. [Endpoints de Stocks](#endpoints-de-stocks)
6. [Endpoints de Recomendaciones](#endpoints-de-recomendaciones)
7. [Endpoints de Suscripciones](#endpoints-de-suscripciones)
8. [Endpoints de Salud](#endpoints-de-salud)
9. [Modelos de Datos](#modelos-de-datos)
10. [Filtros y Paginación](#filtros-y-paginación)
11. [Rate Limiting](#rate-limiting)
12. [Códigos de Error](#códigos-de-error)

## Información General

La API de Stock Tracker es un sistema de recomendaciones de acciones con diferentes niveles de acceso basados en suscripciones. La API está construida en Go usando Chi router y sigue una arquitectura limpia.

**Base URL:** `http://localhost:{PORT}/api/v1`

**Formato de respuesta:** JSON

**Autenticación:** JWT Bearer Token

## Sistema de Autenticación

La API utiliza JWT (JSON Web Tokens) para la autenticación con access tokens y refresh tokens:

- **Access Token:** Válido por períodos cortos, incluido en el header `Authorization: Bearer <token>`
- **Refresh Token:** Válido por períodos largos, usado para obtener nuevos access tokens

## Sistema de Niveles de Usuario

### Niveles de Usuario (UserTier)

1. **GUEST** (`"guest"`)
   - Usuario no registrado
   - Acceso limitado a funcionalidades básicas
   - Rate limit: 100 requests/hora
   - Máximo 10 recomendaciones

2. **BASIC** (`"basic"`)
   - Usuario registrado
   - Acceso a datos en tiempo real y APIs externas
   - Rate limit: 500 requests/hora
   - Máximo 25 recomendaciones
   - Acceso a external_data en recomendaciones

3. **PREMIUM** (`"premium"`)
   - Usuario con suscripción pagada
   - Acceso completo incluyendo IA insights
   - Rate limit: 2000 requests/hora
   - Máximo 100 recomendaciones
   - Acceso a external_data y ai_insights

## Endpoints de Autenticación

### 1. Registro de Usuario
```
POST /api/v1/auth/register
```

**Descripción:** Registra un nuevo usuario en el sistema

**Body (JSON):**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Validaciones:**
- Email válido requerido
- Password mínimo 8 caracteres
- First name y last name requeridos (1-100 caracteres)

**Respuesta exitosa (201):**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "tier": "basic",
    "is_verified": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "tokens": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_in": 3600
  }
}
```

**Funcionamiento interno:**
- Valida que el email no exista previamente
- Hashea la contraseña con bcrypt
- Crea el usuario con tier BASIC por defecto
- Genera tokens JWT
- Crea una sesión en la base de datos

### 2. Login de Usuario
```
POST /api/v1/auth/login
```

**Descripción:** Autentica un usuario existente

**Body (JSON):**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Respuesta exitosa (200):**
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "tier": "basic",
    "is_verified": false,
    "last_login": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "tokens": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_in": 3600
  }
}
```

**Funcionamiento interno:**
- Busca el usuario por email
- Valida la contraseña con bcrypt
- Actualiza last_login
- Genera nuevos tokens JWT
- Crea nueva sesión

### 3. Refresh Token
```
POST /api/v1/auth/refresh
```

**Descripción:** Renueva el access token usando un refresh token válido

**Body (JSON):**
```json
{
  "refresh_token": "eyJ..."
}
```

**Respuesta exitosa (200):**
```json
{
  "tokens": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_in": 3600
  }
}
```

**Funcionamiento interno:**
- Valida el refresh token en la base de datos
- Verifica que no haya expirado
- Elimina la sesión anterior
- Genera nuevos tokens
- Crea nueva sesión

## Endpoints de Stocks

### 1. Obtener Lista de Stocks
```
GET /api/v1/stocks
```

**Descripción:** Obtiene una lista paginada de eventos de stocks con filtros opcionales

**Autenticación:** Opcional (Guest users tienen acceso con limitaciones)

**Query Parameters:**
- `ticker` (string): Filtrar por símbolo de acción
- `company` (string): Filtrar por nombre de empresa
- `brokerage` (string): Filtrar por casa de corretaje
- `action` (string): Filtrar por tipo de acción
- `limit` (int): Número máximo de resultados (default: 50, max: 1000)
- `offset` (int): Número de resultados a omitir para paginación
- `sort_by` (string): Campo para ordenar (default: "event_time")
- `sort_order` (string): Orden ascendente o descendente (default: "desc")

**Ejemplo:**
```
GET /api/v1/stocks?ticker=AAPL&limit=20&sort_by=event_time&sort_order=desc
```

**Respuesta exitosa (200):**
```json
{
  "data": [
    {
      "id": "uuid",
      "ticker": "AAPL",
      "company": "Apple Inc.",
      "brokerage": "Goldman Sachs",
      "action": "upgraded by",
      "rating_from": "buy",
      "rating_to": "strong buy",
      "target_from": 150.0,
      "target_to": 180.0,
      "event_time": "2024-01-01T10:00:00Z",
      "price_close": 155.5,
      "created_at": "2024-01-01T10:05:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total_pages": 10,
    "total_items": 200,
    "has_next": true,
    "has_prev": false
  }
}
```

**Funcionamiento interno:**
- Aplica filtros de búsqueda en la base de datos
- Implementa paginación eficiente
- Ordena resultados según parámetros
- Retorna metadatos de paginación

### 2. Obtener Stocks por Ticker
```
GET /api/v1/stocks/{ticker}
```

**Descripción:** Obtiene todos los eventos relacionados con un ticker específico

**Autenticación:** Opcional

**Path Parameters:**
- `ticker` (string, required): Símbolo de la acción (ej: AAPL, MSFT)

**Ejemplo:**
```
GET /api/v1/stocks/AAPL
```

**Respuesta exitosa (200):**
```json
{
  "data": [
    {
      "id": "uuid",
      "ticker": "AAPL",
      "company": "Apple Inc.",
      "brokerage": "Goldman Sachs",
      "action": "upgraded by",
      "rating_from": "buy",
      "rating_to": "strong buy",
      "target_from": 150.0,
      "target_to": 180.0,
      "event_time": "2024-01-01T10:00:00Z",
      "created_at": "2024-01-01T10:05:00Z"
    }
  ]
}
```

**Funcionamiento interno:**
- Busca en la base de datos todos los eventos para el ticker específico
- Ordena por fecha del evento (más recientes primero)
- Incluye historial completo de cambios de rating y precios objetivo

### 3. Obtener Estadísticas
```
GET /api/v1/stocks/stats
```

**Descripción:** Proporciona estadísticas básicas sobre los datos de stocks

**Autenticación:** Opcional

**Respuesta exitosa (200):**
```json
{
  "data": {
    "total_stocks": 15000,
    "last_updated": "2024-01-01T12:00:00Z"
  }
}
```

**Funcionamiento interno:**
- Cuenta el número total de eventos de stocks en la base de datos
- Proporciona timestamp de última actualización
- Optimizado para ser rápido usando solo metadatos

### 4. Endpoints Protegidos (No Implementados)
Los siguientes endpoints están definidos pero retornan "Not implemented":

- `GET /api/v1/stocks/{id}` - Obtener stock por ID
- `POST /api/v1/stocks/` - Crear nuevo stock (requiere autenticación)
- `PUT /api/v1/stocks/{id}` - Actualizar stock (requiere autenticación)
- `DELETE /api/v1/stocks/{id}` - Eliminar stock (requiere autenticación)

## Endpoints de Recomendaciones

El sistema de recomendaciones es el core de la aplicación, con funcionalidad tiered basada en el nivel del usuario.

### 1. Obtener Recomendaciones
```
GET /api/v1/recommendations
```

**Descripción:** Obtiene recomendaciones de acciones personalizadas según el nivel del usuario

**Autenticación:** Opcional (funcionalidad limitada para guests)

**Query Parameters:**
- `limit` (int): Número máximo de recomendaciones (default: 10)
  - Guest: máximo 10
  - Basic: máximo 25
  - Premium: máximo 100
- `min_score` (float): Score mínimo de recomendación (0.0-1.0)
- `type` (string): Tipo de recomendación ("strong_buy", "buy", "hold", "sell", "strong_sell")
- `exclude` (string): Lista de tickers a excluir (separados por coma)

**Ejemplo:**
```
GET /api/v1/recommendations?limit=20&min_score=0.7&type=buy&exclude=AAPL,MSFT
```

**Respuesta para Usuario GUEST (200):**
```json
{
  "data": [
    {
      "id": "uuid",
      "ticker": "AAPL",
      "company_name": "Apple Inc.",
      "total_events": 15,
      "positive_events": 12,
      "negative_events": 3,
      "avg_target_change": 15.5,
      "latest_target_price": 180.0,
      "broker_consensus": "Strong Buy",
      "basic_score": 0.85,
      "confidence": 0.92,
      "recommendation_type": "Strong Buy",
      "scoring_factors": [
        {
          "factor": "price_target_increase",
          "weight": 0.3,
          "score": 0.9,
          "description": "Consistent price target increases"
        }
      ],
      "tier": "basic",
      "last_event_time": "2024-01-01T10:00:00Z",
      "created_at": "2024-01-01T12:00:00Z",
      "expires_at": "2024-01-02T12:00:00Z"
    }
  ],
  "meta": {
    "count": 10,
    "user_tier": "guest",
    "features": ["basic_recommendations", "market_analytics"],
    "cache_hit": false,
    "generation_time": 250
  }
}
```

**Respuesta para Usuario BASIC (200):**
```json
{
  "data": [
    {
      // ... campos básicos igual que guest
      "external_data": {
        "current_price": 175.50,
        "price_change_24h": 2.5,
        "volume": 85000000,
        "market_cap": 2800000000000,
        "pe_ratio": 25.5,
        "analyst_ratings": {
          "strong_buy": 15,
          "buy": 8,
          "hold": 3,
          "sell": 1,
          "strong_sell": 0
        }
      }
    }
  ],
  "meta": {
    "count": 20,
    "user_tier": "basic",
    "features": ["basic_recommendations", "market_analytics", "real_time_data", "external_apis"],
    "cache_hit": true,
    "generation_time": 150
  }
}
```

**Respuesta para Usuario PREMIUM (200):**
```json
{
  "data": [
    {
      // ... campos básicos + external_data
      "ai_insights": {
        "sentiment_score": 0.75,
        "news_sentiment": "positive",
        "social_media_buzz": 0.82,
        "technical_indicators": {
          "rsi": 45.2,
          "macd": "bullish",
          "moving_averages": "positive_crossover"
        },
        "ai_prediction": "Strong momentum expected for next 30 days",
        "risk_assessment": "moderate"
      }
    }
  ],
  "meta": {
    "count": 50,
    "user_tier": "premium",
    "features": ["basic_recommendations", "market_analytics", "real_time_data", "external_apis", "ai_insights", "sentiment_analysis"],
    "cache_hit": false,
    "generation_time": 500,
    "rate_limit_remaining": 1850
  }
}
```

**Funcionamiento interno:**
- Analiza eventos de stocks de los últimos 30 días
- Calcula scores basados en múltiples factores:
  - Cambios en price targets
  - Mejoras en ratings
  - Consenso de brokers
  - Frecuencia de eventos positivos
- Aplica enrichment según tier del usuario:
  - Guest: Solo datos básicos
  - Basic: Datos básicos + APIs externas (Yahoo Finance, Alpha Vantage)
  - Premium: Todo lo anterior + AI insights
- Implementa caché inteligente con TTL diferenciado por tier
- Procesamiento concurrente optimizado con semáforos

### 2. Obtener Recomendación por Ticker
```
GET /api/v1/recommendations/{ticker}
```

**Descripción:** Obtiene una recomendación específica para un ticker

**Autenticación:** Opcional

**Path Parameters:**
- `ticker` (string, required): Símbolo de la acción (1-10 caracteres, letras y números)

**Ejemplo:**
```
GET /api/v1/recommendations/AAPL
```

**Respuesta exitosa (200):**
```json
{
  "data": {
    // ... estructura similar al endpoint anterior pero para un solo ticker
  },
  "meta": {
    "user_tier": "basic",
    "features": ["basic_recommendations", "market_analytics", "real_time_data"]
  }
}
```

**Funcionamiento interno:**
- Valida y normaliza el ticker (uppercase, caracteres válidos)
- Busca en caché primero
- Calcula recomendación usando datos históricos del ticker
- Aplica enrichment según tier del usuario
- Cachea resultado con TTL apropiado

### 3. Vista Previa de Recomendación
```
GET /api/v1/recommendations/preview/{ticker}
```

**Descripción:** Muestra qué datos adicionales estarían disponibles con un tier superior

**Autenticación:** Requerida (no disponible para guests)

**Path Parameters:**
- `ticker` (string, required): Símbolo de la acción

**Ejemplo:**
```
GET /api/v1/recommendations/preview/AAPL
```

**Respuesta para Usuario BASIC (200):**
```json
{
  "current_tier": {
    "tier": "basic",
    "data": {
      // ... datos disponibles en tier básico
    }
  },
  "premium_preview": {
    // ... muestra los ai_insights que tendría con premium
    "ai_insights": {
      "sentiment_score": 0.75,
      "ai_prediction": "Strong momentum expected...",
      "technical_indicators": { "rsi": 45.2 }
    }
  },
  "upgrade_benefits": [
    "Upgrade to Premium for AI-powered insights",
    "Get up to 100 recommendations",
    "Access to sentiment analysis",
    "Advanced market predictions"
  ]
}
```

**Funcionamiento interno:**
- Solo permite acceso a usuarios autenticados
- Obtiene recomendación para tier actual del usuario
- Si es tier BASIC, muestra preview de tier PREMIUM
- Proporciona lista de beneficios de upgrade

## Endpoints de Suscripciones

### 1. Crear Suscripción
```
POST /api/v1/subscriptions
```

**Descripción:** Crea una nueva suscripción para el usuario autenticado

**Autenticación:** Requerida

**Body (JSON):**
```json
{
  "plan": "monthly"
}
```

**Planes disponibles:**
- `"monthly"` - Suscripción mensual
- `"yearly"` - Suscripción anual

**Respuesta exitosa (201):**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "plan": "monthly",
  "status": "pending",
  "price": 29.99,
  "currency": "USD",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-02-01T00:00:00Z",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Funcionamiento interno:**
- Valida que el usuario exista
- Verifica que no tenga una suscripción activa
- Crea suscripción con status "pending"
- Calcula fecha de fin según el plan
- Asigna precio correspondiente

### 2. Simular Pago
```
POST /api/v1/subscriptions/{id}/payment
```

**Descripción:** Simula el procesamiento de pago para una suscripción

**Autenticación:** Requerida

**Path Parameters:**
- `id` (uuid, required): ID de la suscripción

**Ejemplo:**
```
POST /api/v1/subscriptions/550e8400-e29b-41d4-a716-446655440000/payment
```

**Respuesta exitosa (200):**
```json
{
  "message": "Payment processed successfully"
}
```

**Funcionamiento interno:**
- Obtiene la suscripción por ID
- Valida que esté en status "pending"
- Simula procesamiento de pago (delay de 2 segundos)
- Activa la suscripción con referencia de pago
- Actualiza el tier del usuario a "premium"

## Endpoints de Salud

### 1. Health Check
```
GET /health
```

**Descripción:** Verifica el estado de salud de la aplicación y conectividad a la base de datos

**Autenticación:** No requerida

**Respuesta exitosa (200):**
```json
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

**Respuesta con problemas (503):**
```json
{
  "status": "unhealthy",
  "database": "disconnected"
}
```

### 2. Ping
```
GET /ping
```

**Descripción:** Endpoint simple de ping para verificar que el servidor responde

**Autenticación:** No requerida

**Respuesta:** 200 OK

## Modelos de Datos

### User
```json
{
  "id": "uuid",
  "email": "string",
  "first_name": "string",
  "last_name": "string",
  "tier": "guest|basic|premium",
  "is_verified": "boolean",
  "last_login": "timestamp|null",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Stock Event
```json
{
  "id": "uuid",
  "ticker": "string",
  "company": "string",
  "brokerage": "string",
  "action": "string",
  "rating_from": "string",
  "rating_to": "string",
  "target_from": "float",
  "target_to": "float",
  "event_time": "timestamp",
  "price_close": "float|null",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Recommendation
```json
{
  "id": "uuid",
  "ticker": "string",
  "company_name": "string",
  "total_events": "int",
  "positive_events": "int",
  "negative_events": "int",
  "avg_target_change": "float",
  "latest_target_price": "float",
  "broker_consensus": "string",
  "basic_score": "float",
  "confidence": "float",
  "recommendation_type": "Strong Buy|Buy|Hold|Sell|Strong Sell",
  "scoring_factors": "array",
  "tier": "basic|enriched|premium",
  "external_data": "object|null",
  "ai_insights": "object|null",
  "last_event_time": "timestamp",
  "created_at": "timestamp",
  "expires_at": "timestamp"
}
```

### Subscription
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "plan": "monthly|yearly",
  "status": "pending|active|cancelled|expired",
  "price": "float",
  "currency": "string",
  "start_date": "timestamp",
  "end_date": "timestamp",
  "payment_reference": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

## Filtros y Paginación

### Filtros de Stocks
- `ticker`: Filtro por símbolo
- `company`: Filtro por nombre de empresa
- `brokerage`: Filtro por casa de corretaje
- `action`: Filtro por tipo de acción
- `sort_by`: Campo de ordenamiento (default: "event_time")
- `sort_order`: Orden "asc" o "desc" (default: "desc")

### Paginación
- `limit`: Elementos por página (default: 50, max: 1000)
- `offset`: Elementos a omitir
- Respuesta incluye metadatos de paginación

### Filtros de Recomendaciones
- `limit`: Número de recomendaciones (limitado por tier)
- `min_score`: Score mínimo (0.0-1.0)
- `type`: Tipo de recomendación
- `exclude`: Tickers a excluir (separados por coma)

## Rate Limiting

**Límites por hora por tier:**
- Guest: 100 requests/hora
- Basic: 500 requests/hora
- Premium: 2000 requests/hora

**Headers de respuesta:**
- `X-RateLimit-Remaining`: Requests restantes
- `X-RateLimit-Reset`: Timestamp de reset

**Middleware aplicado a:**
- `/api/v1/auth/*` - Rate limiting para autenticación
- `/api/v1/stocks/*` - Rate limiting basado en tier
- `/api/v1/recommendations/*` - Rate limiting basado en tier
- `/api/v1/subscriptions/*` - Rate limiting para usuarios autenticados

## Códigos de Error

### Códigos HTTP Comunes
- `200 OK` - Solicitud exitosa
- `201 Created` - Recurso creado exitosamente
- `400 Bad Request` - Error en los datos enviados
- `401 Unauthorized` - Autenticación requerida o inválida
- `403 Forbidden` - Sin permisos para acceder al recurso
- `404 Not Found` - Recurso no encontrado
- `429 Too Many Requests` - Rate limit excedido
- `500 Internal Server Error` - Error interno del servidor
- `503 Service Unavailable` - Servicio no disponible

### Formato de Errores
```json
{
  "error": "Descripción del error"
}
```

### Errores Específicos de Validación
```json
{
  "error": "Validation failed: password must be at least 8 characters long"
}
```

## Middleware y Seguridad

### CORS
La API incluye headers CORS para permitir acceso desde diferentes dominios:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Accept, Authorization, Content-Type, X-CSRF-Token`

### Middleware de Producción
- **Compression:** Compresión de respuestas
- **Request ID:** Tracking de requests
- **Recovery:** Recuperación de panics
- **Timeout:** Timeout de 60 segundos
- **Clean Path:** Limpieza de URLs

### Autenticación Middleware
- **OptionalAuth:** Permite acceso con o sin autenticación
- **RequireAuth:** Requiere autenticación válida
- **RequirePremium:** Requiere suscripción premium

## Caché

El sistema implementa caché inteligente:

### TTL por Tier
- Guest: 12 horas (datos básicos cambian menos)
- Basic: 4 horas (datos externos se actualizan más frecuentemente)
- Premium: 2 horas (usuarios premium obtienen datos más frescos)

### Keys de Caché
- Recomendaciones generales: `recommendations:{tier}:{limit}`
- Recomendación específica: `recommendation:{ticker}:{tier}`

## Consideraciones para Frontend

### Autenticación
1. Implementar login/registro con manejo de tokens
2. Almacenar tokens de forma segura (httpOnly cookies recomendado)
3. Implementar refresh automático de tokens
4. Manejar estados de autenticación (guest, basic, premium)

### Tiers de Usuario
1. Mostrar diferentes interfaces según el tier
2. Implementar preview de funcionalidades premium
3. Botones de upgrade para usuarios basic
4. Indicadores visuales de limitaciones

### Recomendaciones
1. Interfaz diferenciada por tier de usuario
2. Filtros avanzados para búsqueda
3. Visualización de scores y confianza
4. Tooltips explicativos para factores de scoring

### Rate Limiting
1. Mostrar contadores de requests restantes
2. Implementar retry automático cuando sea apropiado
3. Alertas cuando se acerque al límite

### Gestión de Suscripciones
1. Flujo de pago (integración con Stripe recomendada)
2. Gestión de estado de suscripción
3. Opciones de cancelación y renovación
4. Facturación y historial de pagos

### Performance
1. Implementar loading states apropiados
2. Caché local para datos frecuentemente accedidos
3. Paginación infinite scroll para listas largas
4. Debouncing para filtros de búsqueda

Esta documentación proporciona una base sólida para estructurar el frontend, entendiendo claramente las capacidades y limitaciones de cada endpoint y nivel de usuario. 