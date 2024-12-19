-- +goose Up
CREATE TABLE "user" (
    "id" UUID NOT NULL UNIQUE,
    "Name" TEXT,
    "team_id" UUID,
    "email" TEXT,
    "is_vitian" BOOLEAN,
    "reg_no" TEXT,
    "password" TEXT,
    "phone_no" TEXT,
    "role" TEXT,
    "is_leader" BOOLEAN,
    "college" TEXT,
    "is_verified" BOOLEAN,
    PRIMARY KEY ("id")
);

CREATE TABLE "team" (
    "id" UUID NOT NULL UNIQUE,
    "name" TEXT,
    "number_of_people" INTEGER,
    "users" UUID,
    "submission" UUID,
    "round_qualified" INTEGER,
    "code" TEXT,
    PRIMARY KEY ("id")
);

CREATE TABLE "score" (
    "id" UUID NOT NULL UNIQUE,
    "panelist" UUID,
    "team_id" UUID,
    "design" INTEGER,
    "implementation" INTEGER,
    "presentation" INTEGER,
    "round" INTEGER,
    PRIMARY KEY ("id")
);

CREATE TABLE "submission" (
    "id" UUID NOT NULL UNIQUE,
    "github_link" TEXT,
    "figma_link" TEXT,
    "ppt_link" TEXT,
    "other_link" TEXT,
    PRIMARY KEY ("id")
);

ALTER TABLE "team" ADD FOREIGN KEY ("users") REFERENCES "user" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "team" ADD FOREIGN KEY ("id") REFERENCES "score" ("team_id") ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "user" ADD FOREIGN KEY ("id") REFERENCES "score" ("panelist") ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "team" ADD FOREIGN KEY ("submission") REFERENCES "submission" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "user" ADD FOREIGN KEY ("team_id") REFERENCES "team" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;

-- +goose Down
DROP table "user";

DROP TABLE "team";

DROP TABLE "score";

DROP TABLE "submission";
