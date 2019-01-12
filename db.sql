
CREATE TABLE IF NOT EXISTS filters
(
    uid serial NOT NULL,
    guildid character varying(100) NOT NULL,
    phrase character varying(500) NOT NULL,
    CONSTRAINT filter_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE);

CREATE TABLE IF NOT EXISTS localxp
(
    uid serial NOT NULL,
    guildid character varying(100) NOT NULL,
    userid character varying(500) NOT NULL,
    xp int NOT NULL,
    CONSTRAINT localxp_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE);

CREATE TABLE IF NOT EXISTS globalxp
(
    uid serial NOT NULL,
    userid character varying(500) NOT NULL,
    xp int NOT NULL,
    CONSTRAINT globalxp_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE);

CREATE TABLE IF NOT EXISTS xpignoredchannels
(
    uid serial NOT NULL,
    channelid character varying(500) NOT NULL,
    CONSTRAINT xpignoredchannel_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE);

CREATE TABLE IF NOT EXISTS discordusers
(
    uid serial NOT NULL,
    userid character varying(100) NOT NULL,
    username character varying(500) NOT NULL,
    discriminator character varying(500) NOT NULL,
    xp int NOT NULL,
    nextxpgaintime timestamp NOT NULL,
    xpexcluded boolean NOT NULL,
    reputation int NOT NULL,
    cangivereptime timestamp NOT NULL,
    CONSTRAINT discorduser_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE);

CREATE TABLE IF NOT EXISTS commandlog
(
    uid serial NOT NULL,
    command character varying(100) NOT NULL,
    args character varying(2000) NOT NULL,
    userid character varying(30) NOT NULL,
    guildid character varying(30)NOT NULL,
    channelid character varying(30) NOT NULL,
    messageid character varying(30) NOT NULL,
    tstamp timestamp NOT NULL,
    CONSTRAINT commandlog_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE); 

CREATE TABLE IF NOT EXISTS discordguilds
(
    uid serial NOT NULL,
    guildid character varying(30) NOT NULL,
    usestrikes boolean NOT NULL,
    maxstrikes int NOT NULL,
    CONSTRAINT discordguild_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE); 

CREATE TABLE IF NOT EXISTS filterignorechannels
(
    uid serial NOT NULL,
    guildid character varying(30) NOT NULL,
    channelid character varying(30) NOT NULL,
    CONSTRAINT filterignorechannel_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE);

CREATE TABLE IF NOT EXISTS userroles
(
    uid serial NOT NULL,
    guildid character varying(30) NOT NULL,
    userid character varying(30) NOT NULL,
    roleid character varying(30) NOT NULL,
    CONSTRAINT userrole_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE);

CREATE TABLE IF NOT EXISTS strikes
(
    uid serial NOT NULL,
    guildid character varying(30) NOT NULL,
    userid character varying(30) NOT NULL,
    reason character varying(250) NOT NULL,
    executorid character varying(30) NOT NULL,
    tstamp timestamp NOT NULL,
    CONSTRAINT strike_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE); 
