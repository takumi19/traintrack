begin;
create type muscle_group as enum (
    'chest',
    'back',
    'shoulders',
    'biceps',
    'triceps',
    'forearms',
    'legs',
    'core',
    'neck',
    'other'
);

create table exercises (
    id            bigserial primary key,
    name          text unique not null,
    notes         text,
    is_rep_based  boolean default false not null, -- time-trackable/rep-trackable
    is_bodyweight boolean default false not null  -- bodyweight/external resistance
);

create table exercises_muscle_groups (
    exercise_id bigint not null,
    worked_muscle_group muscle_group not null,
    primary key (exercise_id, worked_muscle_group),  -- composite key
    foreign key (exercise_id) references exercises(id) on delete cascade
    -- no fk on muscle_group since it's an enum defined locally
);
commit;
