-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id uuid NOT NULL,
	mail text NOT NULL,
	"password" text NOT NULL,
	CONSTRAINT uni_users_mail UNIQUE (mail),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);
