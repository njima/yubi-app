import { makeApi, Zodios, type ZodiosOptions } from "@zodios/core";
import { z } from "zod";

const UserRole = z.union([
  z.literal(0),
  z.literal(1),
  z.literal(2),
  z.literal(3),
  z.literal(4),
]);
const LocationSummary = z
  .object({ location_id: z.string(), name: z.string() })
  .passthrough();
const SiteSummary = z
  .object({ site_id: z.string(), name: z.string() })
  .passthrough();
const UserResponse = z
  .object({
    user_id: z.string(),
    email: z.string(),
    display_name: z.string(),
    role: UserRole.optional(),
    created_at: z.string().datetime({ offset: true }),
    organization_id: z.string(),
    organization_name: z.string(),
    updated_at: z.string().datetime({ offset: true }).optional(),
    locations: z.array(LocationSummary),
    sites: z.array(SiteSummary),
  })
  .passthrough();
const ErrorResponse = z
  .object({ code: z.number().int(), message: z.string() })
  .passthrough();
const MeUpdateRequest = z
  .object({ display_name: z.string().min(1).max(60) })
  .passthrough();
const Pagination = z
  .object({
    count: z.number().int(),
    page: z.number().int(),
    limit: z.number().int(),
  })
  .passthrough();
const UserFilter = z
  .object({
    organization_id: z.string(),
    location_id: z.string(),
    site_id: z.string(),
  })
  .partial()
  .passthrough();
const UserListResponse = z
  .object({
    users: z.array(UserResponse),
    pagination: Pagination,
    filter: UserFilter,
  })
  .passthrough();
const UserCreateRequest = z
  .object({
    email: z.string(),
    display_name: z.string(),
    role: UserRole,
    location_ids: z.array(z.string()).optional(),
    site_ids: z.array(z.string()).optional(),
  })
  .passthrough();
const UserUpdateRequest = z
  .object({ email: z.string(), display_name: z.string() })
  .partial()
  .passthrough();
const UserRoleUpdateRequest = z.object({ role: UserRole }).passthrough();
const UserLocationsUpdateRequest = z
  .object({ location_ids: z.array(z.string()) })
  .passthrough();
const UserImportRequest = z.object({ csv_content: z.string() }).passthrough();
const UserImportValidRow = z
  .object({
    row_number: z.number().int(),
    email: z.string(),
    display_name: z.string(),
    role: z.string(),
  })
  .passthrough();
const UserImportRowError = z
  .object({
    row_number: z.number().int(),
    email: z.string().optional(),
    errors: z.array(z.string()),
  })
  .passthrough();
const UserImportSummary = z
  .object({
    valid_count: z.number().int(),
    duplicate_count: z.number().int(),
    error_count: z.number().int(),
  })
  .passthrough();
const UserImportValidationResponse = z
  .object({
    valid_rows: z.array(UserImportValidRow),
    duplicate_rows: z.array(UserImportRowError),
    error_rows: z.array(UserImportRowError),
    summary: UserImportSummary,
  })
  .passthrough();
const UserImportResponse = z
  .object({
    imported_count: z.number().int(),
    skipped_count: z.number().int(),
    error_count: z.number().int(),
    errors: z.array(UserImportRowError),
  })
  .passthrough();
const OrganizationResponse = z
  .object({
    organization_id: z.string(),
    display_name: z.string(),
    description: z.string().optional(),
    created_at: z.string().datetime({ offset: true }).optional(),
    updated_at: z.string().datetime({ offset: true }).optional(),
  })
  .passthrough();
const OrganizationListResponse = z
  .object({
    organizations: z.array(OrganizationResponse),
    pagination: Pagination,
  })
  .passthrough();
const OrganizationCreateRequest = z
  .object({
    display_name: z.string().min(2).max(100),
    description: z.string().optional(),
  })
  .passthrough();
const OrganizationUpdateRequest = z
  .object({ display_name: z.string(), description: z.string() })
  .partial()
  .passthrough();
const Permission = z.object({ code: z.string() }).passthrough();
const PermissionCreate = z.object({ code: z.string() }).passthrough();
const FleetStatusCount = z
  .object({ operational: z.number().int(), total: z.number().int() })
  .passthrough();
const FleetRobotTypeSummary = z
  .object({ leader: FleetStatusCount.optional(), follower: FleetStatusCount })
  .passthrough();
const FleetSiteSummary = z
  .object({
    site: z.string(),
    site_id: z.string(),
    robot_types: z.record(FleetRobotTypeSummary),
  })
  .passthrough();
const FleetRobotTypeStats = z
  .object({
    robot_type: z.string(),
    robot_uptime: z.number().nullish(),
    uptime_rate: z.number().gte(0).lte(1).nullish(),
    robot_count: z.number().int().nullish(),
    data_collection_time: z.number(),
  })
  .passthrough();
const FleetSiteStats = z
  .object({
    site: z.string(),
    site_id: z.string(),
    robot_types: z.array(FleetRobotTypeStats),
  })
  .passthrough();
const TrendSeries = z
  .object({ label: z.string(), data: z.array(z.number()) })
  .passthrough();
const CollectionTrend = z
  .object({
    labels: z.array(z.string()),
    by_site: z.array(TrendSeries),
    by_robot_type: z.array(TrendSeries),
  })
  .passthrough();
const RobotStatus = z.union([
  z.literal(0),
  z.literal(1),
  z.literal(2),
  z.literal(3),
  z.literal(4),
  z.literal(5),
]);
const LeaderStatus = z.union([z.literal(0), z.literal(1), z.literal(2)]);
const RobotOperator = z
  .object({
    user_id: z.string(),
    display_name: z.string(),
    organization_name: z.string(),
  })
  .passthrough();
const Robot = z
  .object({
    id: z.string(),
    organization_id: z.string().optional(),
    organization_name: z.string().optional(),
    site_id: z.string().optional(),
    site_name: z.string().optional(),
    location_id: z.string().optional(),
    location_name: z.string().optional(),
    name: z.string(),
    robot_type: z.string().optional(),
    status: RobotStatus.optional(),
    leader_status: LeaderStatus.optional(),
    consecutive_fault_days: z.number().int().nullish(),
    leader_consecutive_fault_days: z.number().int().nullish(),
    leader_fault_started_at: z.string().datetime({ offset: true }).nullish(),
    battery_level: z.number().int().optional(),
    last_heartbeat_at: z.string().datetime({ offset: true }).optional(),
    offline_reason: z.string().optional(),
    robot_config: z.object({}).partial().passthrough().nullish(),
    active_episode_id: z.string().optional(),
    active_user_id: z.string().optional(),
    active_operator: RobotOperator.optional(),
  })
  .passthrough();
const RobotListResponse = z
  .object({ robots: z.array(Robot), pagination: Pagination })
  .passthrough();
const RobotCreate = z
  .object({
    organization_id: z.string(),
    location_id: z.string(),
    name: z.string(),
    robot_type: z.string().optional(),
    leader_status: LeaderStatus.optional(),
    robot_config: z.object({}).partial().passthrough().nullish(),
  })
  .passthrough();
const RobotUpdate = z
  .object({
    robot_id: z.string(),
    organization_id: z.string(),
    location_id: z.string(),
    name: z.string(),
    robot_type: z.string(),
    status: RobotStatus,
    leader_status: LeaderStatus,
    battery_level: z.number().int(),
    last_heartbeat_at: z.string().datetime({ offset: true }),
    offline_reason: z.string(),
    robot_config: z.object({}).partial().passthrough().nullable(),
  })
  .partial()
  .passthrough();
