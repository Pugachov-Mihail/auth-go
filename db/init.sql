create table users_my (
    id serial primary key unique,
    email varchar unique,
    pass_hash varchar,
    steam_id bigint,
    login varchar
);

drop table access_token;

create table access_token (
    user_id bigint unique,
    token varchar
);

insert into access_token (user_id, token) VALUES (9, 'commodo anim');