-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.addresses
(
    id uuid NOT NULL,
    ip text COLLATE pg_catalog."default" NOT NULL,
    name text COLLATE pg_catalog."default" DEFAULT ''::text,
    max_rtt integer DEFAULT 0,
    interval integer DEFAULT 100,
    count integer DEFAULT 5,
    timeout integer DEFAULT 1000,
    not_count integer DEFAULT 3,
    period_start integer DEFAULT 0,
    period_end integer DEFAULT 0,
    enabled boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT addresses_pkey PRIMARY KEY (id),
    UNIQUE(ip)
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.addresses
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.addresses;
-- +goose StatementEnd
