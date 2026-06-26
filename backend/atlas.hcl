// Atlas configuration file
// https://atlasgo.io/atlas-schema/projects

variable "db_url" {
  type    = string
  default = getenv("DATABASE_URL")
}

env "local" {
  // Source of truth: SQL schema file
  src = "file://internal/infra/database/schema/schema.up.sql"

  // Target database URL
  url = var.db_url

  // Dev database for computing diffs (auto-managed Docker container)
  dev = "docker://postgres/17/dev?search_path=public"

  migration {
    // Directory containing migration files
    dir = "file://internal/infra/database/migrate"
  }

  diff {
    // Skip destructive changes by default (optional, for safety)
    // skip {
    //   drop_table = true
    // }
  }

  format {
    migrate {
      // Use SQL format for migration files
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "dev" {
  // Same as local but with explicit connection
  src = "file://internal/infra/database/schema/schema.up.sql"
  url = "postgres://postgres:postgres@localhost:5432/airoa?sslmode=disable"
  dev = "docker://postgres/17/dev?search_path=public"

  migration {
    dir = "file://internal/infra/database/migrate"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "prod" {
  // Production environment using environment variables
  url = "postgres://${getenv("DB_USER")}:${getenv("DB_PASSWORD")}@${getenv("DB_HOST")}:${getenv("DB_PORT")}/${getenv("DB_NAME")}?sslmode=require"

  migration {
    dir = "file://migrate"
  }
}
