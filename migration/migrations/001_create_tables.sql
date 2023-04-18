-- +goose Up
create table if not exists users (
    user_id uuid not null primary key,
    user_mame text not null
);
comment on column users.user_mame is 'user name';


create table if not exists word(
    id uuid default gen_random_uuid() not null primary key,
    text text not null,
    lang text not null
);
create unique index if not exists idx_unique_word__text_lang
    on word (text, lang);


create table if not exists language (
    code text not null primary key,
    lang text not null
);
create unique index if not exists idx_unique_languages__lang_code
    on language (lang, code);


create table if not exists example (
    id uuid default gen_random_uuid() not null primary key,
    original text,
    translate text
);


create table if not exists dictionary (
    user_id uuid references users (user_id) not null,
    original_word uuid references word (id) not null,
    original_lang text references language (code) not null,
    translate_lang text references language (code) not null,
    translate_word uuid[] not null,
    pronunciation text,
    example uuid[]
);
create unique index if not exists idx_unique_dictionary__user_id_original_word
    on dictionary (user_id, original_word);



-- +goose Down
drop table if exists users;
drop table if exists language;
drop table if exists example;
drop table if exists word;
drop table if exists dictionary;