ALTER TABLE public.collections
DROP CONSTRAINT IF EXISTS fk_collections_preview_image;

ALTER TABLE public.collections
DROP COLUMN IF EXISTS preview_image_id;
