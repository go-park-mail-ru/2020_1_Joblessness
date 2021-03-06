--
-- PostgreSQL database dump
--

-- Dumped from database version 12.3 (Ubuntu 12.3-1.pgdg18.04+1)
-- Dumped by pg_dump version 12.3 (Ubuntu 12.3-1.pgdg18.04+1)

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
-- Name: education; Type: TABLE; Schema: public; Owner: haha
--

CREATE TABLE public.education (
    summary_id integer NOT NULL,
    institution character varying(60),
    speciality character varying(60),
    graduated timestamp without time zone,
    type character varying(60)
);


ALTER TABLE public.education OWNER TO haha;

--
-- Name: experience; Type: TABLE; Schema: public; Owner: haha
--

CREATE TABLE public.experience (
    company_name character varying(60),
    role character varying(120),
    responsibilities text,
    start timestamp without time zone,
    stop timestamp without time zone,
    summary_id integer NOT NULL
);


ALTER TABLE public.experience OWNER TO haha;

--
-- Name: favorite; Type: TABLE; Schema: public; Owner: haha
--

CREATE TABLE public.favorite (
    user_id integer NOT NULL,
    favorite_id integer NOT NULL
);


ALTER TABLE public.favorite OWNER TO haha;

--
-- Name: message; Type: TABLE; Schema: public; Owner: haha
--

