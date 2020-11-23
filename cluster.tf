resource "example-cluster-config" "my-cluster" {
    api_key = "<my api key>"
    cluster_name = "my cluster"
    config = <<EOF
 {
  "name": "my-third-api-deployment",
  "resources": {
    "elasticsearch": [
      {
        "region": "us-east-1",
        "ref_id": "main-elasticsearch",
        "plan": {
          "cluster_topology": [
            {
              "node_type": {
                "data": true,
                "master": true,
                "ingest": true
              },
              "instance_configuration_id": "aws.data.highio.i3",
              "zone_count": 2,
              "size": {
                "resource": "memory",
                "value": 4096
              }
            }
          ],
          "elasticsearch": {
            "version": "7.8.1"
          },
          "deployment_template": {
            "id": "aws-io-optimized"
          }
        }
      }
    ],
    "kibana": [
      {
        "region": "us-east-1",
        "elasticsearch_cluster_ref_id": "main-elasticsearch",
        "ref_id": "main-kibana",
        "plan": {
          "cluster_topology": [
            {
              "instance_configuration_id": "aws.kibana.r4",
              "zone_count": 1,
              "size": {
                "resource": "memory",
                "value": 1024
              }
            }
          ],
          "kibana": {
            "version": "7.8.1"
          }
        }
      }
    ],
    "apm": [
      {
        "region": "us-east-1",
        "elasticsearch_cluster_ref_id": "main-elasticsearch",
        "ref_id": "main-apm",
        "plan": {
          "apm": {
            "version": "7.8.1"
          },
          "cluster_topology": [
            {
              "instance_configuration_id": "aws.apm.r4",
              "zone_count": 1,
              "size": {
                "value": 512,
                "resource": "memory"
              }
            }
          ]
        }
      }
    ]
  }
}
EOF
}
