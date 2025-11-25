-- Есть 3 сущности - пользователь, чат, сообщение

-- У пользователя есть имя и дата регистрации
-- У чата есть название и дата создания
-- У сообщения есть текст, автор и дата создания
-- Пользователь может состоять в нескольких чатах одновременно
-- Сообщение обязательно принадлежит чату, сообщение не может принадлежать более чем 1 чату одновременно

-- Нужно описать предметную область в виде таблиц

CREATE TABLE Users
(
    id         INT PRIMARY KEY,
    name       TEXT      NOT NULL,
    created_at timestamp NOT NULL
);

CREATE TABLE Chats
(
    id         INT PRIMARY KEY,
    name       TEXT      NOT NULL,
    created_at timestamp NOT NULL
);

CREATE TABLE Messages
(
    id         INT PRIMARY KEY,
    content    TEXT      NOT NULL,
    author_id  INT       NOT NULL,
    created_at timestamp NOT NULL,
    chat_id    INT       NOT NULL
);

CREATE TABLE Users_Chats
(
    user_id INT,
    chat_id INT,
    PRIMARY KEY (user_id, chat_id)
);

------------------------------------------------------------------------------------------------------------------------

SELECT uc.chat_id, c.name AS chat_name
FROM Users u
         LEFT JOIN Users_Chats uc ON uc.user_id = u.id
         LEFT JOIN Chats c ON c.id = uc.chat_id
WHERE u.name = 'Вася';
