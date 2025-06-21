CREATE TABLE IF NOT EXISTS messages (
    id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    payload JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (tenant_id, id)
) PARTITION BY LIST (tenant_id);