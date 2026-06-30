-- organization
CREATE TABLE IF NOT EXISTS "organization"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "name" varchar(255) NOT NULL,
  "kind" varchar(20) NOT NULL DEFAULT 'team',
  "description" text,
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  UNIQUE ("name")
);

-- user
CREATE TABLE IF NOT EXISTS "user"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "google_sub" varchar(255) NOT NULL,
  "name" varchar(255) NOT NULL,
  "email" varchar(255) NOT NULL,
  "avatar_url" text,
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  UNIQUE ("google_sub"),
  UNIQUE ("email")
);

-- organization_membership
CREATE TABLE IF NOT EXISTS "organization_membership"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "user_id" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "role" smallint NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  UNIQUE ("user_id", "organization_id"),
  FOREIGN KEY ("user_id") REFERENCES "user" ("id_natural") ON DELETE CASCADE,
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON DELETE CASCADE
);

-- site
CREATE TABLE IF NOT EXISTS "site"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "name" varchar(255) NOT NULL,
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural")
);

-- location
CREATE TABLE IF NOT EXISTS "location"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "site_id" varchar(36) NOT NULL,
  "name" varchar(255) NOT NULL,
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural"),
  FOREIGN KEY ("site_id") REFERENCES "site" ("id_natural")
);

-- task
CREATE TABLE IF NOT EXISTS "task"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "name" varchar(255) NOT NULL,
  "description" text,
  "manual_url" text,
  "priority" smallint NOT NULL,
  "difficulty" smallint NOT NULL,
  "status" smallint NOT NULL DEFAULT 0,
  "deadline" timestamptz NOT NULL,
  "robot_type" varchar(255),
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural")
);

-- robot
CREATE TABLE IF NOT EXISTS "robot"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "location_id" varchar(36) NOT NULL,
  "name" varchar(255) NOT NULL,
  "robot_type" varchar(255),
  "status" smallint NOT NULL DEFAULT 0,
  "leader_status" smallint,
  "leader_fault_started_at" timestamptz,
  "fault_started_at" timestamptz,
  "last_heartbeat_at" timestamptz,
  "offline_reason" varchar(255),
  "robot_config" jsonb,
  "active_episode_id" varchar(36),
  "active_user_id" varchar(36),
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural")
);

-- task_version
CREATE TABLE IF NOT EXISTS "task_version"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "task_id" varchar(36) NOT NULL,
  "version" varchar(50) NOT NULL,
  "schema_hash" varchar(255),
  "is_active" boolean NOT NULL DEFAULT true,
  "approval_status" integer NOT NULL DEFAULT 0,
  "target_duration_seconds" integer,
  "target_episode_count" integer,
  "target_duration_per_episode_seconds" integer,
  "display_name" varchar(100),
  "parameters" jsonb,
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural")
);

-- subtask
CREATE TABLE IF NOT EXISTS "subtask"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "task_version_id" varchar(36) NOT NULL,
  "order_index" integer NOT NULL,
  "name" varchar(255) NOT NULL,
  "description" text,
  "target_duration_seconds" integer,
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural")
);

-- episode
CREATE TABLE IF NOT EXISTS "episode"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "task_version_id" varchar(36) NOT NULL,
  "location_id" varchar(36) NOT NULL,
  "robot_id" varchar(36) NOT NULL,
  "user_id" varchar(36) NOT NULL,
  "recorded_by" varchar(36),
  "started_at" timestamptz,
  "finished_at" timestamptz,
  "collection_status" smallint NOT NULL DEFAULT 0,
  "error_details" jsonb,
  "parameter_values" jsonb,
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural")
);

-- episode_grade
CREATE TABLE IF NOT EXISTS "episode_grade"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "episode_id" varchar(36) NOT NULL,
  "user_id" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "grade" double precision NOT NULL,
  "comment" text,
  "graded_at" timestamptz NOT NULL,
  PRIMARY KEY ("episode_id",
  "user_id"),
  FOREIGN KEY ("episode_id") REFERENCES "episode" ("id_natural") ON DELETE CASCADE,
  FOREIGN KEY ("user_id") REFERENCES "user" ("id_natural"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural")
);

-- episode_sub_task
CREATE TABLE IF NOT EXISTS "episode_sub_task"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "episode_id" varchar(36) NOT NULL,
  "sub_task_id" varchar(36) NOT NULL,
  "collection_status" smallint NOT NULL DEFAULT 0,
  "task_result" smallint NOT NULL DEFAULT 0,
  CONSTRAINT "episode_sub_task_episode_sub_task_unique" UNIQUE ("episode_id",
  "sub_task_id"),
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural"),
  FOREIGN KEY ("episode_id") REFERENCES "episode" ("id_natural"),
  FOREIGN KEY ("sub_task_id") REFERENCES "subtask" ("id_natural")
);

