CREATE TABLE IF NOT EXISTS "contents" (
	"id" SERIAL PRIMARY KEY,
	"title" VARCHAR(255) NOT NULL,
	"excerpt" VARCHAR(255) NOT NULL,
    "description" text NOT NULL,
    "image" text NULL,
    "status" VARCHAR(20) NOT NULL DEFAULT 'PUBLISH',
	"tags" text NOT NULL,
	"created_by_id" INT REFERENCES users(id) ON DELETE CASCADE,
	"category_id" INT REFERENCES categories(id) ON DELETE CASCADE,
	"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"updated_at" TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_contents_created_by_id ON contents(created_by_id);
CREATE INDEX idx_contents_category_id ON contents(category_id);