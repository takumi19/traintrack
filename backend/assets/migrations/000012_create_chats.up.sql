begin;
create table chats (
    id           bigserial primary key,
    name         text not null,
    created_at   timestamp with time zone not null default now(),
    updated_at   timestamp with time zone not null default now()
);

-- users to chats join table, no access rights
create table users_chats (
    id        bigserial primary key,
    chat_id   bigint references chats(id),
    user_id   bigint references users(id),
    unique    (user_id, chat_id)
);

create table chat_messages (
    id            bigserial primary key,
    author_id     bigint references users(id),
    text_content  text,
    img_content   text,
    sent_at       timestamp with time zone not null,
    edited_at     timestamp with time zone not null
);
commit;
