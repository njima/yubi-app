-- Initial seed data for all tables
-- Run this script to populate the database with sample data

-- 1. Organization (must be first due to foreign key dependencies)
INSERT INTO "organization" (created_at, updated_at, id_natural, name, description)
VALUES (
    NOW(),
    NOW(),
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    'Sample Organization',
    'This is a sample organization for development and testing purposes.'
) ON CONFLICT (id_natural) DO NOTHING;

-- 2. User (Default Admin)
INSERT INTO "user" (created_at, updated_at, id_natural, google_sub, name, email)
VALUES (
    NOW(),
    NOW(),
    '69fad3df-d73f-45e1-9fb4-df52bd4857b0',
    'dev-google-sub-default-admin',
    'Default Admin',
    'admin@example.com'
) ON CONFLICT (id_natural) DO NOTHING;

-- 3. Organization membership (Default Admin)
INSERT INTO "organization_membership" (created_at, updated_at, id_natural, user_id, organization_id, role)
VALUES (
    NOW(),
    NOW(),
    '0d1b16c5-38df-4277-8fd0-90f9d0a05da1',
    '69fad3df-d73f-45e1-9fb4-df52bd4857b0',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    0
) ON CONFLICT (id_natural) DO NOTHING;

-- 4. Site
INSERT INTO "site" (created_at, updated_at, id_natural, organization_id, name)
VALUES (
    NOW(),
    NOW(),
    'c115b89a-55d9-4407-8f08-a25617beea2c',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    'Sample Site'
) ON CONFLICT (id_natural) DO NOTHING;

-- 5. Location
INSERT INTO "location" (created_at, updated_at, id_natural, organization_id, site_id, name)
VALUES (
    NOW(),
    NOW(),
    '91154897-df4b-4b39-8c4c-b48daf4a3b37',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    'c115b89a-55d9-4407-8f08-a25617beea2c',
    'Sample Location'
) ON CONFLICT (id_natural) DO NOTHING;

-- 6. Robot
INSERT INTO "robot" (created_at, updated_at, id_natural, organization_id, location_id, name, robot_type, status, robot_config)
VALUES (
    NOW(),
    NOW(),
    'c2f8e62b-ea23-4a50-8660-d707e4d5c2bc',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    '91154897-df4b-4b39-8c4c-b48daf4a3b37',
    'Sample Robot',
    'yubi-stationary',
    5,
    '{"host": "localhost", "port": 9090, "cameras": [{"namespace": "camera_0", "name": "Front Camera"}]}'::jsonb
) ON CONFLICT (id_natural) DO NOTHING;

-- 7. Task
INSERT INTO "task" (created_at, updated_at, id_natural, organization_id, name, description, priority, difficulty, status, deadline, manual_url)
VALUES (
    NOW(),
    NOW(),
    '6013935a-ab9c-4bd8-b59d-49958f516d47',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    'Sample Task',
    'This is a sample task for demonstration purposes.',
    1,
    1,
    0,
    NOW() + INTERVAL '30 days',
    'https://example.com/manual'
) ON CONFLICT (id_natural) DO NOTHING;

-- 8. Task Version
INSERT INTO "task_version" (created_at, updated_at, id_natural, organization_id, task_id, version, schema_hash, is_active, approval_status)
VALUES (
    NOW(),
    NOW(),
    '5437b101-6d9d-495f-a4e8-45420eb10d99',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    '6013935a-ab9c-4bd8-b59d-49958f516d47',
    'v1.0.0',
    'sha256:abc123def456',
    true,
    1
) ON CONFLICT (id_natural) DO NOTHING;

-- 9. SubTask
INSERT INTO "subtask" (created_at, updated_at, id_natural, organization_id, task_version_id, order_index, name, description)
VALUES (
    NOW(),
    NOW(),
    '7065d47f-8de7-4b6d-af34-4aa924dfa98e',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    '5437b101-6d9d-495f-a4e8-45420eb10d99',
    0,
    'Sample SubTask 1',
    'This is the first sample subtask.'
) ON CONFLICT (id_natural) DO NOTHING;

INSERT INTO "subtask" (created_at, updated_at, id_natural, organization_id, task_version_id, order_index, name, description)
VALUES (
    NOW(),
    NOW(),
    '1d8744a2-cfe0-4cfd-88ea-f38fbc50c640',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    '5437b101-6d9d-495f-a4e8-45420eb10d99',
    1,
    'Sample SubTask 2',
    'This is the second sample subtask.'
) ON CONFLICT (id_natural) DO NOTHING;

