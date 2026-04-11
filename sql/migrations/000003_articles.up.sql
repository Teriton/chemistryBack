CREATE TABLE lessons(
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	title VARCHAR(256) UNIQUE NOT NULL
);


CREATE TABLE lessons_completed (
	user_id INT NOT NULL,
	lesson_id INT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users (id),
	FOREIGN KEY (lesson_id) REFERENCES lessons (id) ON DELETE CASCADE,
	PRIMARY KEY (user_id, lesson_id)
);
