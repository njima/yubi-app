package repository

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type EpisodeSortBy string

const (
	EpisodeSortByTask       EpisodeSortBy = "task"
	EpisodeSortByRobot      EpisodeSortBy = "robot"
	EpisodeSortByRecordedBy EpisodeSortBy = "recorded_by"
	EpisodeSortByStartedAt  EpisodeSortBy = "started_at"
	EpisodeSortByEndedAt    EpisodeSortBy = "ended_at"
	EpisodeSortByError      EpisodeSortBy = "error"
)

type LocationSortBy string

const (
	LocationSortByName LocationSortBy = "name"
)

type RobotSortBy string

const (
	RobotSortByActiveEpisodeID RobotSortBy = "active_episode_id"
	RobotSortByActiveUserID    RobotSortBy = "active_user_id"
	RobotSortByLastHeartbeatAt RobotSortBy = "last_heartbeat_at"
	RobotSortByLeaderStatus    RobotSortBy = "leader_status"
	RobotSortByLocationID      RobotSortBy = "location_id"
	RobotSortByName            RobotSortBy = "name"
	RobotSortByRobotType       RobotSortBy = "robot_type"
	RobotSortByStatus          RobotSortBy = "status"
)

type TaskSortBy string

const (
	TaskSortByDifficulty            TaskSortBy = "difficulty"
	TaskSortByName                  TaskSortBy = "name"
	TaskSortByPriority              TaskSortBy = "priority"
	TaskSortByRecommended           TaskSortBy = "recommended"
	TaskSortByRobotType             TaskSortBy = "robot_type"
	TaskSortByStatus                TaskSortBy = "status"
	TaskSortByTargetDurationSeconds TaskSortBy = "target_duration_seconds"
)

type UserSortBy string

const (
	UserSortByCreatedAt UserSortBy = "created_at"
	UserSortByEmail     UserSortBy = "email"
	UserSortByLocation  UserSortBy = "location"
	UserSortByName      UserSortBy = "name"
	UserSortByRole      UserSortBy = "role"
)
