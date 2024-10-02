CREATE TABLE videos (
	videoId TEXT(32) NOT NULL,
	title TEXT(256),
	description TEXT(1024),
	CONSTRAINT videos_pk PRIMARY KEY (videoId)
);

