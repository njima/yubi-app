-- Create "organization" table
CREATE TABLE "organization" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "organization_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "organization_name_key" UNIQUE ("name")
);
-- Create "robot" table
CREATE TABLE "robot" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "location_id" character varying(36) NOT NULL,
  "name" character varying(255) NOT NULL,
  "robot_type" character varying(255) NULL,
  "status" smallint NOT NULL DEFAULT 0,
  "leader_status" smallint NULL,
  "leader_fault_started_at" timestamptz NULL,
  "fault_started_at" timestamptz NULL,
  "last_heartbeat_at" timestamptz NULL,
  "offline_reason" character varying(255) NULL,
  "robot_config" jsonb NULL,
  "active_episode_id" character varying(36) NULL,
  "active_user_id" character varying(36) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "robot_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "robot_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "robot_organization_id_idx" to table: "robot"
CREATE INDEX "robot_organization_id_idx" ON "robot" ("organization_id");
-- Create "user" table
CREATE TABLE "user" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "name" character varying(255) NOT NULL,
  "email" character varying(255) NOT NULL,
  "role" smallint NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "user_email_key" UNIQUE ("email"),
  CONSTRAINT "user_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "user_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "user_organization_id_idx" to table: "user"
CREATE INDEX "user_organization_id_idx" ON "user" ("organization_id");
-- Create "api_key" table
CREATE TABLE "api_key" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "user_id" character varying(36) NOT NULL,
  "robot_id" character varying(36) NULL,
  "name" character varying(255) NOT NULL,
  "key_hash" character(64) NOT NULL,
  "key_hint" character varying(16) NOT NULL,
  "expires_at" timestamptz NULL,
  "last_used_at" timestamptz NULL,
  "revoked_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "api_key_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "api_key_key_hash_key" UNIQUE ("key_hash"),
  CONSTRAINT "api_key_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "api_key_robot_id_fkey" FOREIGN KEY ("robot_id") REFERENCES "robot" ("id_natural") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "api_key_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "user" ("id_natural") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "api_key_organization_id_idx" to table: "api_key"
CREATE INDEX "api_key_organization_id_idx" ON "api_key" ("organization_id");
-- Create index "api_key_robot_id_idx" to table: "api_key"
CREATE INDEX "api_key_robot_id_idx" ON "api_key" ("robot_id");
-- Create index "api_key_user_id_idx" to table: "api_key"
CREATE INDEX "api_key_user_id_idx" ON "api_key" ("user_id");
-- Create "episode" table
CREATE TABLE "episode" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "task_version_id" character varying(36) NOT NULL,
  "location_id" character varying(36) NOT NULL,
  "robot_id" character varying(36) NOT NULL,
  "user_id" character varying(36) NOT NULL,
  "recorded_by" character varying(36) NULL,
  "started_at" timestamptz NULL,
  "finished_at" timestamptz NULL,
  "collection_status" smallint NOT NULL DEFAULT 0,
  "error_details" jsonb NULL,
  "parameter_values" jsonb NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "episode_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "episode_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "episode_location_id_idx" to table: "episode"
