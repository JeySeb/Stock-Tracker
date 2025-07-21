-- Tabla de brokers para normalizar
CREATE TABLE IF NOT EXISTS brokers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name STRING NOT NULL UNIQUE,
    credibility_score DECIMAL(3,2) DEFAULT 0.60,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Tabla principal de stocks
CREATE TABLE IF NOT EXISTS stocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker STRING NOT NULL,
    company STRING NOT NULL,
    broker_id UUID REFERENCES brokers(id),
    action STRING NOT NULL,
    rating_from STRING,
    rating_to STRING,
    target_from DECIMAL(10,2),
    target_to DECIMAL(10,2),
    event_time TIMESTAMPTZ NOT NULL,
    price_close DECIMAL(10,2), -- Para enriquecimiento futuro
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    
    INDEX idx_ticker (ticker),
    INDEX idx_event_time (event_time DESC),
    INDEX idx_ticker_time (ticker, event_time DESC),
    INDEX idx_broker_id (broker_id)
);

-- Tabla de recomendaciones generadas
CREATE TABLE IF NOT EXISTS recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticker STRING NOT NULL,
    score DECIMAL(5,4) NOT NULL CHECK (score >= 0 AND score <= 1),
    confidence DECIMAL(5,4) NOT NULL CHECK (confidence >= 0 AND confidence <= 1),
    factors JSONB,
    recommendation_type STRING NOT NULL,
    explanation TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    expires_at TIMESTAMPTZ,
    
    INDEX idx_ticker (ticker),
    INDEX idx_score (score DESC),
    INDEX idx_created_at (created_at DESC),
    INDEX idx_expires_at (expires_at)
);

-- Tabla para tracking de ingestion
CREATE TABLE IF NOT EXISTS ingestion_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_id STRING NOT NULL,
    total_records INT NOT NULL,
    successful_records INT NOT NULL,
    failed_records INT NOT NULL,
    status STRING NOT NULL, -- 'running', 'completed', 'failed'
    error_details JSONB,
    started_at TIMESTAMPTZ DEFAULT now(),
    completed_at TIMESTAMPTZ,
    
    INDEX idx_batch_id (batch_id),
    INDEX idx_status (status),
    INDEX idx_started_at (started_at DESC)
);

