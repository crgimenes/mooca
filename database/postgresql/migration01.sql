CREATE TABLE IF NOT EXISTS structs
(
    id serial NOT NULL,
    parent integer NOT NULL DEFAULT 0,
    z_order integer NOT NULL DEFAULT 0,
    type integer NOT NULL,
    name character varying(254) NOT NULL,
    title character varying(254) NOT NULL,
    metadata jsonb,
    PRIMARY KEY (id),
    FOREIGN KEY (type) REFERENCES types(id)
);

CREATE TABLE IF NOT EXISTS types
(
    id serial NOT NULL,
    name character varying(254) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS datasets
(
    id serial NOT NULL,
    parent integer NOT NULL DEFAULT 0,
    struct_id integer NOT NULL,
    value text NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (struct_id) REFERENCES structs(id)
);
