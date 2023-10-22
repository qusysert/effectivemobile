CREATE TABLE IF NOT EXISTS "user"
(
    id         serial CONSTRAINT user_pk PRIMARY KEY,
    name       text NOT NULL,
    surname    text NOT NULL,
    patronymic text,
    age        int,
    gender     text,
    nation     text,
    full_name  text
);

CREATE OR REPLACE FUNCTION update_full_name()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.full_name := CONCAT(NEW.name, ' ', NEW.surname, ' ', NEW.patronymic);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_full_name_trigger
    BEFORE INSERT OR UPDATE OF name, surname, patronymic
    ON "user"
    FOR EACH ROW
EXECUTE FUNCTION update_full_name();
