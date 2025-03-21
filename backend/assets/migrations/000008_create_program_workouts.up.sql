create table program_workouts (
    id             bigserial primary key,
    program_week_id bigint not null references program_weeks(id) on delete cascade,
    workout_index  integer not null,
    title          varchar(255),
    notes          text,
    created_at     timestamp with time zone not null default now(),
    updated_at     timestamp with time zone not null default now()
);
