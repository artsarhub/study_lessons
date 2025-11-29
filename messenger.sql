-- Есть 3 сущности - пользователь, чат, сообщение

-- У пользователя есть имя и дата регистрации
-- У чата есть название и дата создания
-- У сообщения есть текст, автор и дата создания
-- Пользователь может состоять в нескольких чатах одновременно
-- Сообщение обязательно принадлежит чату, сообщение не может принадлежать более чем 1 чату одновременно

-- Нужно описать предметную область в виде таблиц

CREATE DATABASE messenger;

CREATE TABLE users
(
    id         INT PRIMARY KEY,
    name       TEXT      NOT NULL,
    created_at timestamp NOT NULL
);

CREATE TABLE chats
(
    id         INT PRIMARY KEY,
    name       TEXT      NOT NULL,
    created_at timestamp NOT NULL
);

CREATE TABLE messages
(
    id         INT PRIMARY KEY,
    content    TEXT      NOT NULL,
    author_id  INT       NOT NULL,
    created_at timestamp NOT NULL,
    chat_id    INT       NOT NULL
);

CREATE TABLE users_chats
(
    user_id INT,
    chat_id INT,
    PRIMARY KEY (user_id, chat_id)
);

------------------------------------------------------------------------------------------------------------------------

select max(id) from users;

insert into users (id, name, created_at) VALUES (101, 'Вася', '2023-04-01 12:00:00');

explain analyse
SELECT uc.chat_id, c.name AS chat_name
FROM Users u
         LEFT JOIN Users_Chats uc ON uc.user_id = u.id
         LEFT JOIN Chats c ON c.id = uc.chat_id
WHERE u.name = 'Вася';

create index users_name_idx on users(name);

explain analyse
select * from messages where author_id = 85;

create index messages_author_id_idx on messages(author_id);

create table messages_202511_202512 partition of messages for values from () to ();