const RobotOperatorRequest = z
  .object({
    user_id: z.string(),
    display_name: z.string(),
    organization_name: z.string(),
  })
  .passthrough();
const TaskStatus = z.union([
  z.literal(0),
  z.literal(1),
  z.literal(2),
  z.literal(3),
]);
const TaskPriority = z.union([
  z.literal(0),
  z.literal(1),
  z.literal(2),
  z.literal(3),
]);
const TaskDifficulty = z.union([
  z.literal(0),
  z.literal(1),
  z.literal(2),
  z.literal(3),
]);
const TaskTag = z
  .object({
    id: z.string(),
    name: z.string(),
    category_type_id: z.string(),
    category_type_name: z.string(),
  })
  .passthrough();
const Task = z
  .object({
    id: z.string(),
    name: z.string(),
    description: z.string().optional(),
    version: z.string().optional(),
    version_display_name: z.string().optional(),
    manual_url: z.string(),
    priority: TaskPriority,
    difficulty: TaskDifficulty,
    status: TaskStatus,
    deadline: z.string().datetime({ offset: true }),
    robot_type: z.string().optional(),
    target_duration_seconds: z.number().int().optional(),
    target_episode_count: z.number().int().optional(),
    actual_episode_count: z.number().int().optional(),
    tags: z.array(TaskTag).optional(),
  })
  .passthrough();
const TaskListResponse = z
  .object({ tasks: z.array(Task), pagination: Pagination })
  .passthrough();
const TaskCreate = z
  .object({
    organization_id: z.string(),
    name: z.string(),
    description: z.string().optional(),
    manual_url: z.string(),
    priority: TaskPriority,
    difficulty: TaskDifficulty,
    status: TaskStatus,
    deadline: z.string().datetime({ offset: true }),
    robot_type: z.string().optional(),
    tag_ids: z.array(z.string()).optional(),
  })
  .passthrough();
const TaskSummaryResponse = z
  .object({
    total_tasks: z.number().int(),
    target_duration_seconds: z.number().int(),
    target_episode_count: z.number().int(),
  })
  .passthrough();
const TaskTrendGroup = z
  .object({
    label: z.string(),
    target_tasks: z.number().int(),
    actual_tasks: z.number().int(),
    target_duration: z.number().int(),
    actual_duration: z.number().int(),
    target_episodes: z.number().int(),
    actual_episodes: z.number().int(),
  })
  .passthrough();
const TaskTrendPeriod = z
  .object({
    start: z.string().datetime({ offset: true }).optional(),
    end: z.string().datetime({ offset: true }).optional(),
    groups: z.array(TaskTrendGroup),
  })
  .passthrough();
const TaskCompletionTrend = z
  .object({ periods: z.array(TaskTrendPeriod) })
  .passthrough();
const TaskImportRequest = z.object({ csv_content: z.string() }).passthrough();
const TaskImportRow = z
  .object({
    row_number: z.number().int(),
    name: z.string(),
    description: z.string().optional(),
    manual_url: z.string(),
    priority: z.string(),
    difficulty: z.string(),
    status: z.string().optional(),
    deadline: z.string(),
    robot_type: z.string().optional(),
    tags: z.string().optional(),
  })
  .passthrough();
const TaskImportRowError = z
  .object({
    row_number: z.number().int(),
    errors: z.array(z.string()),
    name: z.string().optional(),
  })
  .passthrough();
const TaskImportSummary = z
  .object({
    valid_count: z.number().int(),
    duplicate_count: z.number().int(),
    error_count: z.number().int(),
  })
  .passthrough();
const TaskImportValidationResponse = z
  .object({
    valid_rows: z.array(TaskImportRow),
    duplicate_rows: z.array(TaskImportRowError),
    error_rows: z.array(TaskImportRowError),
    summary: TaskImportSummary,
  })
  .passthrough();
const TaskImportResponse = z
  .object({
    imported_count: z.number().int(),
    skipped_count: z.number().int(),
    error_count: z.number().int(),
    errors: z.array(TaskImportRowError),
  })
  .passthrough();
const TaskUpdate = z
  .object({
    name: z.string(),
    description: z.string(),
    manual_url: z.string(),
    priority: TaskPriority,
    difficulty: TaskDifficulty,
    status: TaskStatus,
    deadline: z.string().datetime({ offset: true }),
    robot_type: z.string(),
    tag_ids: z.array(z.string()),
  })
  .partial()
  .passthrough();
const TaskCategoryType = z
  .object({ id: z.string(), slug: z.string(), name: z.string() })
  .passthrough();
const TaskTagCreate = z
  .object({ name: z.string(), category_type_id: z.string() })
  .passthrough();
const ApprovalStatus = z.union([z.literal(0), z.literal(1)]);
const TaskVersionParameter = z
  .object({ key: z.string().min(1), values: z.array(z.string().min(1)).min(1) })
  .passthrough();
const TaskVersion = z
  .object({
    id: z.string(),
    task_id: z.string(),
    version: z.string(),
    display_name: z.string().max(100).nullish(),
    is_current: z.boolean(),
    approval_status: ApprovalStatus,
    parameters: z.array(TaskVersionParameter).optional(),
    created_at: z.string().datetime({ offset: true }),
    target_duration_seconds: z.number().int().gte(1).optional(),
    target_episode_count: z.number().int().gte(1).optional(),
    target_duration_per_episode_seconds: z.number().int().gte(1).optional(),
    actual_duration_seconds: z.number().int().optional(),
    actual_episode_count: z.number().int().optional(),
  })
  .passthrough();
const TaskVersionCreate = z
  .object({
    version: z.string().min(1).max(50),
    display_name: z.string().max(100).optional(),
    base_task_version_id: z.string().min(1),
    target_duration_seconds: z.number().int().gte(1).optional(),
    target_episode_count: z.number().int().gte(1).optional(),
    target_duration_per_episode_seconds: z.number().int().gte(1).optional(),
    parameters: z.array(TaskVersionParameter).optional(),
  })
  .passthrough();
const TaskVersionUpdate = z
  .object({
    display_name: z.string().max(100),
    target_duration_seconds: z.number().int().gte(1),
    target_episode_count: z.number().int().gte(1),
    target_duration_per_episode_seconds: z.number().int().gte(1),
  })
  .partial()
  .passthrough();
const TaskVersionParametersUpdate = z
  .object({ parameters: z.array(TaskVersionParameter) })
  .passthrough();
const SubTaskStatus = z.union([
  z.literal(0),
  z.literal(1),
  z.literal(2),
  z.literal(3),
]);
const SubTask = z
  .object({
    id: z.string(),
    status: SubTaskStatus.optional(),
    name: z.string(),
    description: z.string().optional(),
    target_duration_seconds: z.number().int().gte(1).optional(),
  })
  .passthrough();
const SubTaskCreate = z
  .object({
    organization_id: z.string(),
    name: z.string(),
    description: z.string().optional(),
    task_id: z.string(),
    task_version_id: z.string(),
    target_duration_seconds: z.number().int().gte(1).optional(),
  })
  .passthrough();
const SubTaskReorder = z
  .object({ task_version_id: z.string(), subtask_ids: z.array(z.string()) })
  .passthrough();