INSERT INTO "subtask" (created_at, updated_at, id_natural, organization_id, task_version_id, order_index, name, description)
VALUES (
    NOW(),
    NOW(),
    '33f37ea4-b36b-40ff-8c59-568e2f5cac57',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    '5437b101-6d9d-495f-a4e8-45420eb10d99',
    2,
    'Sample SubTask 3',
    'This is the third sample subtask.'
) ON CONFLICT (id_natural) DO NOTHING;

-- 10. Episode (all collection_status patterns: 0=Ready, 1=Recording, 2=Cancel, 3=Completed)
-- Episode with collection_status: Ready (0)
INSERT INTO "episode" (created_at, updated_at, id_natural, organization_id, task_version_id, location_id, robot_id, user_id, collection_status)
VALUES (
    NOW(),
    NOW(),
    '358637c3-bf03-4408-845f-f9b189ec767f',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    '5437b101-6d9d-495f-a4e8-45420eb10d99',
    '91154897-df4b-4b39-8c4c-b48daf4a3b37',
    'c2f8e62b-ea23-4a50-8660-d707e4d5c2bc',
    '69fad3df-d73f-45e1-9fb4-df52bd4857b0',
    0
) ON CONFLICT (id_natural) DO NOTHING;

-- Episode with collection_status: Completed (3)
INSERT INTO "episode" (created_at, updated_at, id_natural, organization_id, task_version_id, location_id, robot_id, user_id, collection_status, started_at, finished_at)
VALUES (
    NOW(),
    NOW(),
    'c74cd916-b377-4b2a-8084-195701eca190',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    '5437b101-6d9d-495f-a4e8-45420eb10d99',
    '91154897-df4b-4b39-8c4c-b48daf4a3b37',
    'c2f8e62b-ea23-4a50-8660-d707e4d5c2bc',
    '69fad3df-d73f-45e1-9fb4-df52bd4857b0',
    3,
    NOW() - INTERVAL '2 hours',
    NOW() - INTERVAL '1 hour'
) ON CONFLICT (id_natural) DO NOTHING;

-- Episode with collection_status: Cancel (2)
INSERT INTO "episode" (created_at, updated_at, id_natural, organization_id, task_version_id, location_id, robot_id, user_id, collection_status, started_at, finished_at)
VALUES (
    NOW(),
    NOW(),
    '26bea46c-b781-4e0b-8109-05ba5a72a519',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    '5437b101-6d9d-495f-a4e8-45420eb10d99',
    '91154897-df4b-4b39-8c4c-b48daf4a3b37',
    'c2f8e62b-ea23-4a50-8660-d707e4d5c2bc',
    '69fad3df-d73f-45e1-9fb4-df52bd4857b0',
    2,
    NOW() - INTERVAL '1 hour',
    NOW()
) ON CONFLICT (id_natural) DO NOTHING;

-- Episode with collection_status: Completed (3)
INSERT INTO "episode" (created_at, updated_at, id_natural, organization_id, task_version_id, location_id, robot_id, user_id, collection_status, started_at, finished_at)
VALUES (
    NOW(),
    NOW(),
    '2cc38dbf-fa77-41d9-bc22-47389d413331',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    '5437b101-6d9d-495f-a4e8-45420eb10d99',
    '91154897-df4b-4b39-8c4c-b48daf4a3b37',
    'c2f8e62b-ea23-4a50-8660-d707e4d5c2bc',
    '69fad3df-d73f-45e1-9fb4-df52bd4857b0',
    3,
    NOW() - INTERVAL '2 hours',
    NOW() - INTERVAL '1 hour'
) ON CONFLICT (id_natural) DO NOTHING;

-- 11. Episode SubTask (collection_status: 0=Ready, 1=InProgress, 2=Completed, 3=Skipped, 4=Cancelled)
-- sub_task1: Ready (Ready episode)
INSERT INTO "episode_sub_task" (created_at, updated_at, id_natural, organization_id, episode_id, sub_task_id, collection_status)
VALUES (
    NOW(),
    NOW(),
    'a1c352f7-54c6-4d74-ba6d-aa1352dfa001',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    '358637c3-bf03-4408-845f-f9b189ec767f',
    '7065d47f-8de7-4b6d-af34-4aa924dfa98e',
    0
) ON CONFLICT (id_natural) DO NOTHING;

