package aws

import "github.com/infracost/infracost/internal/schema"

var ResourceRegistry []*schema.RegistryItem = []*schema.RegistryItem{
	GetAPIGatewayRestAPIRegistryItem(),
	GetAPIGatewayStageRegistryItem(),
	GetAPIGatewayv2ApiRegistryItem(),
	GetAutoscalingGroupRegistryItem(),
	GetCloudwatchDashboardRegistryItem(),
	GetCloudwatchLogGroupItem(),
	GetCloudwatchMetricAlarmRegistryItem(),
	GetDBInstanceRegistryItem(),
	GetDMSRegistryItem(),
	GetDocDBClusterInstanceRegistryItem(),
	GetDXGatewayAssociationRegistryItem(),
	GetDynamoDBTableRegistryItem(),
	GetEBSSnapshotCopyRegistryItem(),
	GetEBSSnapshotRegistryItem(),
	GetEBSVolumeRegistryItem(),
	GetEC2TransitGatewayPeeringAttachmentRegistryItem(),
	GetEC2TransitGatewayVpcAttachmentRegistryItem(),
	GetEC2TrafficMirroSessionRegistryItem(),
	GetECRRegistryItem(),
	GetECSServiceRegistryItem(),
	GetEIPRegistryItem(),
	GetElastiCacheClusterItem(),
	GetElastiCacheReplicationGroupItem(),
	GetElasticsearchDomainRegistryItem(),
	GetELBRegistryItem(),
	GetInstanceRegistryItem(),
	GetLambdaFunctionRegistryItem(),
	GetLBRegistryItem(),
	GetLightsailInstanceRegistryItem(),
	GetMSKClusterRegistryItem(),
	GetALBRegistryItem(),
	GetNATGatewayRegistryItem(),
	GetRDSClusterInstanceRegistryItem(),
	GetRoute53RecordRegistryItem(),
	GetRoute53ZoneRegistryItem(),
	GetS3BucketRegistryItem(),
	GetS3BucketAnalyticsConfigurationRegistryItem(),
	GetS3BucketInventoryRegistryItem(),
	GetSNSTopicRegistryItem(),
	GetSNSTopicSubscriptionRegistryItem(),
	GetSQSQueueRegistryItem(),
	GetNewEKSNodeGroupItem(),
	GetNewEKSFargateProfileItem(),
	GetNewEKSClusterItem(),
	GetNewKMSKeyRegistryItem(),
	GetNewKMSExternalKeyRegistryItem(),
	GetVpcEndpointRegistryItem(),
}

