-- +goose Up
-- +goose StatementBegin

create type user_status as ENUM (
  'PENDING',
  'ACTIVE',
  'SUSPENDED',
  'DELETED'
);

create table users(
  id text primary key default gen_random_uuid(),
  email text unique not null,
  username text unique not null,
  password_hash text not null,  -- if oauth signup just fill this with secure random hash
  stream_key text unique,
  status user_status default 'ACTIVE',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table user_profiles(
  user_id text references users(id) on delete cascade,
  display_name text not null,
  bio text,
  avatar_url text,

  -- Denormalized counters and flags
  follower_count int default 0,
  is_live boolean default false,
  updated_at timestamptz default now(),
  primary key (user_id)
);

create table follows(
  follower_id text references users(id) on delete cascade, -- Guy who clicked follow
  following_id text references users(id) on delete cascade, -- Guy being followed
  created_at timestamptz default now(),
  primary key (follower_id, following_id),
  check (follower_id != following_id)
);

create table roles (
  id serial primary key,
  name text not null unique
);

create table permissions(
  id serial primary key,
  slug text unique not null
);

create table role_permissions(
  role_id int references roles(id) on delete cascade,
  permission_id int references permissions(id) on delete cascade,
  primary key (role_id, permission_id)
);

create table user_roles(
  user_id text references users(id) on delete cascade,
  role_id int references roles(id) on delete cascade,
  primary key (user_id, role_id)
);

create table channel_staff(
  channel_id text references users(id) on delete cascade, -- The streamer
  user_id text references users(id) on delete cascade, -- The channel mod
  role_id int references roles(id) on delete cascade,
  created_at timestamptz default now(),

  primary key (channel_id, user_id)

);


-- Stream sessions 
create table streams(
  id text primary key default gen_random_uuid(),
  user_id text references users(id) on delete cascade,
  title text not null,
  created_at timestamptz default now(),
  ended_at timestamptz,

  -- snapshot metrics
  peak_viewers int default 0,
  thumbnail_url text
);

create type video_source_type as ENUM ('UPLOAD', 'STREAM_ARCHIVE', 'CLIP');

create table videos(
  id text primary key default gen_random_uuid(),
  user_id text not null references users(id) on delete cascade,
  stream_id text references streams(id) on delete set null, -- optional, vod of a stream or maybe uploaded vid
  title text not null,
  url text not null, 
  description text,
  duration_seconds int not null,
  source_type video_source_type not null,
  created_at timestamptz default now()

);

-- Index for fast search and "Who is live" lists
CREATE INDEX idx_profiles_display_name ON user_profiles(display_name);
CREATE INDEX idx_profiles_is_live ON user_profiles(is_live) WHERE is_live = TRUE;

-- Foreign key indexes for better join performance
CREATE INDEX idx_follows_follower ON follows(follower_id);
CREATE INDEX idx_follows_following ON follows(following_id);
CREATE INDEX idx_streams_user_id ON streams(user_id);
CREATE INDEX idx_streams_created_at ON streams(created_at DESC);
CREATE INDEX idx_videos_user_id ON videos(user_id);
CREATE INDEX idx_videos_stream_id ON videos(stream_id);
CREATE INDEX idx_videos_created_at ON videos(created_at DESC);
CREATE INDEX idx_channel_staff_channel_id ON channel_staff(channel_id);
CREATE INDEX idx_channel_staff_user_id ON channel_staff(user_id);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists videos;
drop table if exists streams;
drop table if exists follows;
drop table if exists channel_staff;
drop table if exists user_roles;
drop table if exists role_permissions;
drop table if exists permissions;
drop table if exists roles;
drop table if exists user_profiles;
drop table if exists users;
drop type if exists video_source_type;
drop type if exists user_status;
-- +goose StatementEnd
