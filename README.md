# terraform-EC
Spin up and tear down EC deployments

Barebones Terraform provider for spinning up and tearing down EC deployments: based on the following current projects: https://github.com/Ascendon/terraform-provider-ece and https://github.com/phillbaker/terraform-provider-elasticsearch

Required attributes:
1. cluster name
2. API key from elasticsearch cloud
3. config in the form of a string
(see cluster.tf for an example)

