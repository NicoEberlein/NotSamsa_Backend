-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id uuid NOT NULL,
	mail text NULL,
	"password" text NULL,
	CONSTRAINT uni_users_mail UNIQUE (mail),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);
