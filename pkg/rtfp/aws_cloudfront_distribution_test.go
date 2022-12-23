package rtfp

import (
	"encoding/json"
	"errors"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
)

func TestNewDiffAwsCloudfrontDistribution_AwsCloudfrontDistribution(t *testing.T) {
	expected := &diffAwsCloudfrontDistribution{resourceType: "aws_cloudfront_distribution"}
	actual := newDiffAwsCloudfrontDistribution()

	assert.Equal(t, expected, actual)
}

func TestGetTerraformResourceType_AwsCloudfrontDistribution(t *testing.T) {
	expected := "aws_cloudfront_distribution"
	actual := newDiffAwsCloudfrontDistribution().getTerraformResourceType()

	assert.Equal(t, expected, actual)
}

func prepareAwsCloudfrontDistributionTestData(b, a string) tfjson.ResourceChange {
	chg := &tfjson.Change{}
	json.Unmarshal([]byte(b), &chg.Before)
	json.Unmarshal([]byte(a), &chg.After)

	return tfjson.ResourceChange{
		Address: "aws_cloudfront_distribution.test",
		Type:    "aws_cloudfront_distribution",
		Change:  chg,
	}
}

func TestDiffTeraformResource_AwsCloudfrontDistribution(t *testing.T) {
	rc := prepareAwsCloudfrontDistributionTestData(before, after)

	actual, err := newDiffAwsCloudfrontDistribution().diffTerraformResource(&rc)
	assert.Equal(t, "", actual)
	assert.Equal(t, nil, err)
}

func TestDiffTeraformResource_AwsCloudfrontDistribution_ErrorBadJson(t *testing.T) {
	badResourceType := prepareAwsCloudfrontDistributionTestData("{\"foo\": \"bar\"}", after)
	invalidJson := prepareAwsCloudfrontDistributionTestData(before, "bad json")

	actual, err := newDiffAwsCloudfrontDistribution().diffTerraformResource(&badResourceType)
	assert.Equal(t, "", actual)
	assert.Equal(t,
		errors.New("Error unmarshalling before of resource aws_cloudfront_distribution.test. Details: Missing ARN field - likely invalid resource"),
		err,
	)

	actual, err = newDiffAwsCloudfrontDistribution().diffTerraformResource(&invalidJson)
	assert.Equal(t,
		errors.New("Error unmarshalling before of resource aws_cloudfront_distribution.test. Details: Missing ARN field - likely invalid resource"),
		err,
	)
}

func TestHandleMapstructureDecodeErrors(t *testing.T) {
	addr := "aws_cloudfront_distribution.test"
	cfdWithArn := awsCloudfrontDistribution{Arn: "arn:aws:cloudfront::123456789012:distribution/CFDISTROID"}

	// Happy path
	assert.Equal(t, nil, handleMapstructureDecodeErrors(nil, cfdWithArn, addr))

	// soft failure
	assert.Equal(t, nil, handleMapstructureDecodeErrors(errors.New("test decode failure"), cfdWithArn, addr))

	// failed decode, empty key attribute, no other error
	assert.Equal(t,
		errors.New("Error unmarshalling before of resource aws_cloudfront_distribution.test. Details: Missing ARN field - likely invalid resource"),
		handleMapstructureDecodeErrors(nil, awsCloudfrontDistribution{}, addr),
	)

	// failed decode, empty key attribute + other error
	assert.Equal(t,
		errors.New("Error unmarshalling before of resource aws_cloudfront_distribution.test. Details: Test error"),
		handleMapstructureDecodeErrors(errors.New("Test error"), awsCloudfrontDistribution{}, addr),
	)

}

