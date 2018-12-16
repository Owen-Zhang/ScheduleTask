CREATE TABLE `worker` (
  `id`          int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name`        varchar(50)      NOT NULL DEFAULT '' COMMENT '机器名称',
  `key`         varchar(200)     NOT NULL DEFAULT '' COMMENT 'worker标识符',
  `note`        varchar(500)     NOT NULL DEFAULT '' COMMENT '说明',
  `status`      int(11)          NOT NULL DEFAULT '0'COMMENT 'worker的状态',

  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `worker_log` (
  `id`          int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name`        varchar(50)      NOT NULL DEFAULT '' COMMENT '机器名称',
  `key`         varchar(200)     NOT NULL DEFAULT '' COMMENT 'worker标识符',
  `ip`          varchar(30)      NOT NULL DEFAULT '' COMMENT 'worker ip地址',
  `port`        int(11)          NOT NULL DEFAULT '0' COMMENT 'worker的端口',
  `osname`      varchar(50)      NOT NULL DEFAULT '' COMMENT 'worker系统名(windows,linux)',
  `note`        varchar(500)     NOT NULL DEFAULT '' COMMENT '说明',
  `status`      int(11)          NOT NULL DEFAULT '0' COMMENT '状态, 0正常，1失败(服务中心在标识)',

  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `task` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `group_id` int(11) NOT NULL DEFAULT '0' COMMENT '分组ID',
  `system` varchar(50) NOT NULL DEFAULT '' COMMENT '操作系统名称(windows,linux)',
  `task_name` varchar(50) NOT NULL DEFAULT '' COMMENT '任务名称',
  `task_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '任务类型',
  `description` varchar(200) NOT NULL DEFAULT '' COMMENT '任务描述',
  `cron_spec` varchar(100) NOT NULL DEFAULT '' COMMENT '时间表达式',
  `run_file_folder` varchar(200) NOT NULL DEFAULT '' COMMENT '运行程序的文件夹信息',
  `old_zip_file` varchar(200) NOT NULL DEFAULT '' COMMENT '当前运行程序的上传zip包名称(用户上传的名称)',
  `concurrent` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否只允许一个实例',
  `task_api_url` varchar(200) NOT NULL DEFAULT '' COMMENT 'API调用的url: http://www.baidu.com/abc',
  `task_api_method` varchar(10) NOT NULL DEFAULT '' COMMENT 'API调用的Method:post,get',
  `api_header` varchar(500) NOT NULL DEFAULT '' COMMENT 'POST提交时带的body内容',
  `api_body` varchar(8000) NOT NULL DEFAULT '' COMMENT 'POST提交时带的body内容',
  `command` text NOT NULL COMMENT '命令详情',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0停用 1启用',
  `notify` tinyint(4) NOT NULL DEFAULT '0' COMMENT '通知设置',
  `notify_email` text NOT NULL COMMENT '通知人列表',
  `time_out` smallint(6) NOT NULL DEFAULT '0' COMMENT '超时设置',
  `execute_times` int(11) NOT NULL DEFAULT '0' COMMENT '累计执行次数',
  `prev_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '上次执行时间',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `version` int(11) unsigned NOT NULL DEFAULT 0 COMMENT '版本号',
  `zip_file_path` VARCHAR(300) NOT NULL DEFAULT '' COMMENT 'zip获取地址(我们上传到文件服务器生成的文件名)',
  `deleted` int(11) unsigned NOT NULL DEFAULT  0 COMMENT '是否删除,1表示删除,0表示正常',
  `worker_key` varchar(100) NOT NULL DEFAULT '' COMMENT 'worker key(标识此任务由谁在运行)',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `task_group` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `group_name` varchar(50) NOT NULL DEFAULT '' COMMENT '组名',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '说明',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `task_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '任务ID',
  `output` mediumtext NOT NULL COMMENT '任务输出',
  `error` text NOT NULL COMMENT '错误信息',
  `status` tinyint(4) NOT NULL COMMENT '状态',
  `process_time` int(11) NOT NULL DEFAULT '0' COMMENT '消耗时间/毫秒',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户名',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
  `password` char(32) NOT NULL DEFAULT '' COMMENT '密码',
  `salt` char(10) NOT NULL DEFAULT '' COMMENT '密码盐',
  `last_login` int(11) NOT NULL DEFAULT '0' COMMENT '最后登录时间',
  `last_ip` char(15) NOT NULL DEFAULT '' COMMENT '最后登录IP',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0正常 -1禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_name` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `user` (`id`, `user_name`, `email`, `password`, `salt`, `last_login`, `last_ip`, `status`)
VALUES (1,'admin','admin@example.com','e10adc3949ba59abbe56e057f20f883e','',0,'',0);