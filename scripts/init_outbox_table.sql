-- Criar tabela outbox
CREATE TYPE outbox_status as ENUM ('PENDING', 'CONSUMED');
CREATE TABLE IF NOT EXISTS outbox (
    id UUID PRIMARY KEY,
    domain_id UUID,
    paylaod JSONB NOT NULL,
    entity VARCHAR(40) NOT NULL,
    status outbox_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Índices para melhor performance
CREATE INDEX IF NOT EXISTS idx_outbox_status ON outbox(status);