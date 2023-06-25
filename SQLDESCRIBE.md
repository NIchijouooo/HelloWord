# 储能系统数据表说明



## Sqlite

### em_device：设备表：

| 名称              | 类型    | 是否null | 说明     |
| ----------------- | ------- | -------- | -------- |
| id                | INTEGER | NOT NULL | 设备id   |
| name              | TEXT    |          | 设备名称 |
| label             | TEXT    |          | 设备标签 |
| device_type       | TEXT    |          | 设备类型 |
| model_id          | INTEGER |          | 模型id   |
| coll_interface_id | INTEGER |          |          |
| addr              | TEXT    |          | 地址     |
| data              | TEXT    |          | json数据 |

### em_limit_config：越限配置表

| 名称              | 类型    | 是否null | 说明     |
| ----------------- | ------- | -------- | -------- |
| id                | INTEGER | NOT NULL |          |
| device_type       | TEXT    |          | 设备类型 |
| property_code     | TEXT    |          | 测点编码 |
| enable_flag       | integer |          |          |
| notify_min        | text    |          | 最小告知 |
| notify_max        | text    | NOT NULL | 最大告知 |
| notify_rule_id    | integer |          |          |
| secondary_min     | text    |          |          |
| secondary_max     | text    |          | 最小次要 |
| secondary_rule_id | integer |          | 最大次要 |
| serious_min       | text    |          | 最小严重 |
| serious_max       | text    |          | 最大严重 |
| serious_rule_id   | integer |          |          |
| urgent_min        | text    |          | 最小紧急 |
| urgent_max        | text    |          | 最大紧急 |
| urgent_rule_id    | integer |          |          |
| del_flag          | integer | NOT NULL |          |
| create_time       | TEXT    |          |          |
| update_time       | TEXT    |          |          |

## Taos

### charge_discharge：每天充电量表

| 名称               | 类型       | 是否null | 说明   |
| ------------------ | ---------- | -------- | ------ |
| ts                 | time stamp | NOT NULL |        |
| charge_capacity    | double     |          | 充电量 |
| discharge_capacity | double     |          | 放电量 |
| profit             | double     |          |        |
| device_id          | int        |          | 设备id |

### charge_discharge_hour：每小时充电量表

​	表说明：每个小时采集一次，同时采集充放电信息

| 名称               | 类型       | 是否null | 说明   |
| ------------------ | ---------- | -------- | ------ |
| ts                 | time stamp | NOT NULL |        |
| charge_capacity    | double     |          | 充电量 |
| discharge_capacity | double     |          | 放电量 |
| profit             | double     |          |        |
| device_id          | int        |          | 设备id |