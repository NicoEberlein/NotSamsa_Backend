-- Down-Migration

-- Drop Foreign Key Constraints
ALTER TABLE public.collection_participants DROP CONSTRAINT fk_collection_participants_user;
ALTER TABLE public.collection_participants DROP CONSTRAINT fk_collection_participants_collection;

-- Drop Table
DROP TABLE IF EXISTS public.collection_participants;
