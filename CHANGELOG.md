## WIP (**UNRELEASED**)

FEATURES:
* [nparfait] appstream/resource_fleet.go - state, stream_view
* [nparfait] appstream/resource_image_builder.go - access_endpoints
* [nparfait] appstream/resource_stack.go - access_endpoints, application_settings
* [julian3xl] appstream/resource_stack.go - embedded_domain 
* [julian3xl] appstream/resource_stack_attachment.go - new resource 
* [nparfait]  appstream/appstream_usage_report_subscription.go - new resource

ENHANCEMENTS:
* [arnvid] upgraded go modules
* [arnvid] upgraded go api, sdk and interfaces
* [arnvid] appstream/provider.go - support for longer session times
* [julian3xl] appstream/resource_fleet.go - iam_role_arn -> optional
* [julian3xl] appstream/resource_fleet.go - added support for account local image_name
* [julian3xl] appstream/resource_fleet.go - update to vpc_config variable
* [julian3xl] appstream/resource_fleet.go - updated tagging
* [julian3xl] appstream/resource_fleet.go - cleaned up comments and errorlogs
* [julian3xl] appstream/resource_fleet.go - added idle_disconnect_timeout & enable_default_internet_access
* [julian3xl] appstream/resource_stack.go - updated user_settings.
* [other] appstream/resource_stack.go - multiple fixes

BUGFIXES:
* [suhussai] search for stack, fleet and image builder by name instead of listing all #5
* [julian3xl] api inconsistencies
* [arnvid] bug in optional variables forcing recreation

TOOLS:
* [nparfait] Added support for automated builds actions github actions
* [canha] some fixes to automated build tools for newer go lang

## 1.0.8 (June 15, 2020)

FEATURES:
* appstream/resource_stack.go - user_settings patch from dhruv2511

ENHANCEMENTS:
* Upgraded modules
* Updated examples

BUGFIXES:


## 1.0.7 (June 9, 2020)

FEATURES:
* Added support for role ARN

ENHANCEMENTS:

BUGFIXES:

Patch by: Konstantin Odnoralov <kodnoral@pmintl.net>

## 1.0.6 (May 27, 2020)

FEATURES:
* Added Ability to Remote Image 
* Changes to iamge_arn forces new stack

ENHANCEMENTS:
* updated tf lib to 0.12.25

BUGFIXES:
* image_name changed to image_arn

Patch by: Konstantin Odnoralov <kodnoral@pmintl.net>


## 1.0.5 (May 03, 2020)

FEATURES:
* Assume Role authentication

ENHANCEMENTS:
* authentication: Adopted AWS authentication from terraform-provider-aws
* structure: changed structure and build setup of provider

BUGFIXES: 


