
show create table project;

DROP TABLE IF EXISTS `co_chapter`;
CREATE TABLE `co_chapter`  (
  `id` char(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '章节id',
  `name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '章节名称',
  `pid` char(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '父章节id    首级章节设置默认值',
  `course_id` char(20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '课程id',
  `chapter_order` smallint(5) UNSIGNED NOT NULL COMMENT '章节顺序    小数靠前',
  `has_coursetime` tinyint(3) UNSIGNED NOT NULL COMMENT '是否存在课时',
  `version` int(10) UNSIGNED NOT NULL COMMENT '乐观锁',≈
  `del_flag` tinyint(3) UNSIGNED NOT NULL COMMENT '逻辑删除     1-删除',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = COMPACT;

CREATE TABLE `project` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `update_time` datetime(6) NOT NULL COMMENT '更新时间',
  `create_time` datetime(6) NOT NULL COMMENT '创建时间',
  `title` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '项目标题',
  `description` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '项目描述',
  `create_by_id` int NOT NULL COMMENT '项目创建人',
  `delete_time` datetime(6) DEFAULT NULL COMMENT '删除时间',
  `del_flag` tinyint(1) NOT NULL COMMENT '逻辑删除标记',
  PRIMARY KEY (`id`),
  KEY `project_create_by_id_32bf9f46_fk_auth_user_id` (`create_by_id`),
  CONSTRAINT `project_create_by_id_32bf9f46_fk_auth_user_id` FOREIGN KEY (`create_by_id`) REFERENCES `auth_user` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=350 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

CREATE TABLE `auth_user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `last_login` datetime(6) DEFAULT NULL,
  `is_superuser` tinyint(1) NOT NULL,
  `username` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `first_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `last_name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(254) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `is_staff` tinyint(1) NOT NULL,
  `is_active` tinyint(1) NOT NULL,
  `date_joined` datetime(6) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=205 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci


CREATE TABLE `dataset_file_in_version` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `del_flag` tinyint(1) NOT NULL COMMENT '逻辑删除标记',
  `delete_time` datetime(6) DEFAULT NULL COMMENT '删除时间',
  `update_time` datetime(6) NOT NULL COMMENT '更新时间',
  `create_time` datetime(6) NOT NULL COMMENT '创建时间',
  `data_id` bigint NOT NULL COMMENT '版本数据 id',
  `data_version_id` bigint NOT NULL COMMENT '数据集版本 id',
  `dataset_id` bigint NOT NULL COMMENT '数据集 id',
  PRIMARY KEY (`id`),
  KEY `app_dataset_fileofda_data_id_b2ce2372_fk_app_datas` (`data_id`),
  KEY `app_dataset_fileofda_data_version_id_01732413_fk_app_datas` (`data_version_id`),
  KEY `app_dataset_fileofda_dataset_id_c81e285c_fk_app_datas` (`dataset_id`),
  CONSTRAINT `app_dataset_fileofda_data_id_b2ce2372_fk_app_datas` FOREIGN KEY (`data_id`) REFERENCES `dataset_data` (`id`),
  CONSTRAINT `app_dataset_fileofda_data_version_id_01732413_fk_app_datas` FOREIGN KEY (`data_version_id`) REFERENCES `dataset_data_version` (`id`),
  CONSTRAINT `app_dataset_fileofda_dataset_id_c81e285c_fk_app_datas` FOREIGN KEY (`dataset_id`) REFERENCES `dataset` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=190731 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci




CREATE TABLE `auth_user_code` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `code` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '验证码',
  `email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户邮箱',
  `create_time` datetime(6) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_code_unique_user_email` (`username`,`email`),
  UNIQUE KEY `user_code_unique_user` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

CREATE TABLE `annotation_genomes_database` (
  `ftp_path` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件url',
  `taxId` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类学id',
  `organism` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '物种',
  `assembly` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '基因组唯一编码',
  PRIMARY KEY (`assembly`),
  KEY `idx_genomes_database__tax_id` (`taxId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci