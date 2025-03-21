create table program_workout_sets (
    id                         bigserial primary key,
    -- NOTE: it may be slightly confusing that the field below does not refer to the `exercise_id` field of the corresponding exercise, but
    -- instead refers to its pk id. fix: rename the exercises table and the `exercise_id` field in the program_workout_exercises table.
    program_workout_exercise_id bigint not null references program_workout_exercises(id) on delete cascade,
    -- TODO: rename from number to index
    set_number                 integer not null default 1,
    -- TODO: add constraint: not less than 0 and not more than 10
    rpe                        integer,

    suggested_reps_min         integer,
    suggested_reps_max         integer,
    suggested_reps             integer, -- null if using a range

    suggested_weight_min       numeric(6,2),
    suggested_weight_max       numeric(6,2),
    suggested_weight           numeric(6,2), -- exact weight if no range

    suggested_time_min         integer, -- in seconds
    suggested_time_max         integer,
    suggested_time             integer,

    notes                      text,
    created_at                 timestamp with time zone not null default now(),
    updated_at                 timestamp with time zone not null default now()
);
