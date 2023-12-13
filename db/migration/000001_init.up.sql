-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-12-13T10:36:06.496Z

CREATE TABLE "courses" (
  "id" serial PRIMARY KEY,
  "code" varchar NOT NULL,
  "year" int NOT NULL
);

CREATE TABLE "assignments" (
  "id" serial PRIMARY KEY,
  "name" varchar NOT NULL,
  "course" int NOT NULL,
  "reviews_per_submission" int NOT NULL DEFAULT 0,
  "time_due" timestamptz NOT NULL
);

CREATE TABLE "questions" (
  "id" serial PRIMARY KEY,
  "assignment" int NOT NULL,
  "prompt" varchar NOT NULL,
  "order" int NOT NULL
);

CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "display_name" varchar NOT NULL
);

CREATE TABLE "user_course_enrollments" (
  "user" varchar NOT NULL,
  "course" int NOT NULL,
  PRIMARY KEY ("user", "course")
);

CREATE TABLE "submissions" (
  "id" serial PRIMARY KEY,
  "user" varchar NOT NULL,
  "assignment" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "answers" (
  "id" serial PRIMARY KEY,
  "question" int NOT NULL,
  "submission" int NOT NULL,
  "answer" varchar NOT NULL
);

CREATE TABLE "reviews" (
  "id" serial PRIMARY KEY,
  "submission" int NOT NULL,
  "reviewer_id" varchar NOT NULL,
  "finished_at" timestamptz,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "feedback" (
  "id" serial PRIMARY KEY,
  "review" int NOT NULL,
  "answer" int NOT NULL,
  "rating" int NOT NULL,
  "feedback" varchar NOT NULL
);

CREATE UNIQUE INDEX ON "courses" ("code", "year");

CREATE INDEX ON "assignments" ("course");

CREATE INDEX ON "questions" ("order");

CREATE INDEX ON "questions" ("assignment");

CREATE INDEX ON "user_course_enrollments" ("user");

CREATE INDEX ON "user_course_enrollments" ("course");

CREATE INDEX ON "submissions" ("user");

CREATE INDEX ON "submissions" ("assignment");

CREATE INDEX ON "submissions" ("created_at");

CREATE UNIQUE INDEX ON "answers" ("question", "submission");

CREATE INDEX ON "answers" ("submission");

CREATE INDEX ON "answers" ("question");

CREATE UNIQUE INDEX ON "reviews" ("submission", "reviewer_id");

CREATE INDEX ON "reviews" ("submission");

CREATE INDEX ON "reviews" ("reviewer_id");

CREATE UNIQUE INDEX ON "feedback" ("review", "answer");

CREATE INDEX ON "feedback" ("review");

CREATE INDEX ON "feedback" ("answer");

ALTER TABLE "assignments" ADD FOREIGN KEY ("course") REFERENCES "courses" ("id");

ALTER TABLE "questions" ADD FOREIGN KEY ("assignment") REFERENCES "assignments" ("id");

ALTER TABLE "user_course_enrollments" ADD FOREIGN KEY ("user") REFERENCES "users" ("id");

ALTER TABLE "user_course_enrollments" ADD FOREIGN KEY ("course") REFERENCES "courses" ("id");

ALTER TABLE "submissions" ADD FOREIGN KEY ("user") REFERENCES "users" ("id");

ALTER TABLE "submissions" ADD FOREIGN KEY ("assignment") REFERENCES "assignments" ("id");

ALTER TABLE "answers" ADD FOREIGN KEY ("question") REFERENCES "questions" ("id");

ALTER TABLE "answers" ADD FOREIGN KEY ("submission") REFERENCES "submissions" ("id");

ALTER TABLE "reviews" ADD FOREIGN KEY ("submission") REFERENCES "submissions" ("id");

ALTER TABLE "reviews" ADD FOREIGN KEY ("reviewer_id") REFERENCES "users" ("id");

ALTER TABLE "feedback" ADD FOREIGN KEY ("review") REFERENCES "reviews" ("id");

ALTER TABLE "feedback" ADD FOREIGN KEY ("answer") REFERENCES "answers" ("id");
