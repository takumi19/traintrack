INSERT INTO users (full_name, login, email, password_hash)
VALUES
    ('Alice Smith', 'alices', 'alice.smith@example.com', 'hashed_password_1'),
    ('Bob Johnson', 'bobj', 'bob.johnson@example.com', 'hashed_password_2'),
    ('Charlie Brown', 'charlieb', 'charlie.brown@example.com', 'hashed_password_3'),
    ('Diana Miller', 'dianam', 'diana.miller@example.com', 'hashed_password_4'),
    ('Eve Wilson', 'evew', 'eve.wilson@example.com', 'hashed_password_5'),
    ('Frank Garcia', 'frankg', 'frank.garcia@example.com', 'hashed_password_6'),
    ('Grace Rodriguez', 'gracer', 'grace.rodriguez@example.com', 'hashed_password_7'),
    ('Henry Martinez', 'henrym', 'henry.martinez@example.com', 'hashed_password_8'),
    ('Ivy Anderson', 'ivya', 'ivy.anderson@example.com', 'hashed_password_9'),
    ('Jack Thomas', 'jackt', 'jack.thomas@example.com', 'hashed_password_10');

INSERT INTO exercises (id, name, notes, is_rep_based, is_bodyweight) VALUES
(1, 'Back Squat', 'Barbell squat targeting quads, glutes, and core.', TRUE, FALSE),
(2, 'Bench Press', 'Barbell press for chest, shoulders, and triceps.', TRUE, FALSE),
(3, 'Incline Dumbbell Press', 'Dumbbell press on an incline for upper chest.', TRUE, FALSE),
(4, 'Tricep Pushdown', 'Cable exercise isolating the triceps.', TRUE, FALSE),
(5, 'Leg Press', 'Machine-based exercise for quads and glutes.', TRUE, FALSE);

INSERT INTO exercises_muscle_groups (exercise_id, worked_muscle_group) VALUES
-- Back Squat (ID 1)
(1, 'legs'),      -- Primary: Quads, hamstrings, glutes
(1, 'core'),      -- Secondary: Stabilizing muscles (abs, lower back)
-- Bench Press (ID 2)
(2, 'chest'),     -- Primary: Pectoralis major
(2, 'triceps'),   -- Secondary: Triceps
(2, 'shoulders'), -- Secondary: Anterior deltoids
-- Incline Dumbbell Press (ID 3)
(3, 'chest'),     -- Primary: Upper pectoralis major
(3, 'triceps'),   -- Secondary: Triceps
(3, 'shoulders'), -- Secondary: Anterior deltoids
-- Tricep Pushdown (ID 4)
(4, 'triceps'),   -- Primary: Triceps (only major muscle group worked)
-- Leg Press (ID 5)
(5, 'legs'),      -- Primary: Quads, hamstrings, glutes
(5, 'core');      -- Secondary: Stabilizing muscles (minor involvement)

INSERT INTO program_templates (author_id, name, notes, created_at, updated_at) VALUES
(1, 'Strength Program', 'A 4-week program focused on building raw strength.', '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
(1, 'Hypertrophy Program', 'A 6-week program for muscle growth.', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00');

INSERT INTO program_weeks (program_template_id, week_number, notes, created_at, updated_at) VALUES
-- Strength Program (ID 1)
(1, 1, 'Focus on heavy lifts.', '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
(1, 2, 'Increase weight slightly.', '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
(1, 3, 'Peak week.', '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
(1, 4, 'Deload week.', '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
-- Hypertrophy Program (ID 2)
(2, 1, 'High volume start.', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
(2, 2, NULL, '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
(2, 3, 'Introduce supersets.', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
(2, 4, NULL, '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
(2, 5, 'Push to failure.', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
(2, 6, 'Recovery week.', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00');

INSERT INTO program_workouts (program_week_id, workout_index, title, notes, created_at, updated_at) VALUES
-- Strength Program, Week 1 (ID 1)
(1, 1, 'Squat Day', 'Focus on form and heavy weight.', '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
(1, 3, 'Bench Day', NULL, '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
(1, 5, 'Deadlift Day', 'Warm up thoroughly.', '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
-- Hypertrophy Program, Week 1 (ID 5)
(5, 1, 'Chest & Triceps', 'High volume, moderate weight.', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
(5, 2, 'Back & Biceps', NULL, '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
(5, 4, 'Legs', 'Focus on quads and glutes.', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
(5, 5, 'Shoulders & Core', 'Light weights, high reps.', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00');

INSERT INTO program_workout_exercises (program_workout_id, exercise_id, order_index, notes, created_at, updated_at) VALUES
-- Squat Day (ID 1)
(1, 1, 1, 'Keep core tight.', '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
(1, 5, 2, 'Controlled tempo.', '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
-- Chest & Triceps (ID 4)
(4, 2, 1, 'Use a spotter.', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
(4, 3, 2, NULL, '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
(4, 4, 3, 'Full extension.', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00');

INSERT INTO program_workout_sets (
    program_workout_exercise_id, set_number, rpe,
    suggested_reps_min, suggested_reps_max, suggested_reps,
    suggested_weight_min, suggested_weight_max, suggested_weight,
    suggested_time_min, suggested_time_max, suggested_time,
    notes, created_at, updated_at
) VALUES
-- Back Squat (ID 1)
(1, 1, 7, NULL, NULL, 5, NULL, NULL, 100.00, NULL, NULL, NULL, 'Warm-up set', '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
(1, 2, 8, NULL, NULL, 5, NULL, NULL, 120.00, NULL, NULL, NULL, 'Working set', '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
(1, 3, 8, NULL, NULL, 5, NULL, NULL, 120.00, NULL, NULL, NULL, 'Working set', '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
-- Bench Press (ID 3)
(3, 1, 6, 8, 10, NULL, 50.00, 60.00, NULL, NULL, NULL, NULL, 'Moderate effort', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
(3, 2, 7, 8, 10, NULL, 50.00, 60.00, NULL, NULL, NULL, NULL, NULL, '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
(3, 3, 7, 8, 10, NULL, 50.00, 60.00, NULL, NULL, NULL, NULL, 'Push to fatigue', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
-- Leg Press (ID 2, Strength Program, Squat Day)
(2, 1, 6, NULL, NULL, 10, NULL, NULL, 150.00, NULL, NULL, NULL, 'Warm-up with moderate weight', '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
(2, 2, 7, NULL, NULL, 8, NULL, NULL, 180.00, NULL, NULL, NULL, 'Focus on slow negatives', '2025-03-01 10:00:00+00', '2025-03-01 10:00:00+00'),
-- Incline Dumbbell Press (ID 4, Hypertrophy Program, Chest & Triceps)
(4, 1, 6, 10, 12, NULL, 20.00, 25.00, NULL, NULL, NULL, NULL, 'Controlled reps', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
(4, 2, 7, 10, 12, NULL, 20.00, 25.00, NULL, NULL, NULL, NULL, 'Squeeze at the top', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
-- Tricep Pushdown (ID 5, Hypertrophy Program, Chest & Triceps)
(5, 1, 5, 12, 15, NULL, NULL, NULL, 15.00, NULL, NULL, NULL, 'Light set to warm up', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00'),
(5, 2, 6, 12, 15, NULL, NULL, NULL, 20.00, NULL, NULL, NULL, 'Full range of motion', '2025-03-02 14:00:00+00', '2025-03-02 14:00:00+00');
