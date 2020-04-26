--
-- PostgreSQL database dump
--

-- Dumped from database version 12.2 (Ubuntu 12.2-2.pgdg18.04+1)
-- Dumped by pg_dump version 12.2 (Ubuntu 12.2-2.pgdg18.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE haha;
--
-- Name: haha; Type: DATABASE; Schema: -; Owner: haha
--

CREATE DATABASE haha WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'ru_RU.UTF-8' LC_CTYPE = 'ru_RU.UTF-8';


ALTER DATABASE haha OWNER TO haha;

\connect haha

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: users; Type: TABLE; Schema: public; Owner: haha
--

CREATE TABLE public.users (
    id integer NOT NULL,
    login character varying(60) NOT NULL,
    password character varying(60) NOT NULL,
    organization_id integer,
    person_id integer,
    tag character varying(20) DEFAULT ''::character varying,
    email character varying(60) DEFAULT ''::character varying NOT NULL,
    phone character varying(20) DEFAULT ''::character varying NOT NULL,
    registered timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    avatar text DEFAULT 'https://hb.bizmrg.com/imgs-hh/default-avatar.png'::text,
    CONSTRAINT entity_not_empty CHECK (((organization_id IS NOT NULL) OR (person_id IS NOT NULL)))
);


ALTER TABLE public.users OWNER TO haha;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: haha
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO haha;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: haha
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: haha
--

COPY public.users (id, login, password, organization_id, person_id, tag, email, phone, registered, avatar) FROM stdin;
\.


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: haha
--

SELECT pg_catalog.setval('public.users_id_seq', 1, false);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_organization_fkey; Type: FK CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_organization_fkey FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE;


--
-- Name: users users_person_fkey; Type: FK CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_person_fkey FOREIGN KEY (person_id) REFERENCES public.person(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