var (
	// Only changes are modifying order of some arrays
	before = `
	{
		"aliases": [
			"www.zapier-staging.com",
			"zapier-staging.com"
		],
		"arn": "arn:aws:cloudfront::996097627176:distribution/E276WFITWIO0NZ",
		"caller_reference": "terraform-20200917103530068200000001",
		"comment": null,
		"custom_error_response": [],
		"default_cache_behavior": [
			{
				"allowed_methods": [
					"DELETE",
					"GET",
					"HEAD",
					"OPTIONS",
					"PATCH",
					"POST",
					"PUT"
				],
				"cache_policy_id": "9bbb9d47-3d6b-4707-aa35-5cff2b5bb81e",
				"cached_methods": [
					"GET",
					"HEAD"
				],
				"compress": false,
				"default_ttl": 0,
				"field_level_encryption_id": "",
				"forwarded_values": [],
				"function_association": [
					{
						"event_type": "viewer-request",
						"function_arn": "arn:aws:cloudfront::996097627176:function/staging-rebrand-request"
					},
					{
						"event_type": "viewer-response",
						"function_arn": "arn:aws:cloudfront::996097627176:function/staging-rebrand-response"
					}
				],
				"lambda_function_association": [],
				"max_ttl": 0,
				"min_ttl": 0,
				"origin_request_policy_id": "61c8fd07-3977-4794-af44-22b4d1bc68f3",
				"realtime_log_config_arn": "",
				"response_headers_policy_id": "",
				"smooth_streaming": false,
				"target_origin_id": "k8s",
				"trusted_key_groups": [],
				"trusted_signers": [],
				"viewer_protocol_policy": "allow-all"
			}
		],
		"default_root_object": "",
		"domain_name": "doz3t8yyk8ydc.cloudfront.net",
		"enabled": true,
		"etag": "E22BH3I0KLHSC8",
		"hosted_zone_id": "Z2FDTNDATAQYW2",
		"http_version": "http2",
		"id": "E276WFITWIO0NZ",
		"in_progress_validation_batches": 0,
		"is_ipv6_enabled": true,
		"last_modified_time": "2022-10-03 21:05:25.303 +0000 UTC",
		"logging_config": [
			{
				"bucket": "zapier-cloudfront-logs.s3.amazonaws.com",
				"include_cookies": false,
				"prefix": "zapier-staging-com"
			}
		],
		"ordered_cache_behavior": [
			{
				"allowed_methods": [
					"GET",
					"HEAD",
					"OPTIONS"
				],
				"cache_policy_id": "db5f6a88-9a61-4503-8559-aee24edffdb7",
				"cached_methods": [
					"GET",
					"HEAD",
					"OPTIONS"
				],
				"compress": false,
				"default_ttl": 0,
				"field_level_encryption_id": "",
				"forwarded_values": [],
				"function_association": [],
				"lambda_function_association": [],
				"max_ttl": 0,
				"min_ttl": 0,
				"origin_request_policy_id": "91f21ae8-37b5-4c07-a574-9e0acf746615",
				"path_pattern": "/partner/_next*",
				"realtime_log_config_arn": "",
				"response_headers_policy_id": "",
				"smooth_streaming": false,
				"target_origin_id": "partner-staging-subpath",
				"trusted_key_groups": [],
				"trusted_signers": [],
				"viewer_protocol_policy": "redirect-to-https"
			},
			{
				"allowed_methods": [
					"DELETE",
					"GET",
					"HEAD",
					"OPTIONS",
					"PATCH",
					"POST",
					"PUT"
				],
				"cache_policy_id": "9bbb9d47-3d6b-4707-aa35-5cff2b5bb81e",
				"cached_methods": [
					"GET",
					"HEAD",
					"OPTIONS"
				],
				"compress": true,
				"default_ttl": 0,
				"field_level_encryption_id": "",
				"forwarded_values": [],
				"function_association": [
					{
						"event_type": "viewer-request",
						"function_arn": "arn:aws:cloudfront::996097627176:function/staging-rebrand-request"
					},
					{
						"event_type": "viewer-response",
						"function_arn": "arn:aws:cloudfront::996097627176:function/staging-rebrand-response"
					}
				],
				"lambda_function_association": [],
				"max_ttl": 0,
				"min_ttl": 0,
				"origin_request_policy_id": "61c8fd07-3977-4794-af44-22b4d1bc68f3",
				"path_pattern": "/sign-up/team-member*",
				"realtime_log_config_arn": "",
				"response_headers_policy_id": "",
				"smooth_streaming": false,
				"target_origin_id": "k8s",
				"trusted_key_groups": [],
				"trusted_signers": [],
				"viewer_protocol_policy": "redirect-to-https"
			}
		],
		"origin": [
			{
				"connection_attempts": 3,
				"connection_timeout": 10,
				"custom_header": [
					{
						"name": "X-Forwarded-Host",
						"value": "zapier-staging.com"
					},
					{
						"name": "x-vercel-proxy-secret",
						"value": "REDACTED"
					}
				],
				"custom_origin_config": [
					{
						"http_port": 80,
						"https_port": 443,
						"origin_keepalive_timeout": 60,
						"origin_protocol_policy": "https-only",
						"origin_read_timeout": 60,
						"origin_ssl_protocols": [
							"TLSv1.2"
						]
					}
				],
				"domain_name": "account-management.vercel.zapier-staging.com",
				"origin_id": "account-management",
				"origin_path": "",
				"origin_shield": [],
				"s3_origin_config": []
			},
			{
				"connection_attempts": 3,
				"connection_timeout": 10,
				"custom_header": [
					{
						"name": "X-Forwarded-Host",
						"value": "zapier-staging.com"
					},
					{
						"name": "x-vercel-proxy-secret",
						"value": "REDACTED"
					}
				],
				"custom_origin_config": [
					{
						"http_port": 80,
						"https_port": 443,
						"origin_keepalive_timeout": 60,
						"origin_protocol_policy": "https-only",
						"origin_read_timeout": 60,
						"origin_ssl_protocols": [
							"TLSv1.2"
						]
					}
				],
				"domain_name": "app-management.vercel.zapier-staging.com",
				"origin_id": "app-management",
				"origin_path": "",
				"origin_shield": [],
				"s3_origin_config": []
			}
		],
		"origin_group": [],
		"price_class": "PriceClass_All",
		"restrictions": [
			{
				"geo_restriction": [
					{
						"locations": [],
						"restriction_type": "none"
					}
				]
			}
		],
		"retain_on_delete": false,
		"status": "Deployed",
		"tags": {
			"Name": "zapier-staging-com-cdn",
			"environment": "staging",
			"managed_by": "terraform",
			"service": "zapier"
		},
		"tags_all": {
			"Name": "zapier-staging-com-cdn",
			"environment": "staging",
			"managed_by": "terraform",
			"service": "zapier"
		},
		"trusted_key_groups": [
			{
				"enabled": false,
				"items": []
			}
		],
		"trusted_signers": [
			{
				"enabled": false,
				"items": []
			}
		],
		"viewer_certificate": [
			{
				"acm_certificate_arn": "arn:aws:acm:us-east-1:996097627176:certificate/8a8ff10c-19a1-4f02-96f0-28586a2c4f41",
				"cloudfront_default_certificate": false,
				"iam_certificate_id": "",
				"minimum_protocol_version": "TLSv1.2_2018",
				"ssl_support_method": "sni-only"
			}
		],
		"wait_for_deployment": true,
		"web_acl_id": "some:arn"
	}
	`
	after = `
	{
		"aliases": [
			"zapier-staging.com",
			"www.zapier-staging.com"
		],
		"arn": "arn:aws:cloudfront::996097627176:distribution/E276WFITWIO0NZ",
		"caller_reference": "terraform-20200917103530068200000001",
		"comment": null,
		"custom_error_response": [],
		"default_cache_behavior": [
			{
				"allowed_methods": [
					"DELETE",
					"GET",
					"HEAD",
					"OPTIONS",
					"PATCH",
					"POST",
					"PUT"
				],
				"cache_policy_id": "9bbb9d47-3d6b-4707-aa35-5cff2b5bb81e",
				"cached_methods": [
					"GET",
					"HEAD"
				],
				"compress": false,
				"default_ttl": 0,
				"field_level_encryption_id": "",
				"forwarded_values": [],
				"function_association": [
					{
						"event_type": "viewer-request",
						"function_arn": "arn:aws:cloudfront::996097627176:function/staging-rebrand-request"
					},
					{
						"event_type": "viewer-response",
						"function_arn": "arn:aws:cloudfront::996097627176:function/staging-rebrand-response"
					}
				],
				"lambda_function_association": [],
				"max_ttl": 0,
				"min_ttl": 0,
				"origin_request_policy_id": "61c8fd07-3977-4794-af44-22b4d1bc68f3",
				"realtime_log_config_arn": "",
				"response_headers_policy_id": "",
				"smooth_streaming": false,
				"target_origin_id": "k8s",
				"trusted_key_groups": [],
				"trusted_signers": [],
				"viewer_protocol_policy": "allow-all"
			}
		],
		"default_root_object": "",
		"domain_name": "doz3t8yyk8ydc.cloudfront.net",
		"enabled": true,
		"etag": "E22BH3I0KLHSC8",
		"hosted_zone_id": "Z2FDTNDATAQYW2",
		"http_version": "http2",
		"id": "E276WFITWIO0NZ",
		"in_progress_validation_batches": 0,
		"is_ipv6_enabled": true,
		"last_modified_time": "2022-10-03 21:05:25.303 +0000 UTC",
		"logging_config": [
			{
				"bucket": "zapier-cloudfront-logs.s3.amazonaws.com",
				"include_cookies": false,
				"prefix": "zapier-staging-com"
			}
		],
		"ordered_cache_behavior": [
			{
				"allowed_methods": [
					"DELETE",
					"GET",
					"HEAD",
					"OPTIONS",
					"PATCH",
					"POST",
					"PUT"
				],
				"cache_policy_id": "9bbb9d47-3d6b-4707-aa35-5cff2b5bb81e",
				"cached_methods": [
					"GET",
					"HEAD",
					"OPTIONS"
				],
				"compress": true,
				"default_ttl": 0,
				"field_level_encryption_id": "",
				"forwarded_values": [],
				"function_association": [
					{
						"event_type": "viewer-request",
						"function_arn": "arn:aws:cloudfront::996097627176:function/staging-rebrand-request"
					},
					{
						"event_type": "viewer-response",
						"function_arn": "arn:aws:cloudfront::996097627176:function/staging-rebrand-response"
					}
				],
				"lambda_function_association": [],
				"max_ttl": 0,
				"min_ttl": 0,
				"origin_request_policy_id": "61c8fd07-3977-4794-af44-22b4d1bc68f3",
				"path_pattern": "/sign-up/team-member*",
				"realtime_log_config_arn": "",
				"response_headers_policy_id": "",
				"smooth_streaming": false,
				"target_origin_id": "k8s",
				"trusted_key_groups": [],
				"trusted_signers": [],
				"viewer_protocol_policy": "redirect-to-https"
			},
			{
				"allowed_methods": [
					"GET",
					"HEAD",
					"OPTIONS"
				],
				"cache_policy_id": "db5f6a88-9a61-4503-8559-aee24edffdb7",
				"cached_methods": [
					"GET",
					"HEAD",
					"OPTIONS"
				],
				"compress": false,
				"default_ttl": 0,
				"field_level_encryption_id": "",
				"forwarded_values": [],
				"function_association": [],
				"lambda_function_association": [],
				"max_ttl": 0,
				"min_ttl": 0,
				"origin_request_policy_id": "91f21ae8-37b5-4c07-a574-9e0acf746615",
				"path_pattern": "/partner/_next*",
				"realtime_log_config_arn": "",
				"response_headers_policy_id": "",
				"smooth_streaming": false,
				"target_origin_id": "partner-staging-subpath",
				"trusted_key_groups": [],
				"trusted_signers": [],
				"viewer_protocol_policy": "redirect-to-https"
			}
		],
		"origin": [
			{
				"connection_attempts": 3,
				"connection_timeout": 10,
				"custom_header": [
					{
						"name": "X-Forwarded-Host",
						"value": "zapier-staging.com"
					},
					{
						"name": "x-vercel-proxy-secret",
						"value": "REDACTED"
					}
				],
				"custom_origin_config": [
					{
						"http_port": 80,
						"https_port": 443,
						"origin_keepalive_timeout": 60,
						"origin_protocol_policy": "https-only",
						"origin_read_timeout": 60,
						"origin_ssl_protocols": [
							"TLSv1.2"
						]
					}
				],
				"domain_name": "app-management.vercel.zapier-staging.com",
				"origin_id": "app-management",
				"origin_path": "",
				"origin_shield": [],
				"s3_origin_config": []
			},
			{
				"connection_attempts": 3,
				"connection_timeout": 10,
				"custom_header": [
					{
						"name": "X-Forwarded-Host",
						"value": "zapier-staging.com"
					},
					{
						"name": "x-vercel-proxy-secret",
						"value": "REDACTED"
					}
				],
				"custom_origin_config": [
					{
						"http_port": 80,
						"https_port": 443,
						"origin_keepalive_timeout": 60,
						"origin_protocol_policy": "https-only",
						"origin_read_timeout": 60,
						"origin_ssl_protocols": [
							"TLSv1.2"
						]
					}
				],
				"domain_name": "account-management.vercel.zapier-staging.com",
				"origin_id": "account-management",
				"origin_path": "",
				"origin_shield": [],
				"s3_origin_config": []
			}
		],
		"origin_group": [],
		"price_class": "PriceClass_All",
		"restrictions": [
			{
				"geo_restriction": [
					{
						"locations": [],
						"restriction_type": "none"
					}
				]
			}
		],
		"retain_on_delete": false,
		"status": "Deployed",
		"tags": {
			"Name": "zapier-staging-com-cdn",
			"environment": "staging",
			"managed_by": "terraform",
			"service": "zapier"
		},
		"tags_all": {
			"Name": "zapier-staging-com-cdn",
			"environment": "staging",
			"managed_by": "terraform",
			"service": "zapier"
		},
		"trusted_key_groups": [
			{
				"enabled": false,
				"items": []
			}
		],
		"trusted_signers": [
			{
				"enabled": false,
				"items": []
			}
		],
		"viewer_certificate": [
			{
				"acm_certificate_arn": "arn:aws:acm:us-east-1:996097627176:certificate/8a8ff10c-19a1-4f02-96f0-28586a2c4f41",
				"cloudfront_default_certificate": false,
				"iam_certificate_id": "",
				"minimum_protocol_version": "TLSv1.2_2018",
				"ssl_support_method": "sni-only"
			}
		],
		"wait_for_deployment": true,
		"web_acl_id": "some:arn"
	}
	`
)
