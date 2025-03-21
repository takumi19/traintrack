create table program_workout_exercises (
    id                bigserial primary key,
    program_workout_id bigint not null references program_workouts(id) on delete cascade,
    exercise_id       bigint not null references exercises(id),
    order_index       integer not null default 1,  -- the order in which exercises appear
    notes             text,
    created_at        timestamp with time zone not null default now(),
    updated_at        timestamp with time zone not null default now()
);
