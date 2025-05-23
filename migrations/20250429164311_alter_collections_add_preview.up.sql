ALTER TABLE public.collections
ADD COLUMN preview_image_id UUID NULL;

ALTER TABLE public.collections
ADD CONSTRAINT fk_collections_preview_image
FOREIGN KEY (preview_image_id) REFERENCES public.images(id) ON DELETE SET NULL;

