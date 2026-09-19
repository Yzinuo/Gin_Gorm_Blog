CREATE TABLE IF NOT EXISTS `managed_asset` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `key` varchar(120) NOT NULL,
  `kind` varchar(32) NOT NULL,
  `active_version_id` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_managed_asset_key` (`key`),
  KEY `idx_managed_asset_active_version_id` (`active_version_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `managed_asset_version` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `asset_id` bigint NOT NULL,
  `version` varchar(64) NOT NULL,
  `status` varchar(24) NOT NULL,
  `source_manifest` longtext NOT NULL,
  `delivery_manifest` longtext NOT NULL,
  `validation_report` longtext NOT NULL,
  `created_by` bigint DEFAULT 0,
  `published_at` datetime(3) DEFAULT NULL,
  `failure_reason` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_asset_version` (`asset_id`, `version`),
  KEY `idx_managed_asset_version_asset_id` (`asset_id`),
  KEY `idx_managed_asset_version_status` (`status`),
  KEY `idx_managed_asset_version_created_by` (`created_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Assign these resources to the existing administrator role in the Admin UI.
-- Super administrators can use them immediately; non-super users fail closed
-- until the role_resource mappings are explicitly granted.
INSERT IGNORE INTO `resource` (`created_at`, `updated_at`, `name`, `parent_id`, `url`, `method`, `anonymous`) VALUES
  (NOW(), NOW(), '受管资源版本列表', 0, '/asset/:key/versions', 'GET', 0),
  (NOW(), NOW(), '上传受管资源版本', 0, '/asset/:key/versions', 'POST', 0),
  (NOW(), NOW(), '查看受管资源版本', 0, '/asset/:key/versions/:id', 'GET', 0),
  (NOW(), NOW(), '发布或回滚受管资源', 0, '/asset/:key/versions/:id/publish', 'POST', 0),
  (NOW(), NOW(), '删除未发布资源版本', 0, '/asset/:key/versions/:id', 'DELETE', 0);
