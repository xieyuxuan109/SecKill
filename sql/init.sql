CREATE DATABASE IF NOT EXISTS seckill;
USE seckill;

CREATE TABLE `seckill_activity` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `stock` int NOT NULL,
  `start_at` bigint NOT NULL,
  `end_at` bigint NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `seckill_order` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `order_no` varchar(64) NOT NULL,
  `user_id` bigint NOT NULL,
  `act_id` bigint NOT NULL,
  `status` tinyint NOT NULL DEFAULT 1,
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_act` (`user_id`,`act_id`)
);