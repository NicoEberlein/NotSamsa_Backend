-- public.collections definition

-- Drop table

-- DROP TABLE public.collections;

CREATE TABLE public.collections (
	id uuid NOT NULL,
	owner_id uuid NULL,
	"name" text NULL,
	CONSTRAINT collections_pkey PRIMARY KEY (id)
);


-- public.collections foreign keys

ALTER TABLE public.collections ADD CONSTRAINT fk_collections_owner FOREIGN KEY (owner_id) REFERENCES public.users(id);
