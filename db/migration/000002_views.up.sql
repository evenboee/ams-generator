CREATE OR REPLACE VIEW answers_summary AS
SELECT a.*, COUNT(*) AS count, AVG(f.rating) AS rating
FROM "answers" a
LEFT JOIN "feedback" f ON a.id = f.answer
GROUP BY a.id;

CREATE OR REPLACE VIEW question_summary AS
SELECT q.*, COUNT(*) AS count, AVG(a.rating) AS rating
FROM "questions" q
LEFT JOIN "answers_summary" a ON q.id = a.question
GROUP BY q.id;

CREATE OR REPLACE VIEW assignment_summary AS
SELECT a.*, AVG(f.rating) AS rating
    FROM "assignments" a
    LEFT JOIN "submissions" s
        ON a.id = s.assignment
    LEFT JOIN "reviews" r
        ON s.id = r.submission
    LEFT JOIN "feedback" f
        ON r.id = f.review
    GROUP BY a.id;

CREATE OR REPLACE VIEW course_summary AS
SELECT c.*, COUNT(*) AS count, AVG(a.rating) AS avg_rating
FROM "courses" c
LEFT JOIN "assignment_summary" a ON c.id = a.course
GROUP BY c.id;
