import { z } from "zod";

import { schemas } from "../generated/api";

export type Task = z.infer<typeof schemas.Task>;
export type TaskListResponse = z.infer<typeof schemas.TaskListResponse>;
export type TaskCreate = z.infer<typeof schemas.TaskCreate>;
export type TaskUpdate = z.infer<typeof schemas.TaskUpdate>;
export type TaskTag = z.infer<typeof schemas.TaskTag>;
export type TaskTagCreate = z.infer<typeof schemas.TaskTagCreate>;
export type TaskCategoryType = z.infer<typeof schemas.TaskCategoryType>;
export type Robot = z.infer<typeof schemas.Robot>;
export type RobotListResponse = z.infer<typeof schemas.RobotListResponse>;
export type RobotCreate = z.infer<typeof schemas.RobotCreate>;
export type RobotUpdate = z.infer<typeof schemas.RobotUpdate>;
export type Episode = z.infer<typeof schemas.Episode>;
export type EpisodeListResponse = z.infer<typeof schemas.EpisodeListResponse>;
export type EpisodeCreate = z.infer<typeof schemas.EpisodeCreate>;
export type EpisodeUpdate = z.infer<typeof schemas.EpisodeUpdate>;
export type EpisodeBulkCreateRequest = EpisodeCreate & { count: number };
export type EpisodeGrade = z.infer<typeof schemas.EpisodeGrade>;
export type EpisodeGradeUpdate = z.infer<typeof schemas.EpisodeGradeUpdate>;
export type EpisodeGradeListResponse = z.infer<
  typeof schemas.EpisodeGradeListResponse
>;
export type Site = z.infer<typeof schemas.Site>;
export type SiteListResponse = z.infer<typeof schemas.SiteListResponse>;
export type Location = z.infer<typeof schemas.Location>;
export type LocationListResponse = z.infer<typeof schemas.LocationListResponse>;
export type LocationCreate = z.infer<typeof schemas.LocationCreate>;
export type LocationUpdate = z.infer<typeof schemas.LocationUpdate>;
export type UserResponse = z.infer<typeof schemas.UserResponse>;
export type UserListResponse = z.infer<typeof schemas.UserListResponse>;
export type UserRoleUpdateRequest = z.infer<
  typeof schemas.UserRoleUpdateRequest
>;
export type UserUpdateRequest = z.infer<typeof schemas.UserUpdateRequest>;
export type MeUpdateRequest = z.infer<typeof schemas.MeUpdateRequest>;

export type OrganizationResponse = z.infer<typeof schemas.OrganizationResponse>;
export type OrganizationListResponse = z.infer<
  typeof schemas.OrganizationListResponse
>;
export type UserCreateRequest = z.infer<typeof schemas.UserCreateRequest>;
export type SubTask = z.infer<typeof schemas.SubTask>;
export type SubTaskCreate = z.infer<typeof schemas.SubTaskCreate>;
export type SubTaskUpdate = z.infer<typeof schemas.SubTaskUpdate>;
export type SubTaskReorder = z.infer<typeof schemas.SubTaskReorder>;
export type EpisodeRecordingsResponse = z.infer<
  typeof schemas.EpisodeRecordingsResponse
>;
export type EpisodeStatsResponse = z.infer<typeof schemas.EpisodeStatsResponse>;
export type FleetSiteSummary = z.infer<typeof schemas.FleetSiteSummary>;
export type FleetSiteStats = z.infer<typeof schemas.FleetSiteStats>;
export type CollectionTrend = z.infer<typeof schemas.CollectionTrend>;
export type ApiKeyResponse = z.infer<typeof schemas.ApiKeyResponse>;
export type ApiKeyListResponse = z.infer<typeof schemas.ApiKeyListResponse>;
export type ApiKeyCreateRequest = z.infer<typeof schemas.ApiKeyCreateRequest>;
export type ApiKeyCreateResponse = z.infer<typeof schemas.ApiKeyCreateResponse>;
