--
-- PostgreSQL database dump
--

-- Dumped from database version 10.5 (Debian 10.5-2.pgdg90+1)
-- Dumped by pg_dump version 10.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE IF EXISTS forex;
--
-- Name: forex; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE forex WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';


ALTER DATABASE forex OWNER TO postgres;

\connect forex

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
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
-- Name: exchange_rate_datas; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.exchange_rate_datas (
    id integer NOT NULL,
    exchange_rate_id integer NOT NULL,
    rate numeric NOT NULL,
    valid_time date NOT NULL
);


ALTER TABLE public.exchange_rate_datas OWNER TO postgres;

--
-- Name: exchange_rate_datas_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.exchange_rate_datas_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.exchange_rate_datas_id_seq OWNER TO postgres;

--
-- Name: exchange_rate_datas_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.exchange_rate_datas_id_seq OWNED BY public.exchange_rate_datas.id;


--
-- Name: exchange_rates; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.exchange_rates (
    id integer NOT NULL,
    currency_from character varying(3) NOT NULL,
    currency_to character varying(3) NOT NULL
);


ALTER TABLE public.exchange_rates OWNER TO postgres;

--
-- Name: TABLE exchange_rates; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.exchange_rates IS 'Exchange Rate List';


--
-- Name: exchange_rates_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.exchange_rates_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.exchange_rates_id_seq OWNER TO postgres;

--
-- Name: exchange_rates_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.exchange_rates_id_seq OWNED BY public.exchange_rates.id;


--
-- Name: exchange_rate_datas id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.exchange_rate_datas ALTER COLUMN id SET DEFAULT nextval('public.exchange_rate_datas_id_seq'::regclass);


--
-- Name: exchange_rates id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.exchange_rates ALTER COLUMN id SET DEFAULT nextval('public.exchange_rates_id_seq'::regclass);


--
-- Data for Name: exchange_rate_datas; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.exchange_rate_datas (id, exchange_rate_id, rate, valid_time) FROM stdin;
1	2	0.5	2019-03-18
2	2	0.12	2019-03-17
3	2	0.701	2019-03-16
4	2	0.70001	2019-03-15
5	2	0.699	2019-03-14
6	2	0.55	2019-03-13
7	2	0.5	2019-03-12
8	2	0.5	2019-03-11
9	3	0.92	2019-03-18
10	1	0.22	2019-03-18
11	1	0.2775	2019-03-15
\.


--
-- Data for Name: exchange_rates; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.exchange_rates (id, currency_from, currency_to) FROM stdin;
1	IDR	USD
2	IDR	EUR
3	IDR	BHT
\.


--
-- Name: exchange_rate_datas_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.exchange_rate_datas_id_seq', 11, true);


--
-- Name: exchange_rates_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.exchange_rates_id_seq', 3, true);


--
-- Name: exchange_rate_datas exchange_rate_datas_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.exchange_rate_datas
    ADD CONSTRAINT exchange_rate_datas_pkey PRIMARY KEY (id);


--
-- Name: exchange_rates exchange_rates_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.exchange_rates
    ADD CONSTRAINT exchange_rates_pkey PRIMARY KEY (id);


--
-- Name: exchange_rate_datas_exchange_rate_id_valid_time_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX exchange_rate_datas_exchange_rate_id_valid_time_uindex ON public.exchange_rate_datas USING btree (exchange_rate_id, valid_time);


--
-- Name: exchange_rates_currency_from_currency_to_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX exchange_rates_currency_from_currency_to_uindex ON public.exchange_rates USING btree (currency_from, currency_to);


--
-- Name: exchange_rate_datas exchange_rate_datas_exchange_rates_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.exchange_rate_datas
    ADD CONSTRAINT exchange_rate_datas_exchange_rates_id_fk FOREIGN KEY (exchange_rate_id) REFERENCES public.exchange_rates(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

