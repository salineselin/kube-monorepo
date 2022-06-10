package main

import (
	"strings"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/container"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// create a container registry
		reg, err := container.NewRegistry(ctx, "make-registry", &container.RegistryArgs{
			Location: pulumi.String("US"),
		})
		if err != nil {
			return err
		}

		// log the container registrys url
		ctx.Export("registry: ", reg.BucketSelfLink)

		// Create an admin serviceaccount
		saName := "gcr-admin"
		sa, err := serviceaccount.NewAccount(ctx, saName, &serviceaccount.AccountArgs{
			AccountId:   pulumi.String(saName),
			DisplayName: pulumi.String(saName),
		}, pulumi.Protect(true))
		if err != nil {
			return err
		}

		// bind it to the container registry
		_, err = storage.NewBucketIAMMember(ctx, "give-sa-bucket-permissions", &storage.BucketIAMMemberArgs{
			Role: pulumi.String("roles/storage.admin"),
			Bucket: reg.BucketSelfLink.ApplyT(func(link string) string {
				return strings.Split(link, "/v1/b/")[1]
			}).(pulumi.StringOutput),
			Member: sa.Email.ApplyT(func(Email string) string {
				return "serviceAccount:" + Email
			}).(pulumi.StringOutput),
		})
		if err != nil {
			return err
		}
		return nil
	})
}