create table if not exists users (
    user_id bigint PRIMARY KEY,
    user_mame text not null
);

create table if not exists languages (
    lang text PRIMARY KEY,
    code text not null
);

create table if not exists examples (
    id bigserial PRIMARY KEY,
    original text,
    translate text
);

create table if not exists words (
    id bigserial PRIMARY KEY,
    original_lang text not null,
    translate_lang text not null,
    original text not null,
    translate text not null,
    pronunciation text,
    example examples
);


create table if not exists dictionary (
    id bigserial PRIMARY KEY,
    user_id bigint not null,
    word_id bigint not null,
    created timestamptz
);