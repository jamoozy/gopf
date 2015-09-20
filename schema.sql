
drop table if exists playlists;
create table playlists (
  id integer primary key autoincrement,
  name text not null
);

drop table if exists files;
create table files (
  id integer primary key autoincrement,
  path text not null
);

drop table if exists tags;
create table tags (
  id integer primary key autoincrement,
  name text not null
);

drop table if exists playlist_files;
create table playlist_files (
  playlist_id references playlists(id),
  file_id references files(id)
);

drop table if exists file_tags;
create table file_tags (
  file_id references files(id),
  tag_id references tags(id)
);
