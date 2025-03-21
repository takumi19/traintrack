create table programs_permissions (
    permission_id serial primary key,
    user_id       integer not null,
    program_id    integer not null,
    can_view      boolean default false not null,
    can_modify    boolean default false not null,

    constraint fk_user
        foreign key(user_id)
        references users(id) on delete cascade,  -- if user is deleted, their permissions are too.
    constraint fk_program
        foreign key(program_id)
        references program_templates(id) on delete cascade,
    unique (user_id, program_id)
);
