###
INSERT INTO `app`(`id`, `name`, `private_key`, `create_time`, `update_time`) VALUES (1, 'APP-01', '-----BEGIN RSA PRIVATE KEY-----\nMIICWwIBAAKBgQDcGsUIIAINHfRTdMmgGwLrjzfMNSrtgIf4EGsNaYwmC1GjF/bM\nh0Mcm10oLhNrKNYCTTQVGGIxuc5heKd1gOzb7bdTnCDPPZ7oV7p1B9Pud+6zPaco\nqDz2M24vHFWYY2FbIIJh8fHhKcfXNXOLovdVBE7Zy682X1+R1lRK8D+vmQIDAQAB\nAoGAeWAZvz1HZExca5k/hpbeqV+0+VtobMgwMs96+U53BpO/VRzl8Cu3CpNyb7HY\n64L9YQ+J5QgpPhqkgIO0dMu/0RIXsmhvr2gcxmKObcqT3JQ6S4rjHTln49I2sYTz\n7JEH4TcplKjSjHyq5MhHfA+CV2/AB2BO6G8limu7SheXuvECQQDwOpZrZDeTOOBk\nz1vercawd+J9ll/FZYttnrWYTI1sSF1sNfZ7dUXPyYPQFZ0LQ1bhZGmWBZ6a6wd9\nR+PKlmJvAkEA6o32c/WEXxW2zeh18sOO4wqUiBYq3L3hFObhcsUAY8jfykQefW8q\nyPuuL02jLIajFWd0itjvIrzWnVmoUuXydwJAXGLrvllIVkIlah+lATprkypH3Gyc\nYFnxCTNkOzIVoXMjGp6WMFylgIfLPZdSUiaPnxby1FNM7987fh7Lp/m12QJAK9iL\n2JNtwkSR3p305oOuAz0oFORn8MnB+KFMRaMT9pNHWk0vke0lB1sc7ZTKyvkEJW0o\neQgic9DvIYzwDUcU8wJAIkKROzuzLi9AvLnLUrSdI6998lmeYO9x7pwZPukz3era\nzncjRK3pbVkv0KrKfczuJiRlZ7dUzVO0b6QJr8TRAA==\n-----END RSA PRIVATE KEY-----', '2019-10-15 16:49:39', '2019-10-15 16:49:39');

###
INSERT INTO `uid`(`id`, `business_id`, `max_id`, `step`, `description`, `create_time`, `update_time`) VALUES (1, 'APP-01', 1000, 1000, '设备id', '2020-10-15 16:42:05', '2020-11-27 14:39:45');

INSERT INTO `nonsense`.`user`(`id`, `app_id`, `user_id`, `passwd`, `nickname`, `sex`, `birthday`, `mobile`, `avatar_url`, `email`, `extra`, `create_time`, `last_login_time`, `last_login_ip`, `register_ip`, `weixin_openid`) VALUES (1, 1, 1, '', '不讲武德', 1, 0, '', '', '', '', '2020-11-27 15:24:05', '2020-11-27 15:48:21', '', '', '');
INSERT INTO `nonsense`.`user`(`id`, `app_id`, `user_id`, `passwd`, `nickname`, `sex`, `birthday`, `mobile`, `avatar_url`, `email`, `extra`, `create_time`, `last_login_time`, `last_login_ip`, `register_ip`, `weixin_openid`) VALUES (2, 1, 2, '', '马保国', 1, 0, '', '', '', '', '2020-11-27 15:42:18', '2020-11-27 15:48:18', '', '', '');

###
INSERT INTO `group`(`id`, `app_id`, `group_id`, `name`, `introduction`, `user_num`, `type`, `extra`, `create_time`, `update_time`) VALUES (1, 1, 1, '传统武术交流群', '太极拳术', 1, 0, '', '2020-11-27 15:56:02', '2020-11-27 15:56:40');

###
INSERT INTO `group_user`(`id`, `app_id`, `group_id`, `user_id`, `label`, `extra`, `create_time`, `update_time`) VALUES (1, 1, 1, 2, '群主', '', '2020-11-27 15:57:01', '2020-11-27 15:57:01');


######
INSERT INTO `device`(`id`, `device_id`, `app_id`, `user_id`, `type`, `brand`, `model`, `system_version`, `sdk_version`, `status`, `conn_id`, `user_ip`, `create_time`, `update_time`) VALUES (1, 61, 1, 3, 1, 'huawei', 'huawei P30', '1.0.0', '1.0.0', 0, '127.0.0.1:60000', 8, '2020-10-23 17:11:11', '2020-10-26 09:41:24');
INSERT INTO `device`(`id`, `device_id`, `app_id`, `user_id`, `type`, `brand`, `model`, `system_version`, `sdk_version`, `status`, `conn_id`, `user_ip`, `create_time`, `update_time`) VALUES (2, 62, 1, 0, 1, 'huawei', 'huawei P30 pro', '1.0.0', '1.0.0', 0, '', 0, '2020-10-23 17:16:19', '2020-10-23 17:16:19');




