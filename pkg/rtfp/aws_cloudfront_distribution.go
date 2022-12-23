package rtfp

import (
	"errors"
	"fmt"
	"sort"

	"github.com/google/go-cmp/cmp"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

type diffAwsCloudfrontDistribution struct {
	resourceType string
}

func newDiffAwsCloudfrontDistribution() TerraformResourceDiff {
	return &diffAwsCloudfrontDistribution{
		resourceType: "aws_cloudfront_distribution",
	}
}

func (d *diffAwsCloudfrontDistribution) getTerraformResourceType() string {
	return d.resourceType
}

func (d *diffAwsCloudfrontDistribution) diffTerraformResource(rc *tfjson.ResourceChange) (string, error) {
	var b, a awsCloudfrontDistribution

	err := mapstructure.Decode(rc.Change.Before, &b)
	if returnErr := handleMapstructureDecodeErrors(err, b, rc.Address); returnErr != nil {
		return "", returnErr
	}
	err = mapstructure.Decode(rc.Change.After, &a)
	if returnErr := handleMapstructureDecodeErrors(err, a, rc.Address); returnErr != nil {
		return "", returnErr
	}

	return cmp.Diff(
		sortAwsCloudfrontDistribution(b),
		sortAwsCloudfrontDistribution(a),
	), nil
}

func handleMapstructureDecodeErrors(err error, c awsCloudfrontDistribution, addr string) error {
	// Some fields may be missed due to sanitised fields. Only let errors bubble up if it's a dealbreaker
	// ARN will always be there, if it's missing something really blew up
	if c.Arn == "" {
		// If it's only the ARN missing, it's likely it's a different resource type. This shoudln't happen,
		// we should be filtering resource types earlier.
		if err == nil {
			err = errors.New("Missing ARN field - likely invalid resource")
		}
		return fmt.Errorf("Error unmarshalling before of resource %s. Details: %v", addr, err)
	}

	// Not end of world if we get some fields, and mapstructure errors if there are some type mismatches.
	// This often happens with sanitised fields, so we'll just log the error and continue.
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Partial error unmarshalling %s", addr))
	}
	return nil
}

// Function to sort slices in an awsCloudfrontDistribution by unique keys
func sortAwsCloudfrontDistribution(cfd awsCloudfrontDistribution) awsCloudfrontDistribution {
	sort.SliceStable(cfd.Aliases, func(i, j int) bool {
		return cfd.Aliases[i] < cfd.Aliases[j]
	})

	sort.SliceStable(cfd.CustomErrorResponse, func(i, j int) bool {
		return cfd.CustomErrorResponse[i].ResponseCode < cfd.CustomErrorResponse[j].ResponseCode
	})

	sort.SliceStable(cfd.DefaultCacheBehavior, func(i, j int) bool {
		return cfd.DefaultCacheBehavior[i].TargetOriginID < cfd.DefaultCacheBehavior[j].TargetOriginID
	})

	sort.SliceStable(cfd.LoggingConfig, func(i, j int) bool {
		return cfd.LoggingConfig[i].Bucket < cfd.LoggingConfig[j].Bucket
	})

	sort.SliceStable(cfd.OrderedCacheBehavior, func(i, j int) bool {
		return cfd.OrderedCacheBehavior[i].PathPattern < cfd.OrderedCacheBehavior[j].PathPattern
	})

	sort.SliceStable(cfd.Origin, func(i, j int) bool {
		return cfd.Origin[i].OriginID < cfd.Origin[j].OriginID
	})

	sort.SliceStable(cfd.OriginGroup, func(i, j int) bool {
		return cfd.OriginGroup[i].OriginID < cfd.OriginGroup[j].OriginID
	})

	sort.SliceStable(cfd.ViewerCertificate, func(i, j int) bool {
		return cfd.ViewerCertificate[i].AcmCertificateArn < cfd.ViewerCertificate[j].AcmCertificateArn
	})

	return cfd
}

// Struct for aws_cloudfront_distribution resource, interpreted from json
type awsCloudfrontDistribution struct {
	Aliases                     []string              `mapstructure:"aliases"`
	Arn                         string                `mapstructure:"arn"`
	CallerReference             string                `mapstructure:"caller_reference"`
	Comment                     interface{}           `mapstructure:"comment"`
	CustomErrorResponse         []CustomErrorResponse `mapstructure:"custom_error_response"`
	DefaultCacheBehavior        []CacheBehavior       `mapstructure:"default_cache_behavior"`
	DefaultRootObject           string                `mapstructure:"default_root_object"`
	DomainName                  string                `mapstructure:"domain_name"`
	Enabled                     bool                  `mapstructure:"enabled"`
	Etag                        string                `mapstructure:"etag"`
	HostedZoneID                string                `mapstructure:"hosted_zone_id"`
	HTTPVersion                 string                `mapstructure:"http_version"`
	ID                          string                `mapstructure:"id"`
	InProgressValidationBatches int                   `mapstructure:"in_progress_validation_batches"`
	IsIpv6Enabled               bool                  `mapstructure:"is_ipv6_enabled"`
	LastModifiedTime            string                `mapstructure:"last_modified_time"`
	LoggingConfig               []LoggingConfig       `mapstructure:"logging_config"`
	OrderedCacheBehavior        []CacheBehavior       `mapstructure:"ordered_cache_behavior"`
	Origin                      []Origin              `mapstructure:"origin"`
	OriginGroup                 []OriginGroup         `mapstructure:"origin_group"`
	PriceClass                  string                `mapstructure:"price_class"`
	Restrictions                []Restrictions        `mapstructure:"restrictions"`
	RetainOnDelete              bool                  `mapstructure:"retain_on_delete"`
	Status                      string                `mapstructure:"status"`
	Tags                        map[string]string     `mapstructure:"tags"`
	TagsAll                     map[string]string     `mapstructure:"tags_all"`
	TrustedKeyGroups            []TrustedKeyGroups    `mapstructure:"trusted_key_groups"`
	TrustedSigners              []TrustedSigners      `mapstructure:"trusted_signers"`
	ViewerCertificate           []ViewerCertificate   `mapstructure:"viewer_certificate"`
	WaitForDeployment           bool                  `mapstructure:"wait_for_deployment"`
	WebACLID                    string                `mapstructure:"web_acl_id"`
}