-- sub_task2: Ready (Ready episode)
INSERT INTO "episode_sub_task" (created_at, updated_at, id_natural, organization_id, episode_id, sub_task_id, collection_status)
VALUES (
    NOW(),
    NOW(),
    'a1c352f7-54c6-4d74-ba6d-aa1352dfa002',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    '358637c3-bf03-4408-845f-f9b189ec767f',
    '1d8744a2-cfe0-4cfd-88ea-f38fbc50c640',
    0
) ON CONFLICT (id_natural) DO NOTHING;

-- sub_task3: Ready (Ready episode)
INSERT INTO "episode_sub_task" (created_at, updated_at, id_natural, organization_id, episode_id, sub_task_id, collection_status)
VALUES (
    NOW(),
    NOW(),
    'a1c352f7-54c6-4d74-ba6d-aa1352dfa003',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    '358637c3-bf03-4408-845f-f9b189ec767f',
    '33f37ea4-b36b-40ff-8c59-568e2f5cac57',
    0
) ON CONFLICT (id_natural) DO NOTHING;

-- sub_task1: Completed (Completed episode)
INSERT INTO "episode_sub_task" (created_at, updated_at, id_natural, organization_id, episode_id, sub_task_id, collection_status)
VALUES (
    NOW(),
    NOW(),
    'b1c352f7-54c6-4d74-ba6d-aa1352dfaee0',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    'c74cd916-b377-4b2a-8084-195701eca190',
    '7065d47f-8de7-4b6d-af34-4aa924dfa98e',
    2
) ON CONFLICT (id_natural) DO NOTHING;

-- sub_task2: Completed (Completed episode)
INSERT INTO "episode_sub_task" (created_at, updated_at, id_natural, organization_id, episode_id, sub_task_id, collection_status)
VALUES (
    NOW(),
    NOW(),
    '1cfb5edb-ce00-4c6c-8c58-69cdbd06ad5e',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    'c74cd916-b377-4b2a-8084-195701eca190',
    '1d8744a2-cfe0-4cfd-88ea-f38fbc50c640',
    2
) ON CONFLICT (id_natural) DO NOTHING;

-- sub_task3: Completed (Completed episode)
INSERT INTO "episode_sub_task" (created_at, updated_at, id_natural, organization_id, episode_id, sub_task_id, collection_status)
VALUES (
    NOW(),
    NOW(),
    '5cebcf54-d97b-4bee-90f9-373f1246f824',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    'c74cd916-b377-4b2a-8084-195701eca190',
    '33f37ea4-b36b-40ff-8c59-568e2f5cac57',
    2
) ON CONFLICT (id_natural) DO NOTHING;

-- 11. Episode SubTask Execution (execution_status: 0=Ready, 1=Started, 2=Cancelled, 3=Finished)
-- Execution for sub_task1: Finished
INSERT INTO "episode_sub_task_execution" (created_at, updated_at, id_natural, organization_id, episode_sub_task_id, execution_status, started_at, finished_at)
VALUES (
    NOW(),
    NOW(),
    '5abd5a33-60ea-4b40-a7df-750fb1511536',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    'b1c352f7-54c6-4d74-ba6d-aa1352dfaee0',
    3,
    NOW() - INTERVAL '30 minutes',
    NOW() - INTERVAL '20 minutes'
) ON CONFLICT (id_natural) DO NOTHING;

-- Execution for sub_task1: Finished (second execution)
INSERT INTO "episode_sub_task_execution" (created_at, updated_at, id_natural, organization_id, episode_sub_task_id, execution_status, started_at, finished_at)
VALUES (
    NOW(),
    NOW(),
    '8e15ba4c-8912-4a80-b9eb-a657e272f380',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    'b1c352f7-54c6-4d74-ba6d-aa1352dfaee0',
    3,
    NOW() - INTERVAL '15 minutes',
    NOW() - INTERVAL '5 minutes'
) ON CONFLICT (id_natural) DO NOTHING;

-- Execution for sub_task2: Finished
INSERT INTO "episode_sub_task_execution" (created_at, updated_at, id_natural, organization_id, episode_sub_task_id, execution_status, started_at, finished_at)
VALUES (
    NOW(),
    NOW(),
    '7ab30835-b535-4dad-b66f-8919008d65af',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    '1cfb5edb-ce00-4c6c-8c58-69cdbd06ad5e',
    3,
    NOW() - INTERVAL '10 minutes',
    NOW() - INTERVAL '3 minutes'
) ON CONFLICT (id_natural) DO NOTHING;

