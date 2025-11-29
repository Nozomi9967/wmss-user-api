/*
 Navicat Premium Data Transfer

 Source Server         : mysql
 Source Server Type    : MySQL
 Source Server Version : 80402 (8.4.2)
 Source Host           : localhost:3306
 Source Schema         : wmss

 Target Server Type    : MySQL
 Target Server Version : 80402 (8.4.2)
 File Encoding         : 65001

 Date: 29/11/2025 23:36:21
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for customer_bank_card
-- ----------------------------
DROP TABLE IF EXISTS `customer_bank_card`;
CREATE TABLE `customer_bank_card`  (
  `card_id` bigint NOT NULL AUTO_INCREMENT COMMENT '银行卡记录唯一标识，自增',
  `customer_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '客户标识，关联客户信息表',
  `bank_card_number` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '银行卡号，加密存储',
  `bank_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '开户行名称，如“中国工商银行”“中信银行”',
  `card_balance` decimal(18, 2) NOT NULL COMMENT '银行卡余额，虚拟充值金额需标注',
  `is_virtual` tinyint(1) NOT NULL COMMENT '是否虚拟银行卡，1-是，0-否',
  `bind_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '绑定状态，如“正常”“已解绑”',
  `bind_time` datetime NOT NULL COMMENT '绑定时间',
  `unbind_time` datetime NULL DEFAULT NULL COMMENT '解绑时间，解绑时记录',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`card_id`) USING BTREE,
  UNIQUE INDEX `uk_customer_card`(`customer_id` ASC, `bank_card_number` ASC) USING BTREE,
  INDEX `idx_bind_status`(`bind_status` ASC) USING BTREE,
  CONSTRAINT `fk_cbc_customer_id` FOREIGN KEY (`customer_id`) REFERENCES `customer_info` (`customer_id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '客户银行卡表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of customer_bank_card
-- ----------------------------
INSERT INTO `customer_bank_card` VALUES (1, '039d346f35c54359bcead8c23be2b345', '面包车', '蔡文韬', 384.15, 0, '已解绑', '2025-11-29 21:35:47', '2025-11-29 23:34:24', '2025-11-29 21:35:46', '2025-11-29 23:34:24', '2025-11-29 23:34:24');

-- ----------------------------
-- Table structure for customer_behavior
-- ----------------------------
DROP TABLE IF EXISTS `customer_behavior`;
CREATE TABLE `customer_behavior`  (
  `behavior_id` bigint NOT NULL AUTO_INCREMENT COMMENT '行为记录唯一标识，自增',
  `customer_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '客户标识，关联客户信息表',
  `behavior_type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '行为类型，如“产品查询”“申购”“赎回”“风险测评”',
  `behavior_time` datetime NOT NULL COMMENT '行为发生时间',
  `related_product_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '关联产品ID，如查询、申购、赎回的产品',
  `behavior_detail` json NULL COMMENT '行为详情，如查询的时间范围、申购的金额',
  `ip_address` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '行为发生时的IP地址',
  `device_info` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '设备信息，如“PC端-Windows-Chrome”“移动端-Android-微信小程序”',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  PRIMARY KEY (`behavior_id`) USING BTREE,
  INDEX `idx_customer_id`(`customer_id` ASC) USING BTREE,
  INDEX `idx_behavior_type`(`behavior_type` ASC) USING BTREE,
  INDEX `idx_behavior_time`(`behavior_time` ASC) USING BTREE,
  INDEX `idx_related_product`(`related_product_id` ASC) USING BTREE,
  CONSTRAINT `fk_cb_customer_id` FOREIGN KEY (`customer_id`) REFERENCES `customer_info` (`customer_id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `fk_cb_product_id` FOREIGN KEY (`related_product_id`) REFERENCES `product_info` (`product_id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '客户行为分析表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of customer_behavior
-- ----------------------------

-- ----------------------------
-- Table structure for customer_info
-- ----------------------------
DROP TABLE IF EXISTS `customer_info`;
CREATE TABLE `customer_info`  (
  `customer_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '客户唯一标识，客户编号',
  `customer_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '客户姓名（个人）/企业名称（企业）',
  `customer_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '客户类型，如“个人”“企业”',
  `id_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '证件类型，如“身份证”“营业执照”',
  `id_number` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '证件号码，加密存储',
  `risk_level` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '客户风险等级，如“R1”“R2”“R3”“R4”“R5”',
  `risk_evaluation_time` datetime NOT NULL COMMENT '风险测评时间',
  `risk_evaluation_expire_time` datetime NOT NULL COMMENT '风险测评过期时间',
  `contact_phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '联系电话，加密存储',
  `email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '电子邮箱',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '客户创建时间（开户时间）',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '客户信息更新时间',
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`customer_id`) USING BTREE,
  INDEX `idx_customer_type`(`customer_type` ASC) USING BTREE,
  INDEX `idx_risk_level`(`risk_level` ASC) USING BTREE,
  INDEX `idx_id_number`(`id_number` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '客户信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of customer_info
-- ----------------------------
INSERT INTO `customer_info` VALUES ('039d346f35c54359bcead8c23be2b345', '丛三锋', 'min', '75', '39', 'ea', '2026-07-22 16:10:34', '2025-11-30 00:17:37', '95880835066', 'm97k5q_k7562@tom.com', '2025-11-28 13:21:04', '2025-11-29 20:13:36', '2025-11-29 20:13:36');

-- ----------------------------
-- Table structure for product_info
-- ----------------------------
DROP TABLE IF EXISTS `product_info`;
CREATE TABLE `product_info`  (
  `product_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `product_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `product_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `product_sub_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `risk_level` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `product_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `manager` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `custodian` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `purchase_fee_rate` decimal(10, 4) NULL DEFAULT NULL,
  `redemption_fee_rule` json NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `create_by` bigint UNSIGNED NOT NULL,
  PRIMARY KEY (`product_id`) USING BTREE,
  UNIQUE INDEX `idx_product_info_product_name`(`product_name` ASC) USING BTREE,
  INDEX `idx_product_info_product_type`(`product_type` ASC) USING BTREE,
  INDEX `idx_product_info_risk_level`(`risk_level` ASC) USING BTREE,
  INDEX `idx_product_info_product_status`(`product_status` ASC) USING BTREE,
  INDEX `idx_product_info_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of product_info
-- ----------------------------
INSERT INTO `product_info` VALUES ('1a1ec462120a4056bb35241731b8be37', '人体工程学的冷冻自行车', 'labore', 'enim pariatur', 'minim elit', 'non magna ', 'in aliqua in', 'ut amet occaecat in', 3.0000, '{}', '养将过活元价三采样没。量发点。红段究转力市至论专行。干他连却万斯土白价。', '2025-11-06 22:39:59', '2025-11-09 10:38:23', NULL, 1);
INSERT INTO `product_info` VALUES ('3b8d3f9713b548f89d06597925d4742c', '精1', 'anim irure', 'cupid', 'v', 'velit adipi', 'ut cillum', 'nostrud occaecat sint ', 13.0000, '{}', '铁上论手改门至后气采。维油个采样准。公且治权论想。', '2025-11-06 22:45:15', '2025-11-06 22:45:15', NULL, 1);
INSERT INTO `product_info` VALUES ('3bd20921ba2a434f919e4fad5fd90183', '22131', 'anim irure', 'cupid', 'v', 'velit adipi', 'ut cillum', 'nostrud occaecat sint ', 13.0000, '{}', '铁上论手改门至后气采。维油个采样准。公且治权论想。', '2025-11-07 22:46:11', '2025-11-07 22:46:11', NULL, 1);
INSERT INTO `product_info` VALUES ('52aa010242804ba8a3f67e760ab1ec28', '精', 'anim irure', 'cupid', 'v', 'velit adipi', 'ut cillum', 'nostrud occaecat sint ', 13.0000, '{}', '铁上论手改门至后气采。维油个采样准。公且治权论想。', '2025-11-06 22:44:05', '2025-11-06 22:44:05', NULL, 1);
INSERT INTO `product_info` VALUES ('7ee5e00858534acb99d76874e3ad6104', '1', 'anim irure', 'cupid', 'v', 'velit adipi', 'ut cillum', 'nostrud occaecat sint ', 13.0000, '{}', '铁上论手改门至后气采。维油个采样准。公且治权论想。', '2025-11-06 22:45:21', '2025-11-06 22:45:21', NULL, 1);
INSERT INTO `product_info` VALUES ('b51085c6ff5749baae9e2d7c53790ead', '精致的软自行车', 'anim irure', 'cupid', 'v', 'velit adipi', 'ut cillum', 'nostrud occaecat sint ', 13.0000, '{}', '铁上论手改门至后气采。维油个采样准。公且治权论想。', '2025-11-06 22:43:43', '2025-11-09 10:37:43', '2025-11-09 10:37:43.000', 1);

-- ----------------------------
-- Table structure for sys_permission
-- ----------------------------
DROP TABLE IF EXISTS `sys_permission`;
CREATE TABLE `sys_permission`  (
  `permission_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '权限唯一标识',
  `permission_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '权限名称，如“产品管理-新增”“开户操作”“清算处理”',
  `permission_code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '权限编码，如“product:add”“customer:openAccount”“liquidation:process”',
  `permission_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '权限类型，如“菜单权限”“按钮权限”',
  `parent_permission_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '父权限ID，用于构建权限层级',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '权限创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '权限信息更新时间',
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`permission_id`) USING BTREE,
  UNIQUE INDEX `uk_permission_code`(`permission_code` ASC) USING BTREE,
  INDEX `idx_parent_permission`(`parent_permission_id` ASC) USING BTREE,
  INDEX `idx_permission_type`(`permission_type` ASC) USING BTREE,
  CONSTRAINT `fk_sp_parent_permission` FOREIGN KEY (`parent_permission_id`) REFERENCES `sys_permission` (`permission_id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '系统权限表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_permission
-- ----------------------------
INSERT INTO `sys_permission` VALUES ('28', '方颖', '18', 'sint veniam', NULL, '2025-11-22 21:20:15', '2025-11-22 21:20:36', '2025-11-22 21:20:36');
INSERT INTO `sys_permission` VALUES ('47', '字治涛', '47', 'aliqua ut', NULL, '2025-11-14 13:36:15', '2025-11-21 23:52:40', '2025-11-21 22:23:35');
INSERT INTO `sys_permission` VALUES ('P000', '产品管理', 'menu:product', '菜单权限', 'P000', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P001', '产品管理-新增', 'product:add', '按钮权限', 'P000', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P0010', '客户管理', 'menu:customer', '菜单权限', 'P0010', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P002', '产品管理-编辑', 'product:edit', '按钮权限', 'P000', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P0020', '交易管理', 'menu:trade', '菜单权限', 'P0020', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P003', '产品管理-查看', 'product:view', '按钮权限', 'P000', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P0030', '清算管理', 'menu:liquidation', '菜单权限', 'P0030', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P004', '客户开户', 'customer:openAccount', '按钮权限', 'P0010', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P0040', '系统管理', 'menu:system', '菜单权限', 'P0040', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P005', '银行卡绑定', 'customer:bindCard', '按钮权限', 'P0010', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P006', '客户风险测评', 'customer:riskEvaluate', '按钮权限', 'P0010', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P007', '申购申请受理', 'trade:purchase', '按钮权限', 'P0020', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P008', '赎回申请受理', 'trade:redemption', '按钮权限', 'P0020', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P009', '交易确认处理', 'trade:confirm', '按钮权限', 'P0020', '2024-01-01 09:30:00', '2024-01-02 15:20:00', NULL);
INSERT INTO `sys_permission` VALUES ('P010', '清算处理', 'liquidation:process', '按钮权限', 'P0030', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P011', '清算日志查看', 'liquidation:viewLog', '按钮权限', 'P0030', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P012', '角色管理', 'sys:roleManage', '按钮权限', 'P0040', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P013', '用户管理', 'sys:userManage', '按钮权限', 'P0040', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);
INSERT INTO `sys_permission` VALUES ('P014', '权限配置', 'sys:permissionConfig', '按钮权限', 'P0040', '2024-01-01 09:30:00', '2024-01-01 09:30:00', NULL);

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role`  (
  `role_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '角色唯一标识',
  `role_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '角色名称，如“柜台操作员”“理财运营人员”“系统管理员”',
  `role_desc` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '角色描述，说明角色职责和权限范围',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '角色创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '角色信息更新时间',
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`role_id`) USING BTREE,
  UNIQUE INDEX `uk_role_name`(`role_name` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '系统角色表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO `sys_role` VALUES ('b9d19ed38e884a72b897a65ac0257c12', '泉敬阳', 'consec', '2025-11-22 22:42:42', '2025-11-22 22:42:42', NULL);
INSERT INTO `sys_role` VALUES ('R001', '系统管理员', '负责系统配置、角色权限管理、用户账号维护', '2024-01-01 09:00:00', '2024-01-01 09:00:00', NULL);
INSERT INTO `sys_role` VALUES ('R002', '理财运营人员', '负责产品上下架、净值更新、交易清算处理', '2024-01-01 09:00:00', '2024-01-02 14:30:00', NULL);
INSERT INTO `sys_role` VALUES ('R003', '柜台操作员', '负责客户开户、银行卡绑定、申购赎回申请受理', '2024-01-01 09:00:00', '2024-01-01 09:00:00', NULL);
INSERT INTO `sys_role` VALUES ('R004', '数据分析师', '负责客户行为数据统计、业务报表生成', '2024-01-03 10:15:00', '2024-01-03 10:15:00', NULL);

-- ----------------------------
-- Table structure for sys_role_permission
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_permission`;
CREATE TABLE `sys_role_permission`  (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '关联记录唯一标识，自增',
  `role_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '角色ID，关联角色表',
  `permission_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '权限ID，关联权限表',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '关联记录创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_role_permission`(`role_id` ASC, `permission_id` ASC) USING BTREE,
  INDEX `idx_permission_id`(`permission_id` ASC) USING BTREE,
  CONSTRAINT `fk_srp_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `sys_permission` (`permission_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_srp_role_id` FOREIGN KEY (`role_id`) REFERENCES `sys_role` (`role_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 46 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '角色权限关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role_permission
-- ----------------------------
INSERT INTO `sys_role_permission` VALUES (2, 'R001', 'P000', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (3, 'R001', 'P001', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (4, 'R001', 'P002', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (5, 'R001', 'P003', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (6, 'R001', 'P0010', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (7, 'R001', 'P004', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (8, 'R001', 'P005', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (9, 'R001', 'P006', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (10, 'R001', 'P0020', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (11, 'R001', 'P007', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (12, 'R001', 'P008', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (13, 'R001', 'P009', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (14, 'R001', 'P0030', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (15, 'R001', 'P010', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (16, 'R001', 'P011', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (17, 'R001', 'P0040', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (18, 'R001', 'P012', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (19, 'R001', 'P013', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (20, 'R001', 'P014', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (21, 'R002', 'P000', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (22, 'R002', 'P001', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (23, 'R002', 'P002', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (24, 'R002', 'P003', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (25, 'R002', 'P0020', '2024-01-02 15:20:00');
INSERT INTO `sys_role_permission` VALUES (26, 'R002', 'P009', '2024-01-02 15:20:00');
INSERT INTO `sys_role_permission` VALUES (27, 'R002', 'P0030', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (28, 'R002', 'P010', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (29, 'R002', 'P011', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (30, 'R003', 'P0010', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (31, 'R003', 'P004', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (32, 'R003', 'P005', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (33, 'R003', 'P006', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (34, 'R003', 'P000', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (35, 'R003', 'P003', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (36, 'R003', 'P0020', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (37, 'R003', 'P007', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (38, 'R003', 'P008', '2024-01-01 09:30:00');
INSERT INTO `sys_role_permission` VALUES (39, 'R004', 'P000', '2024-01-03 10:15:00');
INSERT INTO `sys_role_permission` VALUES (40, 'R004', 'P003', '2024-01-03 10:15:00');
INSERT INTO `sys_role_permission` VALUES (41, 'R004', 'P0010', '2024-01-03 10:15:00');
INSERT INTO `sys_role_permission` VALUES (42, 'R004', 'P0030', '2024-01-03 10:15:00');
INSERT INTO `sys_role_permission` VALUES (43, 'R004', 'P011', '2024-01-03 10:15:00');
INSERT INTO `sys_role_permission` VALUES (44, 'R004', 'P001', '2025-11-13 21:34:45');
INSERT INTO `sys_role_permission` VALUES (45, 'R004', 'P002', '2025-11-13 21:34:45');

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user`  (
  `user_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '系统用户唯一标识，操作员ID',
  `user_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名（登录账号）',
  `real_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '真实姓名',
  `password` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码，加密存储（如MD5+盐值）',
  `role_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '角色ID，关联角色表',
  `department` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '所属部门，如“理财销售部”“运营部”',
  `position` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '职位，如“柜台操作员”“运营专员”',
  `contact_phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '联系电话',
  `user_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户状态，如“启用”“禁用”',
  `last_login_time` datetime NULL DEFAULT NULL COMMENT '最后登录时间',
  `password_expire_time` datetime NOT NULL COMMENT '密码过期时间（默认90天）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '用户创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '用户信息更新时间',
  `deleted_at` datetime NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`) USING BTREE,
  UNIQUE INDEX `uk_user_name`(`user_name` ASC) USING BTREE,
  INDEX `idx_role_id`(`role_id` ASC) USING BTREE,
  INDEX `idx_user_status`(`user_status` ASC) USING BTREE,
  CONSTRAINT `fk_su_role_id` FOREIGN KEY (`role_id`) REFERENCES `sys_role` (`role_id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '系统用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_user
-- ----------------------------
INSERT INTO `sys_user` VALUES ('e1055e3e3608424282ae6875c5943405', '时勇', '仵治涛', '66998b5d76e0cd3574777c20c065cf1e', 'R001', 'sed ln', 'Ut', '052 9167 8108', '已删除', NULL, '2026-02-26 14:09:05', '2025-11-28 14:09:04', '2025-11-28 17:05:58', '2025-11-28 17:05:58');
INSERT INTO `sys_user` VALUES ('U001', 'admin', '张三', 'e10adc3949ba59abbe56e057f20f883e', 'R001', '技术部', '系统管理员', '13800138000', '启用', '2024-05-20 08:30:00', '2024-07-30 00:00:00', '2024-01-01 09:00:00', '2025-11-29 20:06:03', NULL);
INSERT INTO `sys_user` VALUES ('U002', 'operation_li', '李四', '1bbd886460827015e5d605ed44252251', 'R002', '理财运营部', '运营主管', '13900139000', '启用', '2024-05-20 09:15:00', '2024-08-15 00:00:00', '2024-01-01 09:00:00', '2025-11-29 20:03:59', '2025-11-29 20:03:57');
INSERT INTO `sys_user` VALUES ('U003', 'operation_wang', '夙娜', 'e10adc3949ba59abbe56e057f20f883e', 'R002', 'esse', 'id non ur', '19971312015', 'Excepteua', '2024-05-20 10:00:00', '2025-11-29 00:00:00', '2024-01-02 10:30:00', '2025-11-29 20:04:21', '2025-11-29 20:04:17');
INSERT INTO `sys_user` VALUES ('U004', 'counter_zhang', '赵六', 'e10adc3949ba59abbe56e057f20f883e', 'R003', '理财销售部', '柜台组长', '13600136000', '启用', '2024-05-20 09:45:00', '2024-07-25 00:00:00', '2024-01-01 09:00:00', '2024-05-20 09:45:00', NULL);
INSERT INTO `sys_user` VALUES ('U005', 'counter_li', '孙七', 'e10adc3949ba59abbe56e057f20f883e', 'R003', '理财销售部', '柜台操作员', '13500135000', '启用', '2024-05-20 10:30:00', '2024-09-01 00:00:00', '2024-01-03 11:00:00', '2024-05-20 10:30:00', NULL);
INSERT INTO `sys_user` VALUES ('U006', 'analyst_chen', '陈八', 'e10adc3949ba59abbe56e057f20f883e', 'R004', '数据分析部', '数据分析师', '13400134000', '启用', '2024-05-20 11:15:00', '2024-08-10 00:00:00', '2024-01-03 10:15:00', '2024-05-20 11:15:00', NULL);

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_users_username`(`username` ASC) USING BTREE,
  UNIQUE INDEX `idx_users_email`(`email` ASC) USING BTREE,
  INDEX `idx_users_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (2, 'testuser', '$2a$10$/eL5N1KlUXQcvwO6AY/leu2up0k9N1c.M3cnZRdohfoiTTcjxmLa2', 'jwupu5_ire@foxmail.com', '2025-10-09 13:28:32.047', '2025-10-09 13:28:32.047', NULL);

SET FOREIGN_KEY_CHECKS = 1;
