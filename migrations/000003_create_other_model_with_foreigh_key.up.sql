CREATE TABLE IF NOT EXISTS `gateway_service_info` (
  `id` bigint(20) UNSIGNED AUTO_INCREMENT NOT NULL COMMENT '自增主键',
  `load_type` tinyint(4) NOT NULL COMMENT '负载类型 0=http 1=tcp 2=grpc',
  `service_name` varchar(255) NOT NULL COMMENT '服务名称 6-128 数字字母下划线',
  `service_desc` varchar(255) NOT NULL COMMENT '服务描述',
  `create_at` datetime NOT NULL  COMMENT '添加时间',
  `update_at` datetime NOT NULL  COMMENT '更新时间',
  `is_delete` tinyint(4) COMMENT '是否删除 1=删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='网关基本信息表';

CREATE TABLE IF NOT EXISTS `gateway_service_access_control` (
  `id` bigint(20) UNSIGNED AUTO_INCREMENT NOT NULL COMMENT '自增主键',
  `service_id` bigint(20) UNSIGNED NOT NULL COMMENT '服务id',
  `open_auth` tinyint(4) NOT NULL COMMENT '是否开启权限 1=开启',
  `black_list` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '黑名单ip',
  `white_list` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '白名单ip',
  `white_host_name` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '白名单主机',
  `clientip_flow_limit` int(11) NOT NULL COMMENT '客户端ip限流',
  `service_flow_limit` int(20) NOT NULL COMMENT '服务端限流',
  `service_flow_type` int(20) NOT NULL COMMENT '服务端限流类型',
  `client_flow_type` int(20) NOT NULL COMMENT '客户端限流类型',
  PRIMARY KEY (`id`),
  KEY `fk_access_control_service` (`service_id`),
  CONSTRAINT `fk_access_control_service` FOREIGN KEY (`service_id`) REFERENCES `gateway_service_info` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='网关权限控制表';


CREATE TABLE IF NOT EXISTS `gateway_service_grpc_rule` (
  `id` bigint(20) UNSIGNED AUTO_INCREMENT NOT NULL COMMENT '自增主键',
  `service_id` bigint(20) UNSIGNED NOT NULL  COMMENT '服务id',
  `port` int(5) NOT NULL  COMMENT '端口',
  `header_transfor` varchar(5000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL  COMMENT 'header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔',
  PRIMARY KEY (`id`),
  KEY `fk_service_grpc_rule_service` (service_id),
  CONSTRAINT `fk_service_grpc_rule_service` FOREIGN KEY (`service_id`) REFERENCES `gateway_service_info` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='网关路由匹配表';

CREATE TABLE IF NOT EXISTS `gateway_service_http_rule` (
  `id` bigint(20) UNSIGNED AUTO_INCREMENT NOT NULL COMMENT '自增主键',
  `service_id` bigint(20) UNSIGNED NOT NULL COMMENT '服务id',
  `rule_type` tinyint(4) NOT NULL  COMMENT '匹配类型 0=url前缀url_prefix 1=域名domain ',
  `rule` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'type=domain表示域名，type=url_prefix时表示url前缀',
  `need_https` tinyint(4) NOT NULL  COMMENT '支持https 1=支持',
  `need_strip_uri` tinyint(4) NOT NULL  COMMENT '启用strip_uri 1=启用',
  `need_websocket` tinyint(4) NOT NULL  COMMENT '是否支持websocket 1=支持',
  `url_rewrite` varchar(5000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔',
  `header_transfor` varchar(5000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔',
  PRIMARY KEY (`id`),
  KEY `fk_service_http_rule_service` (`service_id`),
  CONSTRAINT `fk_service_http_rule_service` FOREIGN KEY (`service_id`) REFERENCES `gateway_service_info` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='网关路由匹配表';

CREATE TABLE IF NOT EXISTS `gateway_service_load_balance` (
  `id` bigint(20) UNSIGNED AUTO_INCREMENT NOT NULL COMMENT '自增主键',
  `service_id` bigint(20) UNSIGNED NOT NULL  COMMENT '服务id',
  `check_method` tinyint(20) NOT NULL  COMMENT '检查方法 0=tcpchk,检测端口是否握手成功',
  `check_timeout` int(10) NOT NULL  COMMENT 'check超时时间,单位s',
  `check_interval` int(11) NOT NULL  COMMENT '检查间隔, 单位s',
  `round_type` tinyint(4) NOT NULL COMMENT '轮询方式 0=random 1=round-robin 2=weight_round-robin 3=ip_hash',
  `ip_list` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'ip列表',
  `weight_list` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权重列表',
  `forbid_list` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '禁用ip列表',
  `upstream_connect_timeout` int(11) NOT NULL  COMMENT '建立连接超时, 单位s',
  `upstream_header_timeout` int(11) NOT NULL  COMMENT '获取header超时, 单位s',
  `upstream_idle_timeout` int(10) NOT NULL  COMMENT '链接最大空闲时间, 单位s',
  `upstream_max_idle` int(11) NOT NULL  COMMENT '最大空闲链接数',
  PRIMARY KEY (`id`),
  KEY `fk_service_load_balance_service` (`service_id`),
  CONSTRAINT `fk_service_load_balance_service` FOREIGN KEY (`service_id`) REFERENCES `gateway_service_info` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='网关负载表';


CREATE TABLE IF NOT EXISTS `gateway_service_tcp_rule` (
  `id` bigint(20) UNSIGNED AUTO_INCREMENT NOT NULL COMMENT '自增主键',
  `service_id` bigint(20) UNSIGNED NOT NULL COMMENT '服务id',
  `port` int(5) NOT NULL  COMMENT '端口号',
  PRIMARY KEY (`id`),
  KEY `fk_service_tcp_rule_service` (`service_id`),
  CONSTRAINT `fk_service_tcp_rule_service` FOREIGN KEY (`service_id`) REFERENCES `gateway_service_info` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='网关路由匹配表';