CREATE TABLE public.message (
    user_one_id integer NOT NULL,
    user_two_id integer,
    body text DEFAULT ''::text NOT NULL,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    user_one character varying(120) DEFAULT ''::character varying NOT NULL,
    user_two character varying(120) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.message OWNER TO haha;

--
-- Name: organization; Type: TABLE; Schema: public; Owner: haha
--

CREATE TABLE public.organization (
    id integer NOT NULL,
    name character varying(60) DEFAULT ''::character varying NOT NULL,
    site character varying(60) DEFAULT ''::character varying NOT NULL,
    about text DEFAULT ''::text
);


ALTER TABLE public.organization OWNER TO haha;

--
-- Name: organization_id_seq; Type: SEQUENCE; Schema: public; Owner: haha
--

CREATE SEQUENCE public.organization_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.organization_id_seq OWNER TO haha;

--
-- Name: organization_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: haha
--

ALTER SEQUENCE public.organization_id_seq OWNED BY public.organization.id;


--
-- Name: person; Type: TABLE; Schema: public; Owner: haha
--

CREATE TABLE public.person (
    id integer NOT NULL,
    name character varying(60) DEFAULT ''::character varying NOT NULL,
    gender character varying(10) DEFAULT ''::character varying NOT NULL,
    birthday date DEFAULT '0001-01-01'::date NOT NULL,
    surname character varying(60) DEFAULT ''::character varying NOT NULL
);


ALTER TABLE public.person OWNER TO haha;

--
-- Name: person_id_seq; Type: SEQUENCE; Schema: public; Owner: haha
--

CREATE SEQUENCE public.person_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.person_id_seq OWNER TO haha;

--
-- Name: person_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: haha
--

ALTER SEQUENCE public.person_id_seq OWNED BY public.person.id;


--
-- Name: requirement; Type: TABLE; Schema: public; Owner: haha
--

CREATE TABLE public.requirement (
    summary_id integer,
    driver_license character varying(60),
    has_car boolean DEFAULT false,
    schedule text,
    employment text
);


ALTER TABLE public.requirement OWNER TO haha;

--
-- Name: response; Type: TABLE; Schema: public; Owner: haha
--

CREATE TABLE public.response (
    summary_id integer,
    vacancy_id integer,
    rejected boolean DEFAULT false,
    approved boolean DEFAULT false,
    date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    interview_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT approved_rejected CHECK (((rejected = false) OR (approved = false))),
    CONSTRAINT has_reflection CHECK (((rejected IS NOT NULL) AND (approved IS NOT NULL)))
);


ALTER TABLE public.response OWNER TO haha;

--
-- Name: session; Type: TABLE; Schema: public; Owner: haha
--

CREATE TABLE public.session (
    user_id integer,
    session_id character varying(64) NOT NULL,
    expires timestamp without time zone
);


ALTER TABLE public.session OWNER TO haha;

--
-- Name: summary; Type: TABLE; Schema: public; Owner: haha
--

CREATE TABLE public.summary (
    id integer NOT NULL,
    author integer,
    keywords text,
    name character varying(120) DEFAULT ''::character varying,
    salary_from integer DEFAULT 0 NOT NULL,
    salary_to integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.summary OWNER TO haha;

--
-- Name: summary_id_seq; Type: SEQUENCE; Schema: public; Owner: haha
--

CREATE SEQUENCE public.summary_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.summary_id_seq OWNER TO haha;

--
-- Name: summary_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: haha
--

ALTER SEQUENCE public.summary_id_seq OWNED BY public.summary.id;


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
-- Name: vacancy; Type: TABLE; Schema: public; Owner: haha
--

CREATE TABLE public.vacancy (
    id integer NOT NULL,
    organization_id integer,
    name character varying(60) NOT NULL,
    description text,
    with_tax boolean DEFAULT false,
    responsibilities text,
    conditions text,
    keywords text,
    salary_from integer DEFAULT 0 NOT NULL,
    salary_to integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.vacancy OWNER TO haha;

--
-- Name: vacancy_id_seq; Type: SEQUENCE; Schema: public; Owner: haha
--

CREATE SEQUENCE public.vacancy_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.vacancy_id_seq OWNER TO haha;

--
-- Name: vacancy_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: haha
--

ALTER SEQUENCE public.vacancy_id_seq OWNED BY public.vacancy.id;


--
-- Name: organization id; Type: DEFAULT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.organization ALTER COLUMN id SET DEFAULT nextval('public.organization_id_seq'::regclass);


--
-- Name: person id; Type: DEFAULT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.person ALTER COLUMN id SET DEFAULT nextval('public.person_id_seq'::regclass);


--
-- Name: summary id; Type: DEFAULT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.summary ALTER COLUMN id SET DEFAULT nextval('public.summary_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: vacancy id; Type: DEFAULT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.vacancy ALTER COLUMN id SET DEFAULT nextval('public.vacancy_id_seq'::regclass);


--
-- Data for Name: education; Type: TABLE DATA; Schema: public; Owner: haha
--

COPY public.education (summary_id, institution, speciality, graduated, type) FROM stdin;
\.


--
-- Data for Name: experience; Type: TABLE DATA; Schema: public; Owner: haha
--

COPY public.experience (company_name, role, responsibilities, start, stop, summary_id) FROM stdin;
\.


--
-- Data for Name: favorite; Type: TABLE DATA; Schema: public; Owner: haha
--

COPY public.favorite (user_id, favorite_id) FROM stdin;
\.


--
-- Data for Name: message; Type: TABLE DATA; Schema: public; Owner: haha
--

COPY public.message (user_one_id, user_two_id, body, created, user_one, user_two) FROM stdin;
\.


--
-- Data for Name: organization; Type: TABLE DATA; Schema: public; Owner: haha
--

COPY public.organization (id, name, site, about) FROM stdin;
\.


--
-- Data for Name: person; Type: TABLE DATA; Schema: public; Owner: haha
--

COPY public.person (id, name, gender, birthday, surname) FROM stdin;
\.


--
-- Data for Name: requirement; Type: TABLE DATA; Schema: public; Owner: haha
--

COPY public.requirement (summary_id, driver_license, has_car, schedule, employment) FROM stdin;
\.


--
-- Data for Name: response; Type: TABLE DATA; Schema: public; Owner: haha
--

COPY public.response (summary_id, vacancy_id, rejected, approved, date, interview_date) FROM stdin;
\.


--
-- Data for Name: session; Type: TABLE DATA; Schema: public; Owner: haha
--

COPY public.session (user_id, session_id, expires) FROM stdin;
\.


--
-- Data for Name: summary; Type: TABLE DATA; Schema: public; Owner: haha
--

COPY public.summary (id, author, keywords, name, salary_from, salary_to) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: haha
--

COPY public.users (id, login, password, organization_id, person_id, tag, email, phone, registered, avatar) FROM stdin;
\.


--
-- Data for Name: vacancy; Type: TABLE DATA; Schema: public; Owner: haha
--

COPY public.vacancy (id, organization_id, name, description, with_tax, responsibilities, conditions, keywords, salary_from, salary_to) FROM stdin;
\.


--
-- Name: organization_id_seq; Type: SEQUENCE SET; Schema: public; Owner: haha
--

SELECT pg_catalog.setval('public.organization_id_seq', 4, true);


--
-- Name: person_id_seq; Type: SEQUENCE SET; Schema: public; Owner: haha
--

SELECT pg_catalog.setval('public.person_id_seq', 8, true);


--
-- Name: summary_id_seq; Type: SEQUENCE SET; Schema: public; Owner: haha
--

SELECT pg_catalog.setval('public.summary_id_seq', 2, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: haha
--

SELECT pg_catalog.setval('public.users_id_seq', 12, true);


--
-- Name: vacancy_id_seq; Type: SEQUENCE SET; Schema: public; Owner: haha
--

SELECT pg_catalog.setval('public.vacancy_id_seq', 13, true);


--
-- Name: organization organization_pkey; Type: CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.organization
    ADD CONSTRAINT organization_pkey PRIMARY KEY (id);


--
-- Name: person person_pkey; Type: CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.person
    ADD CONSTRAINT person_pkey PRIMARY KEY (id);


--
-- Name: summary summary_pkey; Type: CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.summary
    ADD CONSTRAINT summary_pkey PRIMARY KEY (id);


--
-- Name: favorite unique_like; Type: CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.favorite
    ADD CONSTRAINT unique_like UNIQUE (user_id, favorite_id);


--
-- Name: response unique_response; Type: CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.response
    ADD CONSTRAINT unique_response UNIQUE (summary_id, vacancy_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: vacancy vacancy_pkey; Type: CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.vacancy
    ADD CONSTRAINT vacancy_pkey PRIMARY KEY (id);


--
-- Name: session_session_id_uindex; Type: INDEX; Schema: public; Owner: haha
--

CREATE UNIQUE INDEX session_session_id_uindex ON public.session USING btree (session_id);


--
-- Name: education education_summary_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.education
    ADD CONSTRAINT education_summary_id_fkey FOREIGN KEY (summary_id) REFERENCES public.summary(id) ON DELETE CASCADE;


--
-- Name: experience experience_summary_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.experience
    ADD CONSTRAINT experience_summary_id_fkey FOREIGN KEY (summary_id) REFERENCES public.summary(id) ON DELETE CASCADE;


--
-- Name: favorite favorite_favorite_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.favorite
    ADD CONSTRAINT favorite_favorite_id_fkey FOREIGN KEY (favorite_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: favorite favorite_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.favorite
    ADD CONSTRAINT favorite_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: message message_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.message
    ADD CONSTRAINT message_users_id_fk FOREIGN KEY (user_one_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: message message_users_id_fk_2; Type: FK CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.message
    ADD CONSTRAINT message_users_id_fk_2 FOREIGN KEY (user_two_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: requirement requirement_summary_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.requirement
    ADD CONSTRAINT requirement_summary_id_fkey FOREIGN KEY (summary_id) REFERENCES public.summary(id) ON DELETE CASCADE;


--
-- Name: response response_summary_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.response
    ADD CONSTRAINT response_summary_id_fkey FOREIGN KEY (summary_id) REFERENCES public.summary(id) ON DELETE CASCADE;


--
-- Name: response response_vacancy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.response
    ADD CONSTRAINT response_vacancy_id_fkey FOREIGN KEY (vacancy_id) REFERENCES public.vacancy(id) ON DELETE CASCADE;


--
-- Name: session session_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: summary summary_author_fkey; Type: FK CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.summary
    ADD CONSTRAINT summary_author_fkey FOREIGN KEY (author) REFERENCES public.users(id) ON DELETE CASCADE;


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
-- Name: vacancy vacancy_organization_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: haha
--

ALTER TABLE ONLY public.vacancy
    ADD CONSTRAINT vacancy_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

