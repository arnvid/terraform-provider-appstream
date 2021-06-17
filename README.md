## Terraform AppStream 2.0 provider
# terraform-provider-appstream

[![release](https://github.com/arnvid/terraform-provider-appstream/actions/workflows/release.yml/badge.svg)](https://github.com/arnvid/terraform-provider-appstream/actions/workflows/release.yml)


+This is an updated Terraform Appstream provider that covers almost all AWS SDK operations on Appstream service. This fork attempts to build a suitable provider that conforms as clossely as possible terraform provider standards.


# Provider usage

```
$ go build ~/.terraform.d/plugins/terraform-provider-appstream (for terraform 0.12)
$ go build ~/.terraform.d/plugin-cache/registry.terraform.io/arnvid/appstream/1.0.x/<platform> (for terraform 0.13+)
$ terraform init
$ terraform plan
$ terraform apply
```

#Development notes
Several other terraform provider projects have been used to reference how a module should be written,
The goal of this version is to be able to run properly with Terraform Cloud and Terraform Enterprise.
Along side with removing the need for access and secret key in variables and only pass the necessary
to be assumed.

Large portions of code for authentication in config.go & provider.go is from:
https://github.com/terraform-providers/terraform-provider-aws

## Authors/Contributors/Forks
This fork is by:
[Arnvid Karstad](https://github.com/arnvid)

Contributors:
[Konstantin Odnoralov](https://github.com/hostmit)
[Syed Hussain](https://github.com/suhussai)

Original code from:
https://github.com/ops-guru/terraform-provider-appstream
[Viktor Berlov](https://github.com/vktr-brlv)

Other forks ref'd:
https://github.com/bluesentry/terraform-provider-appstream
[Chris Mackubin](https://github.com/chris-mackubin)

https://github.com/nparfait/terraform-provider-appstream
[Nicolas Parfait](https://github.com/nparfait)

https://github.com/julian3xl/terraform-provider-appstream
[julian3xl](https://github.com/julian3xl)

https://github.com/PMI-DEEP-Infra/terraform-provider-appstream
[PMI DEEP Infra Team](https://github.com/PMI-DEEP-Infra)
[Jorge Canha](jorge.canha@pmi.com)