-- Execution for sub_task3: Finished
INSERT INTO "episode_sub_task_execution" (created_at, updated_at, id_natural, organization_id, episode_sub_task_id, execution_status, started_at, finished_at)
VALUES (
    NOW(),
    NOW(),
    'e4c21f8a-9b34-4d67-a812-3f5c6e7d8a90',
    '7bfbe942-5fd6-4525-ac13-0356147c202b',
    '5cebcf54-d97b-4bee-90f9-373f1246f824',
    3,
    NOW() - INTERVAL '3 minutes',
    NOW() - INTERVAL '1 minute'
) ON CONFLICT (id_natural) DO NOTHING;

-- 12. Task Category Types
INSERT INTO "task_category_type" ("id", "slug", "name") VALUES
    ('3e7f61e7-bb28-44ff-8633-1738867868a8', 'basic-skill', 'Basic Skill'),
    ('2293d4a1-1a23-4928-8847-b4db2b4d72c6', 'application',  'Application')
ON CONFLICT (id) DO NOTHING;

-- 13. Task Tags
-- Basic Skill
INSERT INTO "task_tag" ("id", "name", "category_type_id") VALUES
    ('640599cc-a9b4-475f-8374-9cbf85089815', 'Pick & Place',                       '3e7f61e7-bb28-44ff-8633-1738867868a8'),
    ('3bf43fd2-d754-4054-8383-2635af357307', 'Reorient & Regrasp',                 '3e7f61e7-bb28-44ff-8633-1738867868a8'),
    ('06275be1-947c-49b6-b7dd-5820f3b0aca2', 'Non-Prehensile (Push/Pull/Slide)',   '3e7f61e7-bb28-44ff-8633-1738867868a8'),
    ('02ff6897-72ce-4e6b-b7d9-4095b6f17af1', 'Clutter / Clearing',                '3e7f61e7-bb28-44ff-8633-1738867868a8'),
    ('6fef84d8-0865-4d26-84b9-a0a33673d309', 'Articulated (Open/Close/Turn)',      '3e7f61e7-bb28-44ff-8633-1738867868a8'),
    ('0e7e5ae0-7a6b-4694-8d5c-0feff0e5ebb4', 'High-Precision (Insertion/Assembly)','3e7f61e7-bb28-44ff-8633-1738867868a8'),
    ('e9129649-fbea-4cc5-be01-847c7d6c998b', 'Tool Use',                          '3e7f61e7-bb28-44ff-8633-1738867868a8'),
    ('b8df828c-f38f-426a-9d6f-76f89c4f495a', 'Deformable (Cloth/Cable/Bag)',       '3e7f61e7-bb28-44ff-8633-1738867868a8'),
    ('0a8306c7-13e7-49a9-8bf6-a92cc863b687', 'Bimanual Coordination',             '3e7f61e7-bb28-44ff-8633-1738867868a8'),
    ('865e57ca-2d6c-475b-9f23-cd3a08b9dfe8', 'Mobile Manipulation',               '3e7f61e7-bb28-44ff-8633-1738867868a8'),
    ('3e4fedd4-b9df-4084-a6b3-1009b4bc0f51', 'Recovery (Failure Correction)',      '3e7f61e7-bb28-44ff-8633-1738867868a8'),
    ('4ddc35d5-6388-4038-b553-dcbff963d080', 'Human-Robot Interaction (Handover)', '3e7f61e7-bb28-44ff-8633-1738867868a8'),
    ('1c156ebb-3935-478c-8186-a561c73c922a', 'Other',                             '3e7f61e7-bb28-44ff-8633-1738867868a8')
ON CONFLICT (id) DO NOTHING;

-- Application
INSERT INTO "task_tag" ("id", "name", "category_type_id") VALUES
    ('d404c7e8-a094-47db-97c3-bb36b50f0dd4', 'Home',         '2293d4a1-1a23-4928-8847-b4db2b4d72c6'),
    ('3358f4c4-9144-42fb-ad66-8bc192145409', 'Retail',       '2293d4a1-1a23-4928-8847-b4db2b4d72c6'),
    ('78dbc482-a193-4897-b326-cc0bf924b14f', 'Manufacturing', '2293d4a1-1a23-4928-8847-b4db2b4d72c6'),
    ('896f997f-efa5-414a-8857-07a4ce689ac7', 'Logistics',    '2293d4a1-1a23-4928-8847-b4db2b4d72c6'),
    ('a3aec8c4-0cd5-4ebf-b7be-284908a54e3b', 'Construction', '2293d4a1-1a23-4928-8847-b4db2b4d72c6'),
    ('d479e487-056b-479c-80fe-9a64fed48b1b', 'Office',       '2293d4a1-1a23-4928-8847-b4db2b4d72c6')
ON CONFLICT (id) DO NOTHING;
