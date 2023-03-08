-- +goose Up
create table if not exists users (
    user_id uuid not null primary key,
    user_mame text not null
);
comment on column users.user_mame is 'user name';


create table if not exists languages (
    code text not null primary key,
    lang text not null
);
create unique index if not exists idx_unique_languages__lang_code
    on languages (lang, code);


create table if not exists examples (
    id uuid default gen_random_uuid() not null primary key,
    original text,
    translate text
);


create table if not exists word(
    id uuid default gen_random_uuid() not null primary key,
    text text not null,
    pronunciation text,
    lang text not null
);
create unique index if not exists idx_unique_word__text_lang
    on word (text, lang);


create table if not exists dictionary (
    user_id uuid not null references users (user_id) not null,
    original_word uuid references word (id) not null,
    original_lang text references languages (code) not null,
    translate_lang text references languages (code) not null,
    translate_word uuid[] not null,
    example uuid[]
);
create unique index if not exists idx_unique_dictionary__user_id_original_word
    on dictionary (user_id, original_word);



-- +goose Down
drop table if exists users;
drop table if exists languages;
drop table if exists examples;
drop table if exists words;
drop table if exists dictionary;