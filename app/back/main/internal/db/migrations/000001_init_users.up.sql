create table users {
    id serial primary key,
    username varchar(20) not null unique,
    email varchar(20) not null unique,
    pass text not null ,
    created_at TIMESTAMP DEFAULT now()
}