// FreeResources grouped alphabetically
var FreeResources []string = []string{
	// AWS API Gateway Rest APIs
	"aws_api_gateway_account",
	"aws_api_gateway_api_key",
	"aws_api_gateway_authorizer",
	"aws_api_gateway_base_path_mapping",
	"aws_api_gateway_client_certificate",
	"aws_api_gateway_deployment",
	"aws_api_gateway_documentation_part",
	"aws_api_gateway_documentation_version",
	"aws_api_gateway_domain_name",
	"aws_api_gateway_response",
	"aws_api_gateway_integration",
	"aws_api_gateway_method",
	"aws_api_gateway_method_response",
	"aws_api_gateway_method_settings",
	"aws_api_gateway_model",
	"aws_api_gateway_request_validator",
	"aws_api_gateway_resource",
	"aws_api_gateway_usage_plan",
	"aws_api_gateway_usage_plan_key",
	"aws_api_gateway_vpc_link",

	// AWS API Gateway v2 HTTP & Websocket API.
	"aws_apigatewayv2_api_mapping",
	"aws_apigatewayv2_authorizer",
	"aws_apigatewayv2_deployment",
	"aws_apigatewayv2_domain_name",
	"aws_apigatewayv2_integration",
	"aws_apigatewayv2_integration_response",
	"aws_apigatewayv2_model",
	"aws_apigatewayv2_route",
	"aws_apigatewayv2_route_response",
	"aws_apigatewayv2_stage",
	"aws_apigatewayv2_vpc_link",

	// AWS EC2,
	"aws_ec2_traffic_mirror_filter",
	"aws_ec2_traffic_mirror_filter_rule",
	"aws_ec2_traffic_mirror_target",

	// AWS Cloudwatch
	"aws_cloudwatch_log_destination",
	"aws_cloudwatch_log_destination_policy",
	"aws_cloudwatch_log_metric_filter",
	"aws_cloudwatch_log_resource_policy",
	"aws_cloudwatch_log_stream",
	"aws_cloudwatch_log_subscription_filter",

	// AWS ECR
	"aws_ecr_lifecycle_policy",
	"aws_ecr_repository_policy",

	// AWS Elastic Load Balancing
	"aws_alb_listener",
	"aws_alb_listener_certificate",
	"aws_alb_listener_rule",
	"aws_alb_target_group",
	"aws_alb_target_group_attachment",
	"aws_lb_listener",
	"aws_lb_listener_certificate",
	"aws_lb_listener_rule",
	"aws_lb_target_group",
	"aws_lb_target_group_attachment",
	"aws_app_cookie_stickiness_policy",
	"aws_elb_attachment",
	"aws_lb_cookie_stickiness_policy",
	"aws_lb_ssl_negotiation_policy",
	"aws_load_balancer_backend_server_policy",
	"aws_load_balancer_listener_policy",
	"aws_load_balancer_policy",

	// AWS Elasticache
	"aws_elasticache_parameter_group",
	"aws_elasticache_security_group",
	"aws_elasticache_subnet_group",

	// AWS IAM aws_iam_* resources
	"aws_iam_access_key",
	"aws_iam_account_alias",
	"aws_iam_account_alias",
	"aws_iam_account_password_policy",
	"aws_iam_group",
	"aws_iam_group",
	"aws_iam_group_membership",
	"aws_iam_group_policy",
	"aws_iam_group_policy_attachment",
	"aws_iam_instance_profile",
	"aws_iam_instance_profile",
	"aws_iam_openid_connect_provider",
	"aws_iam_policy",
	"aws_iam_policy",
	"aws_iam_policy_attachment",
	"aws_iam_role",
	"aws_iam_role",
	"aws_iam_role_policy",
	"aws_iam_role_policy_attachment",
	"aws_iam_saml_provider",
	"aws_iam_server_certificate",
	"aws_iam_server_certificate",
	"aws_iam_service_linked_role",
	"aws_iam_user",
	"aws_iam_user",
	"aws_iam_user_group_membership",
	"aws_iam_user_login_profile",
	"aws_iam_user_policy",
	"aws_iam_user_policy_attachment",
	"aws_iam_user_ssh_key",

	// AWS KMS
	"aws_kms_alias",
	"aws_kms_ciphertext",
	"aws_kms_grant",

	// AWS Others
	"aws_db_instance_role_association",
	"aws_db_parameter_group",
	"aws_db_subnet_group",
	"aws_dms_replication_subnet_group",
	"aws_dms_replication_task",
	"aws_docdb_cluster_parameter_group",
	"aws_ecs_cluster",
	"aws_ecs_task_definition",
	"aws_eip_association",
	"aws_elasticsearch_domain_policy",
	"aws_key_pair",
	"aws_lambda_function_event_invoke_config",
	"aws_launch_configuration",
	"aws_launch_template",
	"aws_lightsail_domain",
	"aws_lightsail_key_pair",
	"aws_lightsail_static_ip",
	"aws_lightsail_static_ip_attachment",
	"aws_msk_configuration",
	"aws_rds_cluster",
	"aws_rds_cluster_endpoint",
	"aws_rds_cluster_parameter_group",
	"aws_route53_zone_association",
	"aws_sqs_queue_policy",
	"aws_volume_attachment",

	// AWS S3
	"aws_s3_access_point",
	"aws_s3_account_public_access_block",
	"aws_s3_bucket_metric",
	"aws_s3_bucket_notification",
	"aws_s3_bucket_object", // Costs are shown at the bucket level
	"aws_s3_bucket_ownership_controls",
	"aws_s3_bucket_policy",
	"aws_s3_bucket_public_access_block",

	// AWS SNS
	"aws_sns_platform_application",
	"aws_sns_sms_preferences",
	"aws_sns_topic_policy",

	// AWS VPC
	"aws_customer_gateway",
	"aws_default_network_acl",
	"aws_default_route_table",
	"aws_default_security_group",
	"aws_default_subnet",
	"aws_default_vpc",
	"aws_default_vpc_dhcp_options",
	"aws_egress_only_internet_gateway",
	"aws_flow_log",
	"aws_internet_gateway",
	"aws_main_route_table_association",
	"aws_network_acl",
	"aws_network_acl_rule",
	"aws_network_interface",
	"aws_network_interface_attachment",
	"aws_network_interface_sg_attachment",
	"aws_route",
	"aws_route_table",
	"aws_route_table_association",
	"aws_security_group",
	"aws_security_group_rule",
	"aws_subnet",
	"aws_vpc",
	"aws_vpc_dhcp_options",
	"aws_vpc_dhcp_options_association",
	"aws_vpc_endpoint_connection_notification",
	"aws_vpc_endpoint_route_table_association",
	"aws_vpc_endpoint_service",
	"aws_vpc_endpoint_service_allowed_principal",
	"aws_vpc_endpoint_subnet_association",
	"aws_vpc_ipv4_cidr_block_association",
	"aws_vpc_peering_connection",
	"aws_vpc_peering_connection_accepter",
	"aws_vpc_peering_connection_options",
	"aws_vpn_connection_route",
	"aws_vpn_gateway",
	"aws_vpn_gateway_attachment",
	"aws_vpn_gateway_route_propagation",

	// Hashicorp
	"null_resource",
	"local_file",
	"template_dir",
	"random_id",
	"random_integer",
	"random_password",
	"random_pet",
	"random_shuffle",
	"random_string",
	"random_uuid",
	"tls_locally_signed_cert",
	"tls_private_key",
	"tls_self_signed_cert",
	"time_offset",
	"time_rotating",
	"time_sleep",
	"time_static",
}
