package settings

import (
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/defaults"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/output"
	"github.com/mailchain/mailchain/cmd/mailchain/internal/settings/values"
	"github.com/mailchain/mailchain/stores/s3store"
)

func sentStoreS3(s values.Store) *SentStoreS3 {
	return &SentStoreS3{
		Bucket:          values.NewDefaultString(defaults.Empty, s, "sentstore.s3.bucket"),
		Region:          values.NewDefaultString(defaults.Empty, s, "sentstore.s3.region"),
		AccessKeyID:     values.NewDefaultString(defaults.Empty, s, "sentstore.s3.accessKeyId"),
		SecretAccessKey: values.NewDefaultString(defaults.Empty, s, "sentstore.s3.secretAccessKey"),
	}
}

// SentStoreS3 configuration element.
type SentStoreS3 struct {
	Bucket          values.String
	Region          values.String
	AccessKeyID     values.String
	SecretAccessKey values.String
}

// Produce `s3store.Sent` based on configuration settings.
func (s SentStoreS3) Produce() (*s3store.Sent, error) {
	return s3store.NewSent(s.Region.Get(), s.Bucket.Get(), s.AccessKeyID.Get(), s.SecretAccessKey.Get())
}

// Output configuration as an `output.Element` for use in exporting configuration.
func (s SentStoreS3) Output() output.Element {
	return output.Element{
		FullName: "sentstore.s3",
		Attributes: []output.Attribute{
			s.Bucket.Attribute(),
			s.Region.Attribute(),
			s.AccessKeyID.Attribute(),
			s.SecretAccessKey.Attribute(),
		},
	}
}
