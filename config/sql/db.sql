--
-- PostgreSQL database dump
--

-- Dumped from database version 10.12 (Ubuntu 10.12-0ubuntu0.18.04.1)
-- Dumped by pg_dump version 10.12 (Ubuntu 10.12-0ubuntu0.18.04.1)

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

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: organization; Type: TABLE; Schema: public; Owner: huvalk
--

CREATE TABLE public.organization (
    id integer NOT NULL,
    name character varying(60) NOT NULL,
    site character varying(60)
);


ALTER TABLE public.organization OWNER TO huvalk;

--
-- Name: organization_id_seq; Type: SEQUENCE; Schema: public; Owner: huvalk
--

CREATE SEQUENCE public.organization_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.organization_id_seq OWNER TO huvalk;

--
-- Name: organization_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: huvalk
--

ALTER SEQUENCE public.organization_id_seq OWNED BY public.organization.id;


--
-- Name: person; Type: TABLE; Schema: public; Owner: huvalk
--

CREATE TABLE public.person (
    id integer NOT NULL,
    name character varying(60) NOT NULL,
    gender character varying(10),
    birthday date
);


ALTER TABLE public.person OWNER TO huvalk;

--
-- Name: person_id_seq; Type: SEQUENCE; Schema: public; Owner: huvalk
--

CREATE SEQUENCE public.person_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.person_id_seq OWNER TO huvalk;

--
-- Name: person_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: huvalk
--

ALTER SEQUENCE public.person_id_seq OWNED BY public.person.id;


--
-- Name: session; Type: TABLE; Schema: public; Owner: huvalk
--

CREATE TABLE public.session (
    user_id integer,
    session_id character varying(64) NOT NULL,
    expires timestamp without time zone
);


ALTER TABLE public.session OWNER TO huvalk;

--
-- Name: users; Type: TABLE; Schema: public; Owner: huvalk
--

CREATE TABLE public.users (
    id integer NOT NULL,
    login character varying(60) NOT NULL,
    password character varying(60) NOT NULL,
    organization integer,
    person integer,
    tag character varying(20),
    email character varying(60),
    phone character varying(20),
    registrated timestamp without time zone,
    CONSTRAINT entity_not_empty CHECK (((organization IS NOT NULL) OR (person IS NOT NULL)))
);


ALTER TABLE public.users OWNER TO huvalk;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: huvalk
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO huvalk;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: huvalk
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: organization id; Type: DEFAULT; Schema: public; Owner: huvalk
--

ALTER TABLE ONLY public.organization ALTER COLUMN id SET DEFAULT nextval('public.organization_id_seq'::regclass);


--
-- Name: person id; Type: DEFAULT; Schema: public; Owner: huvalk
--

ALTER TABLE ONLY public.person ALTER COLUMN id SET DEFAULT nextval('public.person_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: huvalk
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: organization; Type: TABLE DATA; Schema: public; Owner: huvalk
--

COPY public.organization (id, name, site) FROM stdin;
\.


--
-- Data for Name: person; Type: TABLE DATA; Schema: public; Owner: huvalk
--

COPY public.person (id, name, gender, birthday) FROM stdin;
\.


--
-- Data for Name: session; Type: TABLE DATA; Schema: public; Owner: huvalk
--

COPY public.session (user_id, session_id, expires) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: huvalk
--

COPY public.users (id, login, password, organization, person, tag, email, phone, registrated) FROM stdin;
\.


--
-- Name: organization_id_seq; Type: SEQUENCE SET; Schema: public; Owner: huvalk
--

SELECT pg_catalog.setval('public.organization_id_seq', 1, false);


--
-- Name: person_id_seq; Type: SEQUENCE SET; Schema: public; Owner: huvalk
--

SELECT pg_catalog.setval('public.person_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: huvalk
--

SELECT pg_catalog.setval('public.users_id_seq', 1, false);


--
-- Name: organization organization_pkey; Type: CONSTRAINT; Schema: public; Owner: huvalk
--

ALTER TABLE ONLY public.organization
    ADD CONSTRAINT organization_pkey PRIMARY KEY (id);


--
-- Name: person person_pkey; Type: CONSTRAINT; Schema: public; Owner: huvalk
--

ALTER TABLE ONLY public.person
    ADD CONSTRAINT person_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: huvalk
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: session session_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: huvalk
--

ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: users users_organization_fkey; Type: FK CONSTRAINT; Schema: public; Owner: huvalk
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_organization_fkey FOREIGN KEY (organization) REFERENCES public.organization(id);


--
-- Name: users users_person_fkey; Type: FK CONSTRAINT; Schema: public; Owner: huvalk
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_person_fkey FOREIGN KEY (person) REFERENCES public.person(id);


--
-- PostgreSQL database dump complete
--

