create table logged_workout_exercises (
    id            bigserial primary key,
    workout_id    bigint not null references logged_workouts(id) on delete cascade,
    exercise_id   bigint references exercises(id),
    order_index   integer not null default 1,       -- the order in which exercises are performed
    notes         text,
    created_at    timestamp with time zone not null default now(),
    updated_at    timestamp with time zone not null default now()
);
