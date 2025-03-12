-- TODO: create database users with different privileges

------------------------
--      TYPES
------------------------

CREATE TYPE muscle_group AS ENUM (
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

------------------------
--    END OF TYPES
------------------------

------------------------
--      MISC
------------------------

-- TODO: make id UUID
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,
    login TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE exercises (
    id            BIGSERIAL PRIMARY KEY,
    name          TEXT UNIQUE NOT NULL,
    notes         TEXT,
    is_rep_based  BOOLEAN DEFAULT FALSE NOT NULL, -- time-trackable/rep-trackable
    is_bodyweight BOOLEAN DEFAULT FALSE NOT NULL  -- bodyweight/external resistance
);

CREATE TABLE exercises_muscle_groups (
    exercise_id BIGINT NOT NULL,
    worked_muscle_group muscle_group NOT NULL,
    PRIMARY KEY (exercise_id, worked_muscle_group),  -- Composite key
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
    -- No FK on muscle_group since it's an ENUM defined locally
);

------------------------
--    END OF MISC
------------------------

------------------------
--      LOGGING
------------------------

CREATE TABLE logged_workouts (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workout_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    notes        TEXT,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE logged_workout_exercises (
    id            BIGSERIAL PRIMARY KEY,
    workout_id    BIGINT NOT NULL REFERENCES logged_workouts(id) ON DELETE CASCADE,
    exercise_id   BIGINT REFERENCES exercises(id),
    order_index   INTEGER NOT NULL DEFAULT 1,       -- the order in which exercises are performed
    notes         TEXT,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE workout_sets (
    id                   BIGSERIAL PRIMARY KEY,
    workout_exercise_id  BIGINT NOT NULL REFERENCES logged_workout_exercises(id) ON DELETE CASCADE,
    set_number           INTEGER NOT NULL,          -- order
    reps                 INTEGER,                   -- null if isometric
    weight               NUMERIC(6,2),              -- null if bodyweight
    duration_sec         INTEGER,                   -- null if rep-based
    rpe                  INT,
    notes                TEXT,
    created_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

------------------------
--   END OF LOGGING
------------------------

------------------------
--     TEMPLATING
------------------------

CREATE TABLE program_templates (
    id              BIGSERIAL PRIMARY KEY,
    author_id       BIGINT NOT NULL REFERENCES users(id), -- ON DELETE CASCADE?
    name            TEXT NOT NULL,
    notes           TEXT,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE program_weeks (
    id                 BIGSERIAL PRIMARY KEY,
    program_template_id BIGINT NOT NULL REFERENCES program_templates(id) ON DELETE CASCADE,
    week_number        INTEGER NOT NULL,    -- Week number in its program
    notes              TEXT,
    created_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE program_workouts (
    id             BIGSERIAL PRIMARY KEY,
    program_week_id BIGINT NOT NULL REFERENCES program_weeks(id) ON DELETE CASCADE,
    workout_index  INTEGER NOT NULL,
    title          VARCHAR(255),
    notes          TEXT,
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE program_workout_exercises (
    id                BIGSERIAL PRIMARY KEY,
    program_workout_id BIGINT NOT NULL REFERENCES program_workouts(id) ON DELETE CASCADE,
    exercise_id       BIGINT NOT NULL REFERENCES exercises(id),
    order_index       INTEGER NOT NULL DEFAULT 1,  -- the order in which exercises appear
    notes             TEXT,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE program_workout_sets (
    id                         BIGSERIAL PRIMARY KEY,
    -- NOTE: it may be slightly confusing that the field below does not refer to the `exercise_id` field of the corresponding exercise, but
    -- instead refers to its PK id. Fix: rename the exercises table and the `exercise_id` field in the program_workout_exercises table.
    program_workout_exercise_id BIGINT NOT NULL REFERENCES program_workout_exercises(id) ON DELETE CASCADE,
    -- TODO: Rename from number to index
    set_number                 INTEGER NOT NULL DEFAULT 1,
    rpe                        INTEGER,

    suggested_reps_min         INTEGER,
    suggested_reps_max         INTEGER,
    suggested_reps             INTEGER, -- null if using a range

    suggested_weight_min       NUMERIC(6,2),
    suggested_weight_max       NUMERIC(6,2),
    suggested_weight           NUMERIC(6,2), -- exact weight if no range

    suggested_time_min         INTEGER, -- in seconds
    suggested_time_max         INTEGER,
    suggested_time             INTEGER,

    notes                      TEXT,
    created_at                 TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at                 TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

------------------------
--   END OF TEMPLATING
------------------------

------------------------
--   ACCESS RIGHTS
------------------------

-- Defines the relationship between users and program templates
CREATE TABLE programs_permissions (
    permission_id SERIAL PRIMARY KEY,
    user_id       INTEGER NOT NULL,
    program_id    INTEGER NOT NULL,
    can_view      BOOLEAN DEFAULT FALSE NOT NULL,
    can_modify    BOOLEAN DEFAULT FALSE NOT NULL,

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id) ON DELETE CASCADE,  -- If user is deleted, their permissions are too.
    CONSTRAINT fk_program
        FOREIGN KEY(program_id)
        REFERENCES program_templates(id) ON DELETE CASCADE,
    UNIQUE (user_id, program_id)
);

------------------------
-- END OF ACCESS RIGHTS
------------------------

------------------------
--       CHATS
------------------------

CREATE TABLE chats (
    id           BIGSERIAL PRIMARY KEY,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Users to chats join table, no access rights
CREATE TABLE users_chats (
    id        BIGSERIAL PRIMARY KEY,
    chat_id   BIGINT REFERENCES chats(id),
    user_id   BIGINT REFERENCES users(id),
    UNIQUE    (user_id, chat_id)
);

CREATE TABLE chat_messages (
    id            BIGSERIAL PRIMARY KEY,
    author_id     BIGINT REFERENCES users(id),
    text_content  TEXT,
    img_content   TEXT,
    sent_at       TIMESTAMP WITH TIME ZONE NOT NULL,
    edited_at     TIMESTAMP WITH TIME ZONE NOT NULL
);

------------------------
--       CHATS
------------------------
