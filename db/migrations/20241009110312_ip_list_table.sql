-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.ip_list
(
    id uuid NOT NULL,
    ip text COLLATE pg_catalog."default" NOT NULL,
    name text COLLATE pg_catalog."default" DEFAULT ''::text,
    max_rtt text COLLATE pg_catalog."default" DEFAULT ''::text,
    interval text COLLATE pg_catalog."default" DEFAULT '0.1s'::text,
    count integer DEFAULT 5,
    timeout text COLLATE pg_catalog."default" DEFAULT '1s'::text,
    not_count integer DEFAULT 3,
    period_start text COLLATE pg_catalog."default" DEFAULT ''::text,
    period_end text COLLATE pg_catalog."default" DEFAULT ''::text,
    enabled boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT ap_list_pkey PRIMARY KEY (id),
    UNIQUE(ip)
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.ip_list
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.ip_list;
-- +goose StatementEnd
