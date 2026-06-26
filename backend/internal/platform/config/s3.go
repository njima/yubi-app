package config

type S3 struct {
	Region                string `envconfig:"AWS_REGION" default:"ap-northeast-1"`
	BucketName            string `envconfig:"S3_BUCKET_NAME"`
	PresignedURLExpirySec int    `envconfig:"S3_PRESIGNED_URL_EXPIRY_SEC" default:"3600"`
	Endpoint              string `envconfig:"S3_ENDPOINT" default:""`
	PresignEndpoint       string `envconfig:"S3_PRESIGN_ENDPOINT" default:""`
	AccessKeyID           string `envconfig:"S3_ACCESS_KEY_ID" default:""`
	SecretAccessKey       string `envconfig:"S3_SECRET_ACCESS_KEY" default:""`
}
