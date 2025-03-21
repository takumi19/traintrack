create table program_weeks (
    id                 bigserial primary key,
    program_template_id bigint not null references program_templates(id) on delete cascade,
    week_number        integer not null,    -- week number in its program
    notes              text,
    created_at         timestamp with time zone not null default now(),
    updated_at         timestamp with time zone not null default now()
);
