\connect expensedb;
CREATE TABLE IF NOT EXISTS expenses (
  id SERIAL PRIMARY KEY,
  title TEXT,
  amount FLOAT,
  note TEXT,
  tags TEXT[]
);

INSERT INTO expenses ("title", "amount", "note", "tags") values ('strawberry smoothie', 79.0, 'night market promotion discount 10 bath', ARRAY ['food', 'beverage']);
