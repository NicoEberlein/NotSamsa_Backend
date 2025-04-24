ALTER TABLE public.collections
ADD COLUMN description TEXT,
ADD COLUMN latitude NUMERIC(10, 7),
ADD COLUMN longitude NUMERIC(10, 7);
