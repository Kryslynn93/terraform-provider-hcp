// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"log"
	"time"

	packermodels "github.com/hashicorp/hcp-sdk-go/clients/cloud-packer-service/stable/2021-04-30/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-hcp/internal/clients"
)

func dataSourcePackerImageIteration() *schema.Resource {
	return &schema.Resource{
		Description:        "The Packer ImageIteration data source iteration gets the most recent iteration (or build) of an image, given a channel.",
		DeprecationMessage: "The `hcp_packer_image_iteration` data source is deprecated. Use `hcp_packer_image` or `hcp_packer_iteration` instead.",
		ReadContext:        dataSourcePackerImageIterationRead,
		Timeouts: &schema.ResourceTimeout{
			Default: &defaultPackerTimeout,
		},
		Schema: map[string]*schema.Schema{
			// Required inputs
			"bucket_name": {
				Description:      "The slug of the HCP Packer Registry bucket to pull from.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateSlugID,
			},
			"channel": {
				Description:      "The channel that points to the version of the image you want.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateSlugID,
			},
			// Optional inputs
			"project_id": {
				Description: `
The ID of the HCP project where the HCP Packer registry is located.
If not specified, the project specified in the HCP Provider config block will be used, if configured.
If a project is not configured in the HCP Provider config block, the oldest project in the organization will be used.`,
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},
			// computed outputs
			"organization_id": {
				Description: "The ID of the organization this HCP Packer registry is located in.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			// Actual iteration:
			"incremental_version": {
				Description: "Incremental version of this iteration",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"created_at": {
				Description: "Creation time of this iteration",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"revoke_at": {
				Description: "The revocation time of this iteration. This field will be null for any iteration that has not been revoked or scheduled for revocation.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"builds": {
				Description: "Builds for this iteration. An iteration can have more than one build if it took more than one go to build all images.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_provider": {
							Description: "Name of the cloud provider this image is stored-in, if any.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"component_type": {
							Description: "Name of the builder that built this. Ex: 'amazon-ebs.example'",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"created_at": {
							Description: "Creation time of this build.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"id": {
							Description: "HCP ID of this build.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"images": {
							Description: "Output of the build. This will contain the location or cloud identifier for your images.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"created_at": {
										Description: "Creation time of this image.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"id": {
										Description: "HCP ID of this image.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"image_id": {
										Description: "Cloud Image ID, URL string identifying this image for the builder that built it.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"region": {
										Description: "Region this image was built from. If any.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"labels": {
							Description: "Labels for this build.",
							Type:        schema.TypeMap,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"packer_run_uuid": {
							Description: "UUID of this build.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"status": {
							Description: "Status of this build. DONE means that all images tied to this build were successfully built.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"updated_at": {
							Description: "Time this build was last updated.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePackerImageIterationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client)
	loc, err := getAndUpdateLocationResourceData(d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	bucketName := d.Get("bucket_name").(string)
	channelSlug := d.Get("channel").(string)

	log.Printf("[INFO] Reading HCP Packer registry (%s) [project_id=%s, organization_id=%s, channel=%s]", bucketName, loc.ProjectID, loc.OrganizationID, channelSlug)

	channel, err := clients.GetPackerChannelBySlug(ctx, client, loc, bucketName, channelSlug)
	if err != nil {
		return diag.FromErr(err)
	}

	if channel.Iteration == nil {
		return diag.Errorf("no iteration information found for the specified channel %s", channelSlug)
	}

	iteration := channel.Iteration

	d.SetId(iteration.ID)

	if err := d.Set("incremental_version", iteration.IncrementalVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_at", iteration.CreatedAt.String()); err != nil {
		return diag.FromErr(err)
	}
	if !time.Time(iteration.RevokeAt).IsZero() {
		if err := d.Set("revoke_at", iteration.RevokeAt.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("builds", flattenPackerBuildList(iteration.Builds)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenPackerBuildList(builds []*packermodels.HashicorpCloudPackerBuild) (flattened []map[string]interface{}) {
	for _, build := range builds {
		out := map[string]interface{}{
			"cloud_provider":  build.CloudProvider,
			"component_type":  build.ComponentType,
			"created_at":      build.CreatedAt.String(),
			"id":              build.ID,
			"images":          flattenPackerBuildImagesList(build.Images),
			"labels":          build.Labels,
			"packer_run_uuid": build.PackerRunUUID,
			"status":          build.Status,
			"updated_at":      build.UpdatedAt.String(),
		}
		flattened = append(flattened, out)
	}
	return
}

func flattenPackerBuildImagesList(images []*packermodels.HashicorpCloudPackerImage) (flattened []map[string]interface{}) {
	for _, image := range images {
		out := map[string]interface{}{
			"created_at": image.CreatedAt.String(),
			"id":         image.ID,
			"image_id":   image.ImageID,
			"region":     image.Region,
		}
		flattened = append(flattened, out)
	}
	return
}