-- episode_sub_task_execution
CREATE TABLE IF NOT EXISTS "episode_sub_task_execution"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "episode_sub_task_id" varchar(36) NOT NULL,
  "execution_status" smallint NOT NULL DEFAULT 0,
  "started_at" timestamptz,
  "finished_at" timestamptz,
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural"),
  FOREIGN KEY ("episode_sub_task_id") REFERENCES "episode_sub_task" ("id_natural")
);

-- episode_stats_hourly
CREATE TABLE IF NOT EXISTS "episode_stats_hourly"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "location_id" varchar(36) NOT NULL,
  "robot_id" varchar(36) NOT NULL,
  "period_start" timestamptz NOT NULL,
  "total_duration_seconds" bigint NOT NULL DEFAULT 0,
  "episode_count" int NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural"),
  FOREIGN KEY ("location_id") REFERENCES "location" ("id_natural"),
  FOREIGN KEY ("robot_id") REFERENCES "robot" ("id_natural")
);

-- episode_stats_daily
CREATE TABLE IF NOT EXISTS "episode_stats_daily"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "location_id" varchar(36) NOT NULL,
  "robot_id" varchar(36) NOT NULL,
  "period_start" timestamptz NOT NULL,
  "total_duration_seconds" bigint NOT NULL DEFAULT 0,
  "episode_count" int NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural"),
  FOREIGN KEY ("location_id") REFERENCES "location" ("id_natural"),
  FOREIGN KEY ("robot_id") REFERENCES "robot" ("id_natural")
);

-- episode_stats_monthly
CREATE TABLE IF NOT EXISTS "episode_stats_monthly"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "location_id" varchar(36) NOT NULL,
  "robot_id" varchar(36) NOT NULL,
  "period_start" timestamptz NOT NULL,
  "total_duration_seconds" bigint NOT NULL DEFAULT 0,
  "episode_count" int NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural"),
  FOREIGN KEY ("location_id") REFERENCES "location" ("id_natural"),
  FOREIGN KEY ("robot_id") REFERENCES "robot" ("id_natural")
);

-- task_version_stats
CREATE TABLE IF NOT EXISTS "task_version_stats"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "task_version_id" varchar(36) NOT NULL,
  "total_duration_seconds" bigint NOT NULL DEFAULT 0,
  "episode_count" int NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  UNIQUE ("task_version_id"),
  FOREIGN KEY ("task_version_id") REFERENCES "task_version" ("id_natural") ON DELETE CASCADE
);

-- task_category_type
CREATE TABLE IF NOT EXISTS "task_category_type"(
  "id" varchar(36) NOT NULL,
  "slug" varchar(100) NOT NULL,
  "name" varchar(255) NOT NULL,
  PRIMARY KEY ("id")
);

-- task_tag
CREATE TABLE IF NOT EXISTS "task_tag"(
  "id" varchar(36) NOT NULL,
  "name" varchar(255) NOT NULL,
  "category_type_id" varchar(36) NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("category_type_id") REFERENCES "task_category_type" ("id")
);

-- task_tag_assignment
CREATE TABLE IF NOT EXISTS "task_tag_assignment"(
  "task_id" varchar(36) NOT NULL,
  "tag_id" varchar(36) NOT NULL,
  FOREIGN KEY ("task_id") REFERENCES "task" ("id_natural") ON DELETE CASCADE,
  FOREIGN KEY ("tag_id") REFERENCES "task_tag" ("id") ON DELETE CASCADE
);

-- user_location_assignment
CREATE TABLE IF NOT EXISTS "user_location_assignment"(
  "user_id" varchar(36) NOT NULL,
  "location_id" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  PRIMARY KEY ("user_id",
  "location_id"),
  FOREIGN KEY ("user_id") REFERENCES "user" ("id_natural") ON DELETE CASCADE,
  FOREIGN KEY ("location_id") REFERENCES "location" ("id_natural") ON DELETE CASCADE,
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON DELETE CASCADE
);

-- user_site_assignment
CREATE TABLE IF NOT EXISTS "user_site_assignment"(
  "user_id" varchar(36) NOT NULL,
  "site_id" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  PRIMARY KEY ("user_id",
  "site_id"),
  FOREIGN KEY ("user_id") REFERENCES "user" ("id_natural") ON DELETE CASCADE,
  FOREIGN KEY ("site_id") REFERENCES "site" ("id_natural") ON DELETE CASCADE,
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON DELETE CASCADE
);

