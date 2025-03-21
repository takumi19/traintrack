create table workout_sets (
    id                   bigserial primary key,
    workout_exercise_id  bigint not null references logged_workout_exercises(id) on delete cascade,
    set_number           integer not null,          -- order
    reps                 integer,                   -- null if isometric
    weight               numeric(6,2),              -- null if bodyweight
    duration_sec         integer,                   -- null if rep-based
    rpe                  int,
    notes                text,
    created_at           timestamp with time zone not null default now(),
    updated_at           timestamp with time zone not null default now()
);
