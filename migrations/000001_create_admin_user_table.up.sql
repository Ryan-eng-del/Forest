CREATE TABLE IF NOT EXISTS `gateway_admin` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT "用户名",
  `salt` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT "盐值",
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT "密码",
  `create_at` datetime NOT NULL COMMENT "新增时间",
  `update_at` datetime NOT NULL COMMENT "更新时间",
  `is_delete` tinyint(1) NOT NULL COMMENT "是否删除" DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;