CREATE INDEX "episode_location_id_idx" ON "episode" ("location_id");
-- Create index "episode_one_recording_per_robot" to table: "episode"
CREATE UNIQUE INDEX "episode_one_recording_per_robot" ON "episode" ("robot_id") WHERE (collection_status = 1);
-- Create index "episode_org_created_at_idx" to table: "episode"
CREATE INDEX "episode_org_created_at_idx" ON "episode" ("organization_id", "created_at" DESC);
-- Create index "episode_organization_id_idx" to table: "episode"
CREATE INDEX "episode_organization_id_idx" ON "episode" ("organization_id");
-- Create index "episode_robot_id_idx" to table: "episode"
CREATE INDEX "episode_robot_id_idx" ON "episode" ("robot_id");
-- Create index "episode_started_at_idx" to table: "episode"
CREATE INDEX "episode_started_at_idx" ON "episode" ("started_at");
-- Create index "episode_task_version_stats_idx" to table: "episode"
CREATE INDEX "episode_task_version_stats_idx" ON "episode" ("task_version_id", "collection_status") INCLUDE ("started_at", "finished_at") WHERE (collection_status = 3);
-- Create index "episode_user_id_idx" to table: "episode"
CREATE INDEX "episode_user_id_idx" ON "episode" ("user_id");
-- Create "episode_grade" table
CREATE TABLE "episode_grade" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "episode_id" character varying(36) NOT NULL,
  "user_id" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "grade" double precision NOT NULL,
  "comment" text NULL,
  "graded_at" timestamptz NOT NULL,
  PRIMARY KEY ("episode_id", "user_id"),
  CONSTRAINT "episode_grade_episode_id_fkey" FOREIGN KEY ("episode_id") REFERENCES "episode" ("id_natural") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "episode_grade_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "episode_grade_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "user" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "episode_grade_organization_id_idx" to table: "episode_grade"
CREATE INDEX "episode_grade_organization_id_idx" ON "episode_grade" ("organization_id");
-- Create index "episode_grade_user_id_idx" to table: "episode_grade"
CREATE INDEX "episode_grade_user_id_idx" ON "episode_grade" ("user_id");
-- Create "site" table
CREATE TABLE "site" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "name" character varying(255) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "site_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "site_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "site_organization_id_idx" to table: "site"
CREATE INDEX "site_organization_id_idx" ON "site" ("organization_id");
-- Create "location" table
CREATE TABLE "location" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "site_id" character varying(36) NOT NULL,
  "name" character varying(255) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "location_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "location_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "location_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "site" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "location_organization_id_idx" to table: "location"
CREATE INDEX "location_organization_id_idx" ON "location" ("organization_id");
-- Create "episode_stats_daily" table
CREATE TABLE "episode_stats_daily" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "location_id" character varying(36) NOT NULL,
  "robot_id" character varying(36) NOT NULL,
  "period_start" timestamptz NOT NULL,
  "total_duration_seconds" bigint NOT NULL DEFAULT 0,
  "episode_count" integer NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "episode_stats_daily_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "episode_stats_daily_location_id_fkey" FOREIGN KEY ("location_id") REFERENCES "location" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "episode_stats_daily_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "episode_stats_daily_robot_id_fkey" FOREIGN KEY ("robot_id") REFERENCES "robot" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "episode_stats_daily_org_loc_robot_period_idx" to table: "episode_stats_daily"
CREATE UNIQUE INDEX "episode_stats_daily_org_loc_robot_period_idx" ON "episode_stats_daily" ("organization_id", "location_id", "robot_id", "period_start");
-- Create index "episode_stats_daily_period_start_idx" to table: "episode_stats_daily"
CREATE INDEX "episode_stats_daily_period_start_idx" ON "episode_stats_daily" ("period_start");
-- Create "episode_stats_hourly" table
CREATE TABLE "episode_stats_hourly" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "location_id" character varying(36) NOT NULL,
  "robot_id" character varying(36) NOT NULL,
  "period_start" timestamptz NOT NULL,
  "total_duration_seconds" bigint NOT NULL DEFAULT 0,
  "episode_count" integer NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "episode_stats_hourly_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "episode_stats_hourly_location_id_fkey" FOREIGN KEY ("location_id") REFERENCES "location" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "episode_stats_hourly_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "episode_stats_hourly_robot_id_fkey" FOREIGN KEY ("robot_id") REFERENCES "robot" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "episode_stats_hourly_org_loc_robot_period_idx" to table: "episode_stats_hourly"
CREATE UNIQUE INDEX "episode_stats_hourly_org_loc_robot_period_idx" ON "episode_stats_hourly" ("organization_id", "location_id", "robot_id", "period_start");
-- Create index "episode_stats_hourly_period_start_idx" to table: "episode_stats_hourly"
CREATE INDEX "episode_stats_hourly_period_start_idx" ON "episode_stats_hourly" ("period_start");
-- Create "episode_stats_monthly" table
CREATE TABLE "episode_stats_monthly" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "location_id" character varying(36) NOT NULL,
  "robot_id" character varying(36) NOT NULL,
  "period_start" timestamptz NOT NULL,
  "total_duration_seconds" bigint NOT NULL DEFAULT 0,
  "episode_count" integer NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "episode_stats_monthly_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "episode_stats_monthly_location_id_fkey" FOREIGN KEY ("location_id") REFERENCES "location" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "episode_stats_monthly_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "episode_stats_monthly_robot_id_fkey" FOREIGN KEY ("robot_id") REFERENCES "robot" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "episode_stats_monthly_org_loc_robot_period_idx" to table: "episode_stats_monthly"
CREATE UNIQUE INDEX "episode_stats_monthly_org_loc_robot_period_idx" ON "episode_stats_monthly" ("organization_id", "location_id", "robot_id", "period_start");
-- Create index "episode_stats_monthly_period_start_idx" to table: "episode_stats_monthly"
CREATE INDEX "episode_stats_monthly_period_start_idx" ON "episode_stats_monthly" ("period_start");
-- Create "subtask" table
CREATE TABLE "subtask" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "task_version_id" character varying(36) NOT NULL,
  "order_index" integer NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "target_duration_seconds" integer NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "subtask_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "subtask_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "subtask_organization_id_idx" to table: "subtask"
CREATE INDEX "subtask_organization_id_idx" ON "subtask" ("organization_id");
-- Create "episode_sub_task" table
CREATE TABLE "episode_sub_task" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "episode_id" character varying(36) NOT NULL,
  "sub_task_id" character varying(36) NOT NULL,
  "collection_status" smallint NOT NULL DEFAULT 0,
  "task_result" smallint NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "episode_sub_task_episode_sub_task_unique" UNIQUE ("episode_id", "sub_task_id"),
  CONSTRAINT "episode_sub_task_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "episode_sub_task_episode_id_fkey" FOREIGN KEY ("episode_id") REFERENCES "episode" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "episode_sub_task_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "episode_sub_task_sub_task_id_fkey" FOREIGN KEY ("sub_task_id") REFERENCES "subtask" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "episode_sub_task_execution" table
CREATE TABLE "episode_sub_task_execution" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "episode_sub_task_id" character varying(36) NOT NULL,
  "execution_status" smallint NOT NULL DEFAULT 0,
  "started_at" timestamptz NULL,
  "finished_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "episode_sub_task_execution_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "episode_sub_task_execution_episode_sub_task_id_fkey" FOREIGN KEY ("episode_sub_task_id") REFERENCES "episode_sub_task" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "episode_sub_task_execution_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "episode_sub_task_execution_episode_sub_task_id_created_at_idx" to table: "episode_sub_task_execution"
