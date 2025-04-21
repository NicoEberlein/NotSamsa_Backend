-- Down-Migration

-- Drop Foreign Key Constraint
ALTER TABLE public.collections DROP CONSTRAINT fk_collections_owner;

-- Drop Table
DROP TABLE IF EXISTS public.collections;
