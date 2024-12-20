-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.statistics
(
    id uuid NOT NULL,
    ip text COLLATE pg_catalog."default" NOT NULL,
    name text COLLATE pg_catalog."default" DEFAULT ''::text,
    time_start timestamp NOT NULL,
    time_end timestamp,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT statistics_pkey PRIMARY KEY (id)
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.statistics
    OWNER to postgres;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.statistics;
-- +goose StatementEnd