CREATE INDEX "episode_sub_task_execution_episode_sub_task_id_created_at_idx" ON "episode_sub_task_execution" ("episode_sub_task_id", "created_at");
-- Create "robot_uptime_hourly" table
CREATE TABLE "robot_uptime_hourly" (
  "robot_id" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "location_id" character varying(36) NOT NULL,
  "period_start" timestamptz NOT NULL,
  "uptime_seconds" bigint NOT NULL DEFAULT 0,
  PRIMARY KEY ("robot_id", "period_start"),
  CONSTRAINT "robot_uptime_hourly_location_id_fkey" FOREIGN KEY ("location_id") REFERENCES "location" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "robot_uptime_hourly_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "robot_uptime_hourly_robot_id_fkey" FOREIGN KEY ("robot_id") REFERENCES "robot" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "robot_uptime_hourly_period_start_idx" to table: "robot_uptime_hourly"
CREATE INDEX "robot_uptime_hourly_period_start_idx" ON "robot_uptime_hourly" ("period_start");
-- Create "task" table
CREATE TABLE "task" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "manual_url" text NULL,
  "priority" smallint NOT NULL,
  "difficulty" smallint NOT NULL,
  "status" smallint NOT NULL DEFAULT 0,
  "deadline" timestamptz NOT NULL,
  "robot_type" character varying(255) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "task_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "task_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "task_organization_id_idx" to table: "task"
CREATE INDEX "task_organization_id_idx" ON "task" ("organization_id");
-- Create "task_category_type" table
CREATE TABLE "task_category_type" (
  "id" character varying(36) NOT NULL,
  "slug" character varying(100) NOT NULL,
  "name" character varying(255) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "task_tag" table
CREATE TABLE "task_tag" (
  "id" character varying(36) NOT NULL,
  "name" character varying(255) NOT NULL,
  "category_type_id" character varying(36) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "task_tag_category_type_id_fkey" FOREIGN KEY ("category_type_id") REFERENCES "task_category_type" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "task_tag_assignment" table
CREATE TABLE "task_tag_assignment" (
  "task_id" character varying(36) NOT NULL,
  "tag_id" character varying(36) NOT NULL,
  CONSTRAINT "task_tag_assignment_tag_id_fkey" FOREIGN KEY ("tag_id") REFERENCES "task_tag" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "task_tag_assignment_task_id_fkey" FOREIGN KEY ("task_id") REFERENCES "task" ("id_natural") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "task_version" table
CREATE TABLE "task_version" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  "task_id" character varying(36) NOT NULL,
  "version" character varying(50) NOT NULL,
  "schema_hash" character varying(255) NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "approval_status" integer NOT NULL DEFAULT 0,
  "target_duration_seconds" integer NULL,
  "target_episode_count" integer NULL,
  "target_duration_per_episode_seconds" integer NULL,
  "display_name" character varying(100) NULL,
  "parameters" jsonb NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "task_version_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "task_version_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "task_version_organization_id_idx" to table: "task_version"
CREATE INDEX "task_version_organization_id_idx" ON "task_version" ("organization_id");
-- Create "task_version_stats" table
CREATE TABLE "task_version_stats" (
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "id" bigserial NOT NULL,
  "id_natural" character varying(36) NOT NULL,
  "task_version_id" character varying(36) NOT NULL,
  "total_duration_seconds" bigint NOT NULL DEFAULT 0,
  "episode_count" integer NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "task_version_stats_id_natural_key" UNIQUE ("id_natural"),
  CONSTRAINT "task_version_stats_task_version_id_key" UNIQUE ("task_version_id"),
  CONSTRAINT "task_version_stats_task_version_id_fkey" FOREIGN KEY ("task_version_id") REFERENCES "task_version" ("id_natural") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "user_location_assignment" table
CREATE TABLE "user_location_assignment" (
  "user_id" character varying(36) NOT NULL,
  "location_id" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  PRIMARY KEY ("user_id", "location_id"),
  CONSTRAINT "user_location_assignment_location_id_fkey" FOREIGN KEY ("location_id") REFERENCES "location" ("id_natural") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "user_location_assignment_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "user_location_assignment_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "user" ("id_natural") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_user_location_assignment_location_id" to table: "user_location_assignment"
CREATE INDEX "idx_user_location_assignment_location_id" ON "user_location_assignment" ("location_id");
-- Create "user_site_assignment" table
CREATE TABLE "user_site_assignment" (
  "user_id" character varying(36) NOT NULL,
  "site_id" character varying(36) NOT NULL,
  "organization_id" character varying(36) NOT NULL,
  PRIMARY KEY ("user_id", "site_id"),
  CONSTRAINT "user_site_assignment_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organization" ("id_natural") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "user_site_assignment_site_id_fkey" FOREIGN KEY ("site_id") REFERENCES "site" ("id_natural") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "user_site_assignment_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "user" ("id_natural") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_user_site_assignment_site_id" to table: "user_site_assignment"
CREATE INDEX "idx_user_site_assignment_site_id" ON "user_site_assignment" ("site_id");
