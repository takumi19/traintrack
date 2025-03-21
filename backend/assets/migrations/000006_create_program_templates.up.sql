create table program_templates (
    id              bigserial primary key,
    author_id       bigint not null references users(id), -- on delete cascade?
    name            text not null,
    notes           text,
    created_at      timestamp with time zone not null default now(),
    updated_at      timestamp with time zone not null default now()
);
