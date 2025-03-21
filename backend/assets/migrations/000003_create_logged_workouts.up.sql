create table logged_workouts (
    id           bigserial primary key,
    user_id      bigint not null references users(id) on delete cascade,
    workout_date timestamp with time zone not null default now(),
    notes        text,
    created_at   timestamp with time zone not null default now(),
    updated_at   timestamp with time zone not null default now()
);
