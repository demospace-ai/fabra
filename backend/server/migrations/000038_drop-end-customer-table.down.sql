CREATE TABLE IF NOT EXISTS end_customers (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX end_customers_organization_id_idx ON end_customers(organization_id);

ALTER TABLE sources ADD CONSTRAINT sources_end_customer_id_fkey FOREIGN KEY (end_customer_id) REFERENCES end_customers;