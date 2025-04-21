-- public.collection_participants definition

-- Drop table

-- DROP TABLE public.collection_participants;

CREATE TABLE public.collection_participants (
	collection_id uuid NOT NULL,
	user_id uuid NOT NULL,
	CONSTRAINT collection_participants_pkey PRIMARY KEY (collection_id, user_id)
);


-- public.collection_participants foreign keys

ALTER TABLE public.collection_participants ADD CONSTRAINT fk_collection_participants_collection FOREIGN KEY (collection_id) REFERENCES public.collections(id);
ALTER TABLE public.collection_participants ADD CONSTRAINT fk_collection_participants_user FOREIGN KEY (user_id) REFERENCES public.users(id);
