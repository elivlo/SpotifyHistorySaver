CREATE TABLE `tracks` (
  `id` varchar(255) PRIMARY KEY,
  `name` varchar(255),
  `track_number` int,
  `disc_number` int,
  `explicit` boolean
);

CREATE TABLE `history_entries` (
    `id` int PRIMARY KEY AUTO_INCREMENT,
    `track_id` varchar(255),
    `played_at` datetime
);

CREATE TABLE `artists_tracks` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `artist_id` varchar(255),
  `track_id` varchar(255)
);

CREATE TABLE `artists` (
   `id` varchar(255) PRIMARY KEY,
   `name` varchar(255)
);

ALTER TABLE `history_entries` ADD FOREIGN KEY (`track_id`) REFERENCES `tracks` (`id`);

ALTER TABLE `artists_tracks` ADD FOREIGN KEY (`artist_id`) REFERENCES `artists` (`id`);

ALTER TABLE `artists_tracks` ADD FOREIGN KEY (`track_id`) REFERENCES `tracks` (`id`);