const SubTaskUpdate = z
  .object({
    name: z.string(),
    description: z.string(),
    is_completed: z.boolean(),
    target_duration_seconds: z.number().int().gte(1),
  })
  .partial()
  .passthrough();
const EpisodeCollectionStatus = z.union([
  z.literal(0),
  z.literal(1),
  z.literal(2),
  z.literal(3),
]);
const SubTaskCollectionStatus = z.union([
  z.literal(0),
  z.literal(1),
  z.literal(2),
  z.literal(3),
  z.literal(4),
]);
const ExecutionStatus = z.union([
  z.literal(0),
  z.literal(1),
  z.literal(2),
  z.literal(3),
]);
const EpisodeSubTaskExecution = z
  .object({
    id: z.string(),
    status: ExecutionStatus,
    started_at: z.string().datetime({ offset: true }).optional(),
    finished_at: z.string().datetime({ offset: true }).optional(),
  })
  .passthrough();
const EpisodeSubTask = z
  .object({
    id: z.string(),
    subtask_id: z.string(),
    name: z.string(),
    order_index: z.number().int(),
    status: SubTaskCollectionStatus,
    executions: z.array(EpisodeSubTaskExecution).optional(),
  })
  .passthrough();
const Episode = z
  .object({
    id: z.string(),
    user_id: z.string(),
    location_id: z.string(),
    robot_id: z.string(),
    task_id: z.string(),
    task_name: z.string().optional(),
    task_description: z.string().optional(),
    task_version_id: z.string(),
    task_version_display_name: z.string().optional(),
    status: EpisodeCollectionStatus,
    error_details: z.string().optional(),
    started_at: z.string().datetime({ offset: true }).optional(),
    ended_at: z.string().datetime({ offset: true }).optional(),
    created_at: z.string().datetime({ offset: true }),
    recorded_by: z.string().optional(),
    parameter_values: z.record(z.string()).optional(),
    subtasks: z.array(EpisodeSubTask).optional(),
    average_grade: z.number().gte(0).lte(1).nullish(),
    grade_count: z.number().int().gte(0).optional(),
  })
  .passthrough();
const EpisodeListResponse = z
  .object({ episodes: z.array(Episode), pagination: Pagination })
  .passthrough();
const EpisodeCreate = z
  .object({
    organization_id: z.string(),
    location_id: z.string(),
    robot_id: z.string(),
    task_id: z.string(),
    recorded_by: z.string().optional(),
    task_version_id: z.string().optional(),
    parameter_values: z.record(z.string()).optional(),
  })
  .passthrough();
const EpisodeBulkCreate = EpisodeCreate.and(
  z.object({ count: z.number().int().gte(1).lte(100) }).passthrough()
);
const EpisodeUpdate = z
  .object({
    start_time: z.string().datetime({ offset: true }),
    end_time: z.string().datetime({ offset: true }),
    status: EpisodeCollectionStatus,
    error_details: z.string(),
    recorded_by: z.string(),
  })
  .partial()
  .passthrough();
const EpisodeRecordingsResponse = z
  .object({ recordings: z.record(z.string()) })
  .passthrough();
const EpisodeFeatureStats = z
  .object({
    min: z.unknown(),
    max: z.unknown(),
    mean: z.unknown(),
    std: z.unknown(),
    count: z.number().int(),
  })
  .passthrough();
const EpisodeStatsResponse = z
  .object({ stats: z.record(EpisodeFeatureStats) })
  .passthrough();
const EpisodeGrade = z
  .object({
    episode_id: z.string(),
    user_id: z.string(),
    user_name: z.string(),
    grade: z.number().gte(0).lte(1),
    comment: z.string().max(10000).nullish(),
    graded_at: z.string().datetime({ offset: true }),
    created_at: z.string().datetime({ offset: true }),
    updated_at: z.string().datetime({ offset: true }).optional(),
  })
  .passthrough();
const EpisodeGradeListResponse = z
  .object({ grades: z.array(EpisodeGrade), pagination: Pagination })
  .passthrough();
const EpisodeGradeUpdate = z
  .object({
    grade: z.number().gte(0).lte(1),
    comment: z.string().max(10000).nullish(),
  })
  .passthrough();
const UserSitesUpdateRequest = z
  .object({ site_ids: z.array(z.string()) })
  .passthrough();
const Site = z
  .object({
    id: z.string(),
    name: z.string(),
    organization_id: z.string(),
    created_at: z.string().datetime({ offset: true }).optional(),
    updated_at: z.string().datetime({ offset: true }).optional(),
  })
  .passthrough();
const SiteListResponse = z
  .object({ sites: z.array(Site), pagination: Pagination })
  .passthrough();
const SiteCreate = z
  .object({ organization_id: z.string(), name: z.string().min(1).max(100) })
  .passthrough();
const SiteUpdate = z.object({ name: z.string() }).partial().passthrough();
const Location = z
  .object({
    id: z.string(),
    name: z.string(),
    site_id: z.string(),
    site_name: z.string(),
  })
  .passthrough();
const LocationListResponse = z
  .object({ locations: z.array(Location), pagination: Pagination })
  .passthrough();
const LocationCreate = z
  .object({
    name: z.string(),
    organization_id: z.string(),
    site_id: z.string(),
  })
  .passthrough();
const LocationUpdate = z.object({ name: z.string() }).partial().passthrough();
const EpisodeActionRequest = z
  .object({ occurred_at: z.string().datetime({ offset: true }) })
  .passthrough();
const ExecutionCreateResponse = z
  .object({ execution_id: z.string() })
  .passthrough();
const ExecutionActionRequest = z
  .object({ occurred_at: z.string().datetime({ offset: true }) })
  .passthrough();
const BatteryStatus = z
  .object({ pct: z.number().int(), charging: z.boolean() })
  .passthrough();
const ConnectionStatus = z
  .object({ quality_pct: z.number().int() })
  .passthrough();
const RobotMetric = z
  .object({
    name: z.string(),
    type: z.enum(["scalar", "vector", "map", "array"]),
    unit: z.string(),
    value: z.unknown(),
    labels: z.record(z.string()).optional(),
  })
  .passthrough();
const GateCondition = z
  .object({
    name: z.string(),
    passed: z.boolean(),
    reason: z.string(),
    escalation: z.number().int(),
  })
  .passthrough();
const GateGroupStatus = z
  .object({
    level: z.number().int(),
    settled: z.boolean(),
    conditions: z.array(GateCondition),
  })
  .passthrough();
const GateConditionStatus = z
  .object({ gate_level: z.number().int(), groups: z.record(GateGroupStatus) })
  .passthrough();
const RobotStatusDetail = z
  .object({
    battery: BatteryStatus,
    connection: ConnectionStatus,
    uptime_sec: z.number(),
    metrics: z.array(RobotMetric).optional(),
    gate_conditions: GateConditionStatus.optional(),
  })
  .passthrough();
const RobotStatusUpdateRequest = z
  .object({
    robot_type: z.string(),
    reported_at: z.string().datetime({ offset: true }),
    status: RobotStatusDetail,
  })
  .passthrough();