type CustomErrorResponse struct {
	ErrorCode          int    `mapstructure:"error_code"`
	ErrorCachingMinTTL int    `mapstructure:"error_caching_min_ttl"`
	ResponseCode       string `mapstructure:"response_code"`
	ResponsePagePath   string `mapstructure:"response_page_path"`
}

type FunctionAssociation struct {
	EventType   string `mapstructure:"event_type"`
	FunctionArn string `mapstructure:"function_arn"`
}

type LoggingConfig struct {
	Bucket         string `mapstructure:"bucket"`
	IncludeCookies bool   `mapstructure:"include_cookies"`
	Prefix         string `mapstructure:"prefix"`
}

type CacheBehavior struct {
	AllowedMethods            []string      `mapstructure:"allowed_methods"`
	CachePolicyID             string        `mapstructure:"cache_policy_id"`
	CachedMethods             []string      `mapstructure:"cached_methods"`
	Compress                  bool          `mapstructure:"compress"`
	DefaultTTL                float64       `mapstructure:"default_ttl"`
	FieldLevelEncryptionID    string        `mapstructure:"field_level_encryption_id"`
	ForwardedValues           []interface{} `mapstructure:"forwarded_values"`
	FunctionAssociation       []interface{} `mapstructure:"function_association"`
	LambdaFunctionAssociation []interface{} `mapstructure:"lambda_function_association"`
	MaxTTL                    float64       `mapstructure:"max_ttl"`
	MinTTL                    float64       `mapstructure:"min_ttl"`
	OriginRequestPolicyID     string        `mapstructure:"origin_request_policy_id"`
	// PathPattern not present on DefaultCacheBehaviour
	PathPattern             string        `mapstructure:"path_pattern"`
	RealtimeLogConfigArn    string        `mapstructure:"realtime_log_config_arn"`
	ResponseHeadersPolicyID string        `mapstructure:"response_headers_policy_id"`
	SmoothStreaming         bool          `mapstructure:"smooth_streaming"`
	TargetOriginID          string        `mapstructure:"target_origin_id"`
	TrustedKeyGroups        []interface{} `mapstructure:"trusted_key_groups"`
	TrustedSigners          []interface{} `mapstructure:"trusted_signers"`
	ViewerProtocolPolicy    string        `mapstructure:"viewer_protocol_policy"`
}

type CustomHeader struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

type CustomOriginConfig struct {
	HTTPPort               int      `mapstructure:"http_port"`
	HTTPSPort              int      `mapstructure:"https_port"`
	OriginKeepaliveTimeout int      `mapstructure:"origin_keepalive_timeout"`
	OriginProtocolPolicy   string   `mapstructure:"origin_protocol_policy"`
	OriginReadTimeout      int      `mapstructure:"origin_read_timeout"`
	OriginSslProtocols     []string `mapstructure:"origin_ssl_protocols"`
}

type Origin struct {
	ConnectionAttempts float64              `mapstructure:"connection_attempts"`
	ConnectionTimeout  float64              `mapstructure:"connection_timeout"`
	CustomHeader       []CustomHeader       `mapstructure:"custom_header"`
	CustomOriginConfig []CustomOriginConfig `mapstructure:"custom_origin_config"`
	DomainName         string               `mapstructure:"domain_name"`
	OriginID           string               `mapstructure:"origin_id"`
	OriginPath         string               `mapstructure:"origin_path"`
	OriginShield       []interface{}        `mapstructure:"origin_shield"`
	S3OriginConfig     []interface{}        `mapstructure:"s3_origin_config"`
}

type OriginGroup struct {
	FailoverCriteria []struct {
		StatusCodes []int `mapstructure:"status_codes"`
	} `mapstructure:"failover_criteria"`
	Member []struct {
		OriginId string `mapstructure:"origin_id"`
	} `mapstructure:"member"`
	OriginID string `mapstructure:"origin_id"`
}

type GeoRestriction struct {
	Locations       []string `mapstructure:"locations"`
	RestrictionType string   `mapstructure:"restriction_type"`
}

type Restrictions struct {
	GeoRestriction []GeoRestriction `mapstructure:"geo_restriction"`
}

type TrustedKeyGroups struct {
	Enabled bool           `mapstructure:"enabled"`
	Items   []TrustedItems `mapstructure:"items"`
}

type TrustedSigners struct {
	Enabled bool           `mapstructure:"enabled"`
	Items   []TrustedItems `mapstructure:"items"`
}

type TrustedItems struct {
	AwsAccountNumber string `mapstructure:"aws_account_number"`
	KeyPairIds       string `mapstructure:"key_pair_ids"`
}

type ViewerCertificate struct {
	AcmCertificateArn            string `mapstructure:"acm_certificate_arn"`
	CloudfrontDefaultCertificate bool   `mapstructure:"cloudfront_default_certificate"`
	IamCertificateID             string `mapstructure:"iam_certificate_id"`
	MinimumProtocolVersion       string `mapstructure:"minimum_protocol_version"`
	SslSupportMethod             string `mapstructure:"ssl_support_method"`
}
