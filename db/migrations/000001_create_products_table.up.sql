CREATE TABLE IF NOT EXISTS public.products (
     id UUID PRIMARY KEY,
     name TEXT NOT NULL,
     description TEXT NOT NULL,
     price NUMERIC(12, 2) NOT NULL,
     created_at TIMESTAMP(3) WITH TIME ZONE NOT NULL,
     updated_at TIMESTAMP(3) WITH TIME ZONE NOT NULL
);