-- Down-Migration

-- Drop Foreign Key Constraint
ALTER TABLE public.images DROP CONSTRAINT fk_collections_images;

-- Drop Table
DROP TABLE IF EXISTS public.images;
