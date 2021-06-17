---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "appstream_fleet Resource - terraform-provider-appstream"
subcategory: ""
description: |-
  
---

# appstream_fleet (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **compute_capacity** (Block List, Min: 1) (see [below for nested schema](#nestedblock--compute_capacity))
- **instance_type** (String)
- **name** (String)

### Optional

- **description** (String)
- **disconnect_timeout** (Number)
- **display_name** (String)
- **domain_info** (Block List) (see [below for nested schema](#nestedblock--domain_info))
- **enable_default_internet_access** (Boolean)
- **fleet_type** (String)
- **iam_role_arn** (String)
- **id** (String) The ID of this resource.
- **image_arn** (String)
- **image_name** (String)
- **max_user_duration** (Number)
- **stack_name** (String)
- **state** (String)
- **stream_view** (String)
- **tags** (Map of String)
- **vpc_config** (Block List) (see [below for nested schema](#nestedblock--vpc_config))

<a id="nestedblock--compute_capacity"></a>
### Nested Schema for `compute_capacity`

Required:

- **desired_instances** (Number)


<a id="nestedblock--domain_info"></a>
### Nested Schema for `domain_info`

Optional:

- **directory_name** (String)
- **organizational_unit_distinguished_name** (String)


<a id="nestedblock--vpc_config"></a>
### Nested Schema for `vpc_config`

Optional:

- **security_group_ids** (Set of String)
- **subnet_ids** (Set of String)

