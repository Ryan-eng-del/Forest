package grom

/* gateway_admin
id
username
salt
password
create_at | update_at | delete_at |is_delete
*/

/* gateway_app
id
app_id
name
secret
white_ips
qps
qpd
create_at | update_at | delete_at |is_delete
*/

// gateway_service_access_control

/* gateway_service_grpc_rule
id
service_id
port
header_transform
*/

/* gateway_service_http_rule
id
service_id
rule_type
rule
need_strip_url
read_https
need_websocket
url_rewrite
header_transform
create_at | update_at | delete_at |is_delete
*/

/*gateway_service_tcp_rule
service_id
port
*/

/* gateway_service_info
id
loadBalance_type
service_name
service_desc
create_at | update_at | delete_at |is_delete
*/

/* gateway_service_load_balancer
id
service_id
check_method
check_timeout
check_interval
round_robin_typ
ip_list
weight_list
forbid_list
upstream_connect_timeout


*/