-- robot_uptime_hourly
CREATE TABLE IF NOT EXISTS "robot_uptime_hourly"(
  "robot_id" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "location_id" varchar(36) NOT NULL,
  "period_start" timestamptz NOT NULL,
  "uptime_seconds" bigint NOT NULL DEFAULT 0,
  PRIMARY KEY ("robot_id",
  "period_start"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  FOREIGN KEY ("location_id") REFERENCES "location" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  FOREIGN KEY ("robot_id") REFERENCES "robot" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- api_key
CREATE TABLE IF NOT EXISTS "api_key"(
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" BIGSERIAL NOT NULL,
  "id_natural" varchar(36) NOT NULL,
  "organization_id" varchar(36) NOT NULL,
  "user_id" varchar(36) NOT NULL,
  "robot_id" varchar(36),
  "name" varchar(255) NOT NULL,
  "key_hash" char(64) NOT NULL,
  "key_hint" varchar(16) NOT NULL,
  "expires_at" TIMESTAMPTZ,
  "last_used_at" TIMESTAMPTZ,
  "revoked_at" TIMESTAMPTZ,
  PRIMARY KEY ("id"),
  UNIQUE ("id_natural"),
  UNIQUE ("key_hash"),
  FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural"),
  FOREIGN KEY ("user_id") REFERENCES "user" ("id_natural") ON DELETE CASCADE,
  FOREIGN KEY ("robot_id") REFERENCES "robot" ("id_natural") ON DELETE CASCADE
);

-- indexes
CREATE INDEX "episode_organization_id_idx" ON "episode" ("organization_id");
CREATE INDEX "episode_location_id_idx" ON "episode" ("location_id");
CREATE INDEX "episode_robot_id_idx" ON "episode" ("robot_id");
CREATE INDEX "episode_user_id_idx" ON "episode" ("user_id");
CREATE INDEX "episode_started_at_idx" ON "episode" ("started_at");
CREATE INDEX "episode_task_version_stats_idx" ON "episode" ("task_version_id", "collection_status") INCLUDE ("started_at", "finished_at") WHERE (collection_status = 3);
CREATE INDEX "episode_org_created_at_idx" ON "episode" (organization_id, created_at DESC);
CREATE UNIQUE INDEX "episode_one_recording_per_robot" ON "episode" ("robot_id") WHERE (collection_status = 1);
CREATE INDEX "episode_grade_user_id_idx" ON "episode_grade" ("user_id");
CREATE INDEX "episode_grade_organization_id_idx" ON "episode_grade" ("organization_id");
CREATE INDEX "organization_membership_user_id_idx" ON "organization_membership" ("user_id");
CREATE INDEX "organization_membership_organization_id_idx" ON "organization_membership" ("organization_id");
CREATE INDEX "task_organization_id_idx" ON "task" ("organization_id");
CREATE INDEX "task_version_organization_id_idx" ON "task_version" ("organization_id");
CREATE INDEX "subtask_organization_id_idx" ON "subtask" ("organization_id");
CREATE INDEX "robot_organization_id_idx" ON "robot" ("organization_id");
CREATE INDEX "location_organization_id_idx" ON "location" ("organization_id");
CREATE INDEX "site_organization_id_idx" ON "site" ("organization_id");
CREATE UNIQUE INDEX "episode_stats_hourly_org_loc_robot_period_idx" ON "episode_stats_hourly" ("organization_id", "location_id", "robot_id", "period_start");
CREATE INDEX "episode_stats_hourly_period_start_idx" ON "episode_stats_hourly" ("period_start");
CREATE UNIQUE INDEX "episode_stats_daily_org_loc_robot_period_idx" ON "episode_stats_daily" ("organization_id", "location_id", "robot_id", "period_start");
CREATE INDEX "episode_stats_daily_period_start_idx" ON "episode_stats_daily" ("period_start");
CREATE UNIQUE INDEX "episode_stats_monthly_org_loc_robot_period_idx" ON "episode_stats_monthly" ("organization_id", "location_id", "robot_id", "period_start");
CREATE INDEX "episode_stats_monthly_period_start_idx" ON "episode_stats_monthly" ("period_start");
CREATE INDEX "idx_user_location_assignment_location_id" ON "user_location_assignment" ("location_id");
CREATE INDEX "idx_user_site_assignment_site_id" ON "user_site_assignment" ("site_id");
CREATE INDEX "robot_uptime_hourly_period_start_idx" ON "robot_uptime_hourly" ("period_start");
CREATE INDEX "episode_sub_task_execution_episode_sub_task_id_created_at_idx" ON "episode_sub_task_execution" ("episode_sub_task_id", "created_at");
CREATE INDEX "api_key_organization_id_idx" ON "api_key" ("organization_id");
CREATE INDEX "api_key_robot_id_idx" ON "api_key" ("robot_id");
CREATE INDEX "api_key_user_id_idx" ON "api_key" ("user_id");
