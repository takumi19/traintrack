create table users (
    id bigserial primary key,
    full_name text not null,
    -- name text not null,
    -- surname text not null,
    login text unique not null,
    email text unique not null,
    password_hash text not null
);

-- create trigger on_user_update
-- before update on users
-- for each row execute procedure timestamp_update();
