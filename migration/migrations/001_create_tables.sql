-- +goose Up
create table if not exists users (
    id uuid default gen_random_uuid() primary key,
    name text not null,
    password_hash text not null,
    last_visit date
);

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
    word_id uuid,
    example text
);


create table if not exists dictionary (
    user_id uuid references users (id) not null,
    original_word uuid references word (id) not null,
    pronunciation text,
    translate_word uuid[] not null,
    examples uuid[]
);
create unique index if not exists idx_unique_dictionary__user_id_original_word
    on dictionary (user_id, original_word);



-- +goose Down
drop table if exists users;
drop table if exists language;
drop table if exists example;
drop table if exists word;
drop table if exists dictionary;