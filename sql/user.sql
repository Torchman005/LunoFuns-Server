CREATE TABLE `user` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID，自增主键',
    `username` VARCHAR(32) NOT NULL COMMENT '用户名，登录账号',
    `password` VARCHAR(128) NOT NULL COMMENT '密码，bcrypt加密存储',
    `nickname` VARCHAR(32) NOT NULL COMMENT '昵称，展示名称',
    `avatar` VARCHAR(512) DEFAULT '' COMMENT '头像URL，存储对象存储地址',
    `email` VARCHAR(128) NOT NULL COMMENT '邮箱，用于找回密码',
    `phone` VARCHAR(20) DEFAULT '' COMMENT '手机号',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-正常',
    `role` TINYINT NOT NULL DEFAULT 0 COMMENT '角色：0-普通用户，1-VIP，2-管理员',
    `last_login_at` DATETIME DEFAULT NULL COMMENT '最后登录时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_email` (`email`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';