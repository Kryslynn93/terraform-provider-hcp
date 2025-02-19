---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "hcp_packer_bucket_names Data Source - terraform-provider-hcp"
subcategory: ""
description: |-
  The Packer Bucket Names data source gets the names of all of the buckets in a single HCP Packer registry.
---

# hcp_packer_bucket_names (Data Source)

The Packer Bucket Names data source gets the names of all of the buckets in a single HCP Packer registry.

## Example Usage

```terraform
data "hcp_packer_bucket_names" "all" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `project_id` (String) The ID of the HCP project where the HCP Packer registry is located.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.
- `names` (Set of String) A set of names for all buckets in the HCP Packer registry.
- `organization_id` (String) The ID of the organization where the HCP Packer registry is located.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `default` (String)