const ApiKeyResponse = z
  .object({
    id: z.string(),
    name: z.string(),
    user_id: z.string(),
    user_name: z.string(),
    robot_id: z.string().nullish(),
    robot_name: z.string().nullish(),
    organization_id: z.string(),
    key_hint: z.string(),
    expires_at: z.string().datetime({ offset: true }).nullish(),
    last_used_at: z.string().datetime({ offset: true }).nullish(),
    revoked_at: z.string().datetime({ offset: true }).nullish(),
    created_at: z.string().datetime({ offset: true }),
    updated_at: z.string().datetime({ offset: true }),
  })
  .passthrough();
const ApiKeyListResponse = z
  .object({ api_keys: z.array(ApiKeyResponse), pagination: Pagination })
  .passthrough();
const ApiKeyCreateRequest = z
  .object({
    name: z.string().min(1).max(255),
    robot_id: z.string(),
    expires_at: z.string().datetime({ offset: true }).nullish(),
  })
  .passthrough();
const ApiKeyCreateResponse = ApiKeyResponse.and(
  z.object({ key: z.string() }).passthrough()
);
const ApiKeyUpdateRequest = z
  .object({
    name: z.string().min(1).max(255),
    expires_at: z.string().datetime({ offset: true }).nullable(),
    clear_expires_at: z.boolean(),
  })
  .partial()
  .passthrough();
const EpisodeTaskResult = z.union([z.literal(0), z.literal(1), z.literal(2)]);
const EpisodeUploadStatus = z.union([
  z.literal(0),
  z.literal(1),
  z.literal(2),
  z.literal(3),
]);
const SubTaskTaskResult = z.union([z.literal(0), z.literal(1), z.literal(2)]);
const SubTaskActionRequest = z
  .object({ occurred_at: z.string().datetime({ offset: true }) })
  .passthrough();
const PermissionDelete = z.object({ code: z.string() }).passthrough();
const RobotStatusStreamDetail = z
  .object({
    battery_pct: z.number().int(),
    connection_pct: z.number().int(),
    uptime_sec: z.number().int(),
    gate_conditions: GateConditionStatus.optional(),
  })
  .passthrough();
const RobotStatusStreamResponse = z
  .object({
    robot_id: z.string(),
    robot_type: z.string(),
    status: RobotStatusStreamDetail,
  })
  .passthrough();

export const schemas = {
  UserRole,
  LocationSummary,
  SiteSummary,
  UserResponse,
  ErrorResponse,
  MeUpdateRequest,
  Pagination,
  UserFilter,
  UserListResponse,
  UserCreateRequest,
  UserUpdateRequest,
  UserRoleUpdateRequest,
  UserLocationsUpdateRequest,
  UserImportRequest,
  UserImportValidRow,
  UserImportRowError,
  UserImportSummary,
  UserImportValidationResponse,
  UserImportResponse,
  OrganizationResponse,
  OrganizationListResponse,
  OrganizationCreateRequest,
  OrganizationUpdateRequest,
  Permission,
  PermissionCreate,
  FleetStatusCount,
  FleetRobotTypeSummary,
  FleetSiteSummary,
  FleetRobotTypeStats,
  FleetSiteStats,
  TrendSeries,
  CollectionTrend,
  RobotStatus,
  LeaderStatus,
  RobotOperator,
  Robot,
  RobotListResponse,
  RobotCreate,
  RobotUpdate,
  RobotOperatorRequest,
  TaskStatus,
  TaskPriority,
  TaskDifficulty,
  TaskTag,
  Task,
  TaskListResponse,
  TaskCreate,
  TaskSummaryResponse,
  TaskTrendGroup,
  TaskTrendPeriod,
  TaskCompletionTrend,
  TaskImportRequest,
  TaskImportRow,
  TaskImportRowError,
  TaskImportSummary,
  TaskImportValidationResponse,
  TaskImportResponse,
  TaskUpdate,
  TaskCategoryType,
  TaskTagCreate,
  ApprovalStatus,
  TaskVersionParameter,
  TaskVersion,
  TaskVersionCreate,
  TaskVersionUpdate,
  TaskVersionParametersUpdate,
  SubTaskStatus,
  SubTask,
  SubTaskCreate,
  SubTaskReorder,
  SubTaskUpdate,
  EpisodeCollectionStatus,
  SubTaskCollectionStatus,
  ExecutionStatus,
  EpisodeSubTaskExecution,
  EpisodeSubTask,
  Episode,
  EpisodeListResponse,
  EpisodeCreate,
  EpisodeBulkCreate,
  EpisodeUpdate,
  EpisodeRecordingsResponse,
  EpisodeFeatureStats,
  EpisodeStatsResponse,
  EpisodeGrade,
  EpisodeGradeListResponse,
  EpisodeGradeUpdate,
  UserSitesUpdateRequest,
  Site,
  SiteListResponse,
  SiteCreate,
  SiteUpdate,
  Location,
  LocationListResponse,
  LocationCreate,
  LocationUpdate,
  EpisodeActionRequest,
  ExecutionCreateResponse,
  ExecutionActionRequest,
  BatteryStatus,
  ConnectionStatus,
  RobotMetric,
  GateCondition,
  GateGroupStatus,
  GateConditionStatus,
  RobotStatusDetail,
  RobotStatusUpdateRequest,
  ApiKeyResponse,
  ApiKeyListResponse,
  ApiKeyCreateRequest,
  ApiKeyCreateResponse,
  ApiKeyUpdateRequest,
  EpisodeTaskResult,
  EpisodeUploadStatus,
  SubTaskTaskResult,
  SubTaskActionRequest,
  PermissionDelete,
  RobotStatusStreamDetail,
  RobotStatusStreamResponse,
};

