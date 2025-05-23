-- public.images definition

-- Drop table

-- DROP TABLE public.images;

CREATE TABLE public.images (
	id uuid NOT NULL,
	collection_id uuid NULL,
	"path" text NULL,
	"name" text NULL,
	"size" int8 NULL,
	format text NULL,
	uploaded_at timestamptz NULL,
	CONSTRAINT images_pkey PRIMARY KEY (id)
);


-- public.images foreign keys

ALTER TABLE public.images ADD CONSTRAINT fk_collections_images FOREIGN KEY (collection_id) REFERENCES public.collections(id) ON DELETE CASCADE;
