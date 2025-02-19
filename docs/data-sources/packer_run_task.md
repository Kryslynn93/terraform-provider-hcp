---
page_title: "hcp_packer_run_task Data Source - terraform-provider-hcp"
subcategory: ""
description: |-
  The Packer Run Task data source gets the configuration information needed to set up an HCP Packer Registry's run task.
---

# hcp_packer_run_task (Data Source)

-> **Note:** This data source is currently in public beta.

-> **Note:** Use of this data source in the same workspace as an 
`hcp_packer_run_task` resource (pointing to the same HCP Project) is
discouraged. If this is not possible (ex: using a module containing the data
source in the same workspace as a copy of the resource), use the `depends_on`
meta-argument to mark the data source as dependent on the resource.

The Packer Run Task data source gets the configuration information needed to set up an HCP Packer Registry's run task.

## Example Usage

```terraform
data "hcp_packer_run_task" "registry" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `project_id` (String) The ID of the HCP project where the HCP Packer Registry is located. 
If not specified, the project specified in the HCP Provider config block will be used, if configured.
If a project is not configured in the HCP Provider config block, the oldest project in the organization will be used.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `endpoint_url` (String) A unique HCP Packer URL, specific to your HCP organization and HCP Packer registry. The Terraform Cloud run task will send a payload to this URL for image validation.
- `hmac_key` (String, Sensitive) A secret key that lets HCP Packer verify the run task request.
- `id` (String) The ID of this resource.
- `organization_id` (String) The ID of the HCP organization where this channel is located. Always the same as the associated channel.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `default` (String)