const endpoints = makeApi([
  {
    method: "get",
    path: "/api-keys",
    alias: "listApiKeys",
    description: `List API keys in the caller&#x27;s organization. Admin only.`,
    requestFormat: "json",
    parameters: [
      {
        name: "limit",
        type: "Query",
        schema: z.number().int().gte(1).lte(1000).optional().default(50),
      },
      {
        name: "page",
        type: "Query",
        schema: z.number().int().gte(1).optional().default(1),
      },
      {
        name: "robot_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "user_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "include_revoked",
        type: "Query",
        schema: z.boolean().optional().default(false),
      },
    ],
    response: ApiKeyListResponse,
    errors: [
      {
        status: 401,
        description: `Unauthorized`,
        schema: ErrorResponse,
      },
      {
        status: 403,
        description: `Forbidden`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/api-keys",
    alias: "createApiKey",
    description: `Create a new API key for a robot in the caller&#x27;s organization. The raw
key is returned exactly once in the response and cannot be retrieved
afterwards.
`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: ApiKeyCreateRequest,
      },
    ],
    response: ApiKeyCreateResponse,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 401,
        description: `Unauthorized`,
        schema: ErrorResponse,
      },
      {
        status: 403,
        description: `Forbidden`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/api-keys/:api_key_id",
    alias: "getApiKey",
    requestFormat: "json",
    parameters: [
      {
        name: "api_key_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: ApiKeyResponse,
    errors: [
      {
        status: 401,
        description: `Unauthorized`,
        schema: ErrorResponse,
      },
      {
        status: 403,
        description: `Forbidden`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "patch",
    path: "/api-keys/:api_key_id",
    alias: "updateApiKey",
    description: `Update name and/or expires_at on an API key. The robot_id binding is
immutable — to change it, revoke the key and issue a new one.
`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: ApiKeyUpdateRequest,
      },
      {
        name: "api_key_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: ApiKeyResponse,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 401,
        description: `Unauthorized`,
        schema: ErrorResponse,
      },
      {
        status: 403,
        description: `Forbidden`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "delete",
    path: "/api-keys/:api_key_id",
    alias: "revokeApiKey",
    description: `Soft-delete the key by setting revoked_at. Idempotent.`,
    requestFormat: "json",
    parameters: [
      {
        name: "api_key_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 401,
        description: `Unauthorized`,
        schema: ErrorResponse,
      },
      {
        status: 403,
        description: `Forbidden`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/episodes",
    alias: "listEpisodes",
    requestFormat: "json",
    parameters: [
      {
        name: "task_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "task_version_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "robot_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "user_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "status",
        type: "Query",
        schema: z.array(EpisodeCollectionStatus).optional(),
      },
      {
        name: "started_at_from",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "started_at_to",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "sort_by",
        type: "Query",
        schema: z
          .enum([
            "task",
            "robot",
            "recorded_by",
            "started_at",
            "ended_at",
            "error",
          ])
          .optional(),
      },
      {
        name: "sort_order",
        type: "Query",
        schema: z.enum(["asc", "desc"]).optional(),
      },
      {
        name: "page",
        type: "Query",
        schema: z.number().int().optional(),
      },
      {
        name: "limit",
        type: "Query",
        schema: z.number().int().optional(),
      },
    ],
    response: EpisodeListResponse,
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/episodes",
    alias: "createEpisode",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the episode to create`,
        type: "Body",
        schema: EpisodeCreate,
      },
    ],
    response: Episode,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/episodes/:episodeId",
    alias: "getEpisodeById",
    requestFormat: "json",
    parameters: [
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: Episode,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/episodes/:episodeId",
    alias: "updateEpisodeById",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the episode to update`,
        type: "Body",
        schema: EpisodeUpdate,
      },
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: Episode,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "delete",
    path: "/episodes/:episodeId",
    alias: "deleteEpisodeById",
    requestFormat: "json",
    parameters: [
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/episodes/:episodeId/grades",
    alias: "listEpisodeGrades",
    requestFormat: "json",
    parameters: [
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "page",
        type: "Query",
        schema: z.number().int().gte(1).optional().default(1),
      },
      {
        name: "limit",
        type: "Query",
        schema: z.number().int().gte(1).lte(1000).optional().default(20),
      },
    ],
    response: EpisodeGradeListResponse,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/episodes/:episodeId/grades/me",
    alias: "getMyEpisodeGrade",
    requestFormat: "json",
    parameters: [
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: EpisodeGrade,
    errors: [
      {
        status: 404,
        description: `Current user has not graded this episode yet`,
        schema: z.void(),
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/episodes/:episodeId/grades/me",
    alias: "updateMyEpisodeGrade",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: EpisodeGradeUpdate,
      },
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: EpisodeGrade,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/episodes/:episodeId/recordings",
    alias: "getEpisodeRecordings",
    requestFormat: "json",
    parameters: [
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: EpisodeRecordingsResponse,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/episodes/:episodeId/stats",
    alias: "getEpisodeStats",
    requestFormat: "json",
    parameters: [
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: EpisodeStatsResponse,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/episodes/bulk",
    alias: "createEpisodesBulk",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for episodes to create in bulk`,
        type: "Body",
        schema: EpisodeBulkCreate,
      },
    ],
    response: z.array(Episode),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/episodes/export",
    alias: "exportEpisodes",
    description: `Exports episodes matching the given filters as a UTF-8 CSV file.
Maximum 30,000 rows. If exceeded, a 400 error is returned — apply filters to reduce the result set.
`,
    requestFormat: "json",
    parameters: [
      {
        name: "task_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "task_version_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "robot_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "user_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "status",
        type: "Query",
        schema: z.array(EpisodeCollectionStatus).optional(),
      },
      {
        name: "started_at_from",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "started_at_to",
        type: "Query",
        schema: z.string().optional(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/fleet/collection-trend",
    alias: "getFleetCollectionTrend",
    description: `Time series data collection trend by site and robot type`,
    requestFormat: "json",
    parameters: [
      {
        name: "granularity",
        type: "Query",
        schema: z.enum(["hourly", "daily", "monthly"]),
      },
      {
        name: "from",
        type: "Query",
        schema: z.string().datetime({ offset: true }),
      },
      {
        name: "to",
        type: "Query",
        schema: z.string().datetime({ offset: true }),
      },
    ],
    response: CollectionTrend,
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/fleet/stats",
    alias: "getFleetStats",
    description: `Period-based data collection time statistics by site and robot type`,
    requestFormat: "json",
    parameters: [
      {
        name: "from",
        type: "Query",
        schema: z.string().datetime({ offset: true }),
      },
      {
        name: "to",
        type: "Query",
        schema: z.string().datetime({ offset: true }),
      },
    ],
    response: z.array(FleetSiteStats),
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/fleet/summary",
    alias: "getFleetSummary",
    description: `Real-time robot status aggregation by site, robot type, and leader/follower`,
    requestFormat: "json",
    response: z.array(FleetSiteSummary),
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/locations",
    alias: "listLocations",
    requestFormat: "json",
    parameters: [
      {
        name: "page",
        type: "Query",
        schema: z.number().int().gte(1).optional().default(1),
      },
      {
        name: "limit",
        type: "Query",
        schema: z.number().int().gte(1).lte(1000).optional().default(20),
      },
      {
        name: "search",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "site_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "sort_by",
        type: "Query",
        schema: z.literal("name").optional(),
      },
      {
        name: "sort_order",
        type: "Query",
        schema: z.enum(["asc", "desc"]).optional(),
      },
    ],
    response: LocationListResponse,
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/locations",
    alias: "createLocation",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the location to create`,
        type: "Body",
        schema: LocationCreate,
      },
    ],
    response: Location,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/locations/:locationId",
    alias: "getLocationById",
    requestFormat: "json",
    parameters: [
      {
        name: "locationId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: Location,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/locations/:locationId",
    alias: "updateLocationById",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the location to update`,
        type: "Body",
        schema: z.object({ name: z.string() }).partial().passthrough(),
      },
      {
        name: "locationId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: Location,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "delete",
    path: "/locations/:locationId",
    alias: "deleteLocationById",
    requestFormat: "json",
    parameters: [
      {
        name: "locationId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/me",
    alias: "getMe",
    description: `Retrieve the currently authenticated user&#x27;s information`,
    requestFormat: "json",
    response: UserResponse,
    errors: [
      {
        status: 401,
        description: `Unauthorized`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/me",
    alias: "updateMe",
    description: `Update the currently authenticated user&#x27;s own profile.
Available to all authenticated roles. Only fields exposed by
&#x60;MeUpdateRequest&#x60; may be updated via this endpoint.
`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Self-updatable fields`,
        type: "Body",
        schema: z
          .object({ display_name: z.string().min(1).max(60) })
          .passthrough(),
      },
    ],
    response: UserResponse,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 401,
        description: `Unauthorized`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/organizations",
    alias: "listOrganizations",
    description: `Retrieve list of organizations`,
    requestFormat: "json",
    parameters: [
      {
        name: "limit",
        type: "Query",
        schema: z.number().int().gte(1).lte(1000).optional().default(50),
      },
      {
        name: "page",
        type: "Query",
        schema: z.number().int().gte(1).optional().default(1),
      },
    ],
    response: OrganizationListResponse,
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/organizations",
    alias: "createOrganization",
    description: `Create an organization`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the organization to create`,
        type: "Body",
        schema: OrganizationCreateRequest,
      },
    ],
    response: OrganizationResponse,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/organizations/:organization_id",
    alias: "getOrganizationById",
    description: `Retrieve organization details`,
    requestFormat: "json",
    parameters: [
      {
        name: "organization_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: OrganizationResponse,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/organizations/:organization_id",
    alias: "updateOrganizationById",
    description: `Update organization information`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the organization to update`,
        type: "Body",
        schema: OrganizationUpdateRequest,
      },
      {
        name: "organization_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: OrganizationResponse,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "delete",
    path: "/organizations/:organization_id",
    alias: "deleteOrganizationById",
    description: `Delete an organization`,
    requestFormat: "json",
    parameters: [
      {
        name: "organization_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/reports/operator-yield/export",
    alias: "exportOperatorYield",
    description: `Exports per-operator, per-task daily yield report as a UTF-8 CSV file.

One row per (date × operator × task). Columns include working window,
non-working / discarded / collected minutes, EP count and yield ratio.

NOTE: Cleansing-related filtering is not yet implemented. As a result,
&quot;discarded data&quot; only includes Cancel episodes (not cleansed-out data),
and &quot;collected data&quot; includes all Completed episodes (no quality filter).
Working time and non-working time are unaffected by cleansing.

Maximum 30,000 rows. If exceeded, a 400 error is returned — apply filters to reduce the result set.
`,
    requestFormat: "json",
    parameters: [
      {
        name: "date_from",
        type: "Query",
        schema: z.string(),
      },
      {
        name: "date_to",
        type: "Query",
        schema: z.string(),
      },
      {
        name: "location_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "task_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "user_id",
        type: "Query",
        schema: z.string().optional(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/robot-types",
    alias: "listRobotTypes",
    requestFormat: "json",
    parameters: [
      {
        name: "site_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "location_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "status",
        type: "Query",
        schema: z
          .union([
            z.literal(0),
            z.literal(1),
            z.literal(2),
            z.literal(3),
            z.literal(4),
            z.literal(5),
          ])
          .optional(),
      },
    ],
    response: z.array(z.string()),
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/robot/episodes",
    alias: "listRobotEpisodes",
    description: `Retrieve queued (ready) episodes assigned to the currently authenticated robot`,
    requestFormat: "json",
    response: z.array(Episode),
    errors: [
      {
        status: 401,
        description: `Unauthorized`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/robot/episodes/:episodeId",
    alias: "getRobotEpisodeById",
    description: `Retrieve episode details for robot access`,
    requestFormat: "json",
    parameters: [
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: Episode,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/robot/episodes/:episodeId/cancel",
    alias: "cancelRobotEpisode",
    description: `Cancel recording for the specified episode`,
    requestFormat: "json",
    parameters: [
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/robot/episodes/:episodeId/finish",
    alias: "finishRobotEpisode",
    description: `Finish recording for the specified episode`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: z
          .object({ occurred_at: z.string().datetime({ offset: true }) })
          .passthrough(),
      },
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/robot/episodes/:episodeId/start",
    alias: "startRobotEpisode",
    description: `Start recording for the specified episode`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: z
          .object({ occurred_at: z.string().datetime({ offset: true }) })
          .passthrough(),
      },
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/robot/episodes/:episodeId/subtasks/:subtaskId/complete",
    alias: "completeRobotSubTask",
    description: `Mark subtask as completed`,
    requestFormat: "json",
    parameters: [
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "subtaskId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/robot/episodes/:episodeId/subtasks/:subtaskId/executions",
    alias: "createRobotExecution",
    description: `Create a new execution for the specified subtask`,
    requestFormat: "json",
    parameters: [
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "subtaskId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.object({ execution_id: z.string() }).passthrough(),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/robot/episodes/:episodeId/subtasks/:subtaskId/executions/:executionId/cancel",
    alias: "cancelRobotExecution",
    description: `Cancel the specified execution`,
    requestFormat: "json",
    parameters: [
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "subtaskId",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "executionId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/robot/episodes/:episodeId/subtasks/:subtaskId/executions/:executionId/finish",
    alias: "finishRobotExecution",
    description: `Finish the specified execution`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: z
          .object({ occurred_at: z.string().datetime({ offset: true }) })
          .passthrough(),
      },
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "subtaskId",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "executionId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/robot/episodes/:episodeId/subtasks/:subtaskId/executions/:executionId/start",
    alias: "startRobotExecution",
    description: `Start the specified execution`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: z
          .object({ occurred_at: z.string().datetime({ offset: true }) })
          .passthrough(),
      },
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "subtaskId",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "executionId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/robot/episodes/:episodeId/subtasks/:subtaskId/skip",
    alias: "skipRobotSubTask",
    description: `Skip the specified subtask`,
    requestFormat: "json",
    parameters: [
      {
        name: "episodeId",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "subtaskId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/robot/episodes/repeat-last",
    alias: "repeatLastRobotEpisode",
    description: `Create a new episode for the same task as the robot&#x27;s most recently completed or cancelled episode`,
    requestFormat: "json",
    response: Episode,
    errors: [
      {
        status: 401,
        description: `Unauthorized`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/robot/me",
    alias: "getRobotMe",
    description: `Retrieve the profile of the currently authenticated robot`,
    requestFormat: "json",
    response: Robot,
    errors: [
      {
        status: 401,
        description: `Unauthorized`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/robot/status",
    alias: "updateRobotStatus",
    description: `Update the status and metrics for a robot`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Robot status update payload`,
        type: "Body",
        schema: RobotStatusUpdateRequest,
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 401,
        description: `Unauthorized`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/robots",
    alias: "listRobots",
    requestFormat: "json",
    parameters: [
      {
        name: "organization_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "site_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "location_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "status",
        type: "Query",
        schema: z
          .union([
            z.literal(0),
            z.literal(1),
            z.literal(2),
            z.literal(3),
            z.literal(4),
            z.literal(5),
          ])
          .optional(),
      },
      {
        name: "robot_type",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "page",
        type: "Query",
        schema: z.number().int().gte(1).optional().default(1),
      },
      {
        name: "limit",
        type: "Query",
        schema: z.number().int().gte(1).lte(1000).optional().default(20),
      },
      {
        name: "search",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "sort_by",
        type: "Query",
        schema: z
          .enum([
            "name",
            "location_id",
            "robot_type",
            "status",
            "leader_status",
            "last_heartbeat_at",
            "active_episode_id",
            "active_user_id",
          ])
          .optional(),
      },
      {
        name: "sort_order",
        type: "Query",
        schema: z.enum(["asc", "desc"]).optional(),
      },
    ],
    response: RobotListResponse,
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/robots",
    alias: "createRobot",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the robot to create`,
        type: "Body",
        schema: RobotCreate,
      },
    ],
    response: Robot,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/robots/:robotId",
    alias: "getRobotById",
    requestFormat: "json",
    parameters: [
      {
        name: "robotId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: Robot,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/robots/:robotId",
    alias: "updateRobotById",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the robot to update`,
        type: "Body",
        schema: RobotUpdate,
      },
      {
        name: "robotId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: Robot,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "delete",
    path: "/robots/:robotId",
    alias: "deleteRobotById",
    requestFormat: "json",
    parameters: [
      {
        name: "robotId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/robots/:robotId/operator",
    alias: "getRobotOperator",
    description: `Returns the current operator for this robot, or 204 if none.`,
    requestFormat: "json",
    parameters: [
      {
        name: "robotId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: RobotOperator,
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/robots/:robotId/operator",
    alias: "setRobotOperator",
    description: `Register or refresh the active operator for this robot. Stored in Redis with a short TTL. Returns 409 if another user already holds the operator lock.`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Operator identity`,
        type: "Body",
        schema: RobotOperatorRequest,
      },
      {
        name: "robotId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 409,
        description: `Robot locked by another operator`,
        schema: RobotOperator,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "delete",
    path: "/robots/:robotId/operator",
    alias: "clearRobotOperator",
    description: `Release the operator lock for this robot.`,
    requestFormat: "json",
    parameters: [
      {
        name: "robotId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 403,
        description: `Forbidden`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/sites",
    alias: "listSites",
    requestFormat: "json",
    parameters: [
      {
        name: "page",
        type: "Query",
        schema: z.number().int().gte(1).optional().default(1),
      },
      {
        name: "limit",
        type: "Query",
        schema: z.number().int().gte(1).lte(1000).optional().default(20),
      },
      {
        name: "search",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "organization_id",
        type: "Query",
        schema: z.string().optional(),
      },
    ],
    response: SiteListResponse,
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/sites",
    alias: "createSite",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the site to create`,
        type: "Body",
        schema: SiteCreate,
      },
    ],
    response: Site,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/sites/:siteId",
    alias: "getSiteById",
    requestFormat: "json",
    parameters: [
      {
        name: "siteId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: Site,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/sites/:siteId",
    alias: "updateSiteById",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the site to update`,
        type: "Body",
        schema: z.object({ name: z.string() }).partial().passthrough(),
      },
      {
        name: "siteId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: Site,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "delete",
    path: "/sites/:siteId",
    alias: "deleteSiteById",
    requestFormat: "json",
    parameters: [
      {
        name: "siteId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/sub-tasks",
    alias: "listSubTasks",
    requestFormat: "json",
    parameters: [
      {
        name: "limit",
        type: "Query",
        schema: z.number().int().gte(1).lte(100).optional().default(50),
      },
      {
        name: "page",
        type: "Query",
        schema: z.number().int().gte(1).optional().default(1),
      },
      {
        name: "task_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "task_version_id",
        type: "Query",
        schema: z.string().optional(),
      },
    ],
    response: z.array(SubTask),
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/sub-tasks",
    alias: "createSubTask",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the sub-task to create`,
        type: "Body",
        schema: SubTaskCreate,
      },
    ],
    response: SubTask,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/sub-tasks/:subtaskId",
    alias: "getSubTaskById",
    requestFormat: "json",
    parameters: [
      {
        name: "subtaskId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: SubTask,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/sub-tasks/:subtaskId",
    alias: "updateSubTaskById",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the sub-task to update`,
        type: "Body",
        schema: SubTaskUpdate,
      },
      {
        name: "subtaskId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: SubTask,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "delete",
    path: "/sub-tasks/:subtaskId",
    alias: "deleteSubTaskById",
    requestFormat: "json",
    parameters: [
      {
        name: "subtaskId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/sub-tasks/:subtaskId/complete",
    alias: "completeSubTask",
    requestFormat: "json",
    parameters: [
      {
        name: "subtaskId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: SubTask,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/sub-tasks/reorder",
    alias: "reorderSubTasks",
    description: `Reorder sub-tasks by providing an ordered list of sub-task IDs. The order_index of each sub-task will be updated to match the position in the array.`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Ordered list of sub-task IDs`,
        type: "Body",
        schema: SubTaskReorder,
      },
    ],
    response: z.array(SubTask),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 409,
        description: `Conflict`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/task-category-types",
    alias: "listTaskCategoryTypes",
    requestFormat: "json",
    response: z.array(TaskCategoryType),
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/task-tags",
    alias: "listTaskTags",
    requestFormat: "json",
    parameters: [
      {
        name: "category_type_id",
        type: "Query",
        schema: z.string().optional(),
      },
    ],
    response: z.array(TaskTag),
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/task-tags",
    alias: "createTaskTag",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: TaskTagCreate,
      },
    ],
    response: TaskTag,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/tasks",
    alias: "listTasks",
    requestFormat: "json",
    parameters: [
      {
        name: "has_approved_version",
        type: "Query",
        schema: z.boolean().optional(),
      },
      {
        name: "page",
        type: "Query",
        schema: z.number().int().gte(1).optional().default(1),
      },
      {
        name: "limit",
        type: "Query",
        schema: z.number().int().gte(1).lte(1000).optional().default(20),
      },
      {
        name: "sort_by",
        type: "Query",
        schema: z
          .enum([
            "name",
            "robot_type",
            "priority",
            "difficulty",
            "target_duration_seconds",
            "status",
            "recommended",
          ])
          .optional(),
      },
      {
        name: "sort_order",
        type: "Query",
        schema: z.enum(["asc", "desc"]).optional(),
      },
      {
        name: "status",
        type: "Query",
        schema: z.array(TaskStatus).optional(),
      },
      {
        name: "priority",
        type: "Query",
        schema: z.array(TaskPriority).optional(),
      },
      {
        name: "difficulty",
        type: "Query",
        schema: z.array(TaskDifficulty).optional(),
      },
      {
        name: "robot_type",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "search",
        type: "Query",
        schema: z.string().optional(),
      },
    ],
    response: TaskListResponse,
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/tasks",
    alias: "createTask",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the task to create`,
        type: "Body",
        schema: TaskCreate,
      },
    ],
    response: Task,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/tasks/:taskId",
    alias: "getTaskById",
    requestFormat: "json",
    parameters: [
      {
        name: "taskId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: Task,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/tasks/:taskId",
    alias: "updateTaskById",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the task to update`,
        type: "Body",
        schema: TaskUpdate,
      },
      {
        name: "taskId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: Task,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "delete",
    path: "/tasks/:taskId",
    alias: "deleteTaskById",
    requestFormat: "json",
    parameters: [
      {
        name: "taskId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/tasks/:taskId/versions",
    alias: "listTaskVersions",
    requestFormat: "json",
    parameters: [
      {
        name: "taskId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.array(TaskVersion),
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/tasks/:taskId/versions",
    alias: "createTaskVersion",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for creating a new task version`,
        type: "Body",
        schema: TaskVersionCreate,
      },
      {
        name: "taskId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: TaskVersion,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 409,
        description: `Conflict`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "patch",
    path: "/tasks/:taskId/versions/:versionId",
    alias: "updateTaskVersion",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Fields to update on the task version`,
        type: "Body",
        schema: TaskVersionUpdate,
      },
      {
        name: "taskId",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "versionId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: TaskVersion,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 409,
        description: `Conflict`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/tasks/:taskId/versions/:versionId/approve",
    alias: "approveTaskVersion",
    description: `Approve a draft task version, making it the active version for data collection`,
    requestFormat: "json",
    parameters: [
      {
        name: "taskId",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "versionId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: TaskVersion,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 409,
        description: `Conflict`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/tasks/:taskId/versions/:versionId/parameters",
    alias: "updateTaskVersionParameters",
    description: `Replace the parameter definitions on a draft task version. Fails if the version is already approved.`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `New parameter definitions`,
        type: "Body",
        schema: TaskVersionParametersUpdate,
      },
      {
        name: "taskId",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "versionId",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: TaskVersion,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 409,
        description: `Conflict`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/tasks/available-tags",
    alias: "getTaskAvailableTags",
    requestFormat: "json",
    parameters: [
      {
        name: "robot_type",
        type: "Query",
        schema: z.array(z.string()).optional(),
      },
      {
        name: "category_type_id",
        type: "Query",
        schema: z.string().optional(),
      },
    ],
    response: z.array(TaskTag),
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/tasks/completion-trend",
    alias: "getTaskCompletionTrend",
    requestFormat: "json",
    parameters: [
      {
        name: "robot_type",
        type: "Query",
        schema: z.array(z.string()).optional(),
      },
      {
        name: "category_type_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "tag_id",
        type: "Query",
        schema: z.array(z.string()).optional(),
      },
      {
        name: "group_by",
        type: "Query",
        schema: z.enum(["category", "status"]),
      },
      {
        name: "interval",
        type: "Query",
        schema: z.enum(["1week", "2week", "month"]).optional().default("2week"),
      },
      {
        name: "from",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "to",
        type: "Query",
        schema: z.string().optional(),
      },
    ],
    response: TaskCompletionTrend,
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/tasks/export",
    alias: "exportTasks",
    description: `Exports tasks matching the given filters as a UTF-8 CSV file.
Maximum 5,000 rows; apply filters to reduce the result set if exceeded.
`,
    requestFormat: "json",
    parameters: [
      {
        name: "status",
        type: "Query",
        schema: z.array(TaskStatus).optional(),
      },
      {
        name: "priority",
        type: "Query",
        schema: z.array(TaskPriority).optional(),
      },
      {
        name: "difficulty",
        type: "Query",
        schema: z.array(TaskDifficulty).optional(),
      },
      {
        name: "robot_type",
        type: "Query",
        schema: z.string().optional(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/tasks/import",
    alias: "importTasks",
    description: `Imports tasks from a CSV string. Duplicate rows (by name) are skipped.
Each imported task gets an initial v1.0.0 version created automatically.
`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: z.object({ csv_content: z.string() }).passthrough(),
      },
    ],
    response: TaskImportResponse,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/tasks/import/validate",
    alias: "validateTaskImport",
    description: `Validates a CSV string for task import. Returns a preview of valid rows,
duplicate rows (by name), and rows with errors.
`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: z.object({ csv_content: z.string() }).passthrough(),
      },
    ],
    response: TaskImportValidationResponse,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/tasks/summary",
    alias: "getTaskSummary",
    requestFormat: "json",
    parameters: [
      {
        name: "robot_type",
        type: "Query",
        schema: z.array(z.string()).optional(),
      },
      {
        name: "category_type_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "tag_id",
        type: "Query",
        schema: z.array(z.string()).optional(),
      },
      {
        name: "from",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "to",
        type: "Query",
        schema: z.string().optional(),
      },
    ],
    response: TaskSummaryResponse,
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/users",
    alias: "listUsers",
    description: `Retrieve list of users`,
    requestFormat: "json",
    parameters: [
      {
        name: "limit",
        type: "Query",
        schema: z.number().int().gte(1).lte(1000).optional().default(50),
      },
      {
        name: "page",
        type: "Query",
        schema: z.number().int().gte(1).optional().default(1),
      },
      {
        name: "organization_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "location_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "site_id",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "search",
        type: "Query",
        schema: z.string().optional(),
      },
      {
        name: "sort_by",
        type: "Query",
        schema: z
          .enum(["name", "email", "role", "location", "created_at"])
          .optional(),
      },
      {
        name: "sort_order",
        type: "Query",
        schema: z.enum(["asc", "desc"]).optional(),
      },
    ],
    response: UserListResponse,
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/users",
    alias: "createUser",
    description: `Create User`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the create user`,
        type: "Body",
        schema: UserCreateRequest,
      },
    ],
    response: UserResponse,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/users/:user_id",
    alias: "getUserById",
    description: `Retrieve user details`,
    requestFormat: "json",
    parameters: [
      {
        name: "user_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: UserResponse,
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/users/:user_id",
    alias: "updateUserById",
    description: `Update a user`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the user to update`,
        type: "Body",
        schema: UserUpdateRequest,
      },
      {
        name: "user_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: UserResponse,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "delete",
    path: "/users/:user_id",
    alias: "deleteUserById",
    requestFormat: "json",
    parameters: [
      {
        name: "user_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/users/:user_id/locations",
    alias: "updateUserLocations",
    description: `Assign locations to a user (admin only)`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: UserLocationsUpdateRequest,
      },
      {
        name: "user_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: UserResponse,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 403,
        description: `Forbidden`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "get",
    path: "/users/:user_id/permissions",
    alias: "listUserPermissions",
    requestFormat: "json",
    parameters: [
      {
        name: "user_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.array(Permission),
    errors: [
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/users/:user_id/permissions",
    alias: "grantUserPermission",
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data of the permission to grant`,
        type: "Body",
        schema: z.object({ code: z.string() }).passthrough(),
      },
      {
        name: "user_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "delete",
    path: "/users/:user_id/permissions",
    alias: "revokeUserPermission",
    requestFormat: "json",
    parameters: [
      {
        name: "user_id",
        type: "Path",
        schema: z.string(),
      },
      {
        name: "code",
        type: "Query",
        schema: z.string(),
      },
    ],
    response: z.void(),
    errors: [
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/users/:user_id/role",
    alias: "updateUserRole",
    description: `Update a user&#x27;s role`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        description: `Data for the user role to update`,
        type: "Body",
        schema: UserRoleUpdateRequest,
      },
      {
        name: "user_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: UserResponse,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 403,
        description: `Forbidden`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "put",
    path: "/users/:user_id/sites",
    alias: "updateUserSites",
    description: `Assign sites to a user (admin only)`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: UserSitesUpdateRequest,
      },
      {
        name: "user_id",
        type: "Path",
        schema: z.string(),
      },
    ],
    response: UserResponse,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 403,
        description: `Forbidden`,
        schema: ErrorResponse,
      },
      {
        status: 404,
        description: `Not found`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/users/import",
    alias: "importUsers",
    description: `Imports users from a CSV string. Valid rows are processed independently;
duplicate rows are skipped. Each valid row creates a user in the database.
`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: z.object({ csv_content: z.string() }).passthrough(),
      },
    ],
    response: UserImportResponse,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
  {
    method: "post",
    path: "/users/import/validate",
    alias: "validateUserImport",
    description: `Validates a CSV string for user import. Returns a preview of valid rows,
duplicate rows (by email), and rows with errors. Does not create users.
`,
    requestFormat: "json",
    parameters: [
      {
        name: "body",
        type: "Body",
        schema: z.object({ csv_content: z.string() }).passthrough(),
      },
    ],
    response: UserImportValidationResponse,
    errors: [
      {
        status: 400,
        description: `Bad request`,
        schema: ErrorResponse,
      },
      {
        status: 500,
        description: `Internal server error`,
        schema: ErrorResponse,
      },
    ],
  },
]);

export const api = new Zodios(endpoints);

export function createApiClient(baseUrl: string, options?: ZodiosOptions) {
  return new Zodios(baseUrl, endpoints, options);
}
