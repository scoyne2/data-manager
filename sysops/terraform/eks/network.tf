locals {
  envs = { for tuple in regexall("(.*)=(.*)", file("../../.env")) : tuple[0] => sensitive(tuple[1]) }
}

provider "aws" {
  shared_config_files      = ["~/.aws/config"]
  shared_credentials_files = ["~/.aws/credentials"]
  profile                  = local.envs["AWS_PROFILE"]
}

data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "data-manager-vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    "Name"                                      = "data-manager-vpc"
    "kubernetes.io/cluster/${var.cluster-name}" = "shared"
  }
}

resource "aws_subnet" "data-manager-public-subnet" {
  count                   = 2
  map_public_ip_on_launch = true
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  cidr_block              = "10.0.${20+count.index}.0/24"
  vpc_id                  = aws_vpc.data-manager-vpc.id

  tags = {
    "Name"                                      = "data-manager-public-subnet"
    "kubernetes.io/cluster/${var.cluster-name}" = "shared"
    "kubernetes.io/role/elb"                    = "1"
  }
}

resource "aws_subnet" "data-manager-private-subnet" {
  count             = 2
  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = "10.0.${10+count.index}.0/24"
  vpc_id            = aws_vpc.data-manager-vpc.id
  tags = {
    "Name"          = "data-manager-private-subnet"
  }
}

resource "aws_eip" "data-manager-eip" {
  vpc                       = true
  associate_with_private_ip = "10.0.0.213"
  tags = {
    "Name"                  = "data-manager-eip"
  }
  depends_on                = [aws_internet_gateway.data-manager-gw]
}

resource "aws_internet_gateway" "data-manager-gw" {
  vpc_id = aws_vpc.data-manager-vpc.id
  tags = {
    "Name" = "data-manager-internet-gateway"
  }
  depends_on = [aws_vpc.data-manager-vpc]
}

resource "aws_nat_gateway" "data-manager-nat-gw" {
  allocation_id = aws_eip.data-manager-eip.id
  subnet_id     = aws_subnet.data-manager-public-subnet[0].id
  tags = {
    "Name" = "data-manager-nat-gw"
  }
  depends_on    = [aws_internet_gateway.data-manager-gw]
}

resource "aws_route_table" "data-manager-nat-rt" {
  vpc_id = aws_vpc.data-manager-vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    nat_gateway_id = aws_internet_gateway.data-manager-gw.id
  }
  tags = {
    "Name" = "data-manager-nat-rt"
  }
  depends_on = [aws_internet_gateway.data-manager-gw]
}

resource "aws_route_table_association" "data-manager-nat-rta" {
  count          = 2
  subnet_id      = aws_subnet.data-manager-private-subnet[count.index].id
  route_table_id = aws_route_table.data-manager-nat-rt.id
}

resource "aws_route_table" "data-manager-public-rt" {
  vpc_id = aws_vpc.data-manager-vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.data-manager-gw.id
  }
  tags = {
    "Name" = "data-manager-public-rt"
  }
  depends_on = [aws_internet_gateway.data-manager-gw]
}

resource "aws_route_table_association" "data-manager-public-rta" {
  count          = 2
  subnet_id      = aws_subnet.data-manager-public-subnet[count.index].id
  route_table_id = aws_route_table.data-manager-public-rt.id
}

resource "aws_route_table" "data-manager-private-rt" {
  vpc_id = aws_vpc.data-manager-vpc.id
  tags = {
    "Name" = "data-manager-private-rt"
  }
  depends_on = [aws_internet_gateway.data-manager-gw]
}

resource "aws_route_table_association" "data-manager-private-rta" {
  count          = 2
  subnet_id      = aws_subnet.data-manager-private-subnet[count.index].id
  route_table_id = aws_route_table.data-manager-private-rt.id
}

resource "aws_vpc_endpoint" "data-manager-vpce" {
  vpc_id          = aws_vpc.data-manager-vpc.id
  service_name    = "com.amazonaws.us-west-2.s3"
  route_table_ids = ["${aws_route_table.data-manager-private-rt.id}"]

  tags = {
    Name = "data-manager-s3-endpoint"
  }
}

# This data source looks up the public DNS zone
data "aws_route53_zone" "datamanager_route53" {
  name         = local.envs["DOMAIN_NAME"]
  private_zone = false
}

# This creates an SSL certificate
resource "aws_acm_certificate" "cert" {
  domain_name       = local.envs["DOMAIN_NAME"]
  validation_method = "DNS"
  lifecycle {
    create_before_destroy = true
  }
  subject_alternative_names = [
    format("*.%s", local.envs["DOMAIN_NAME"])
  ]
}

# Prepare to validate the domains
resource "aws_route53_record" "cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
   # Skips the domain if it doesn't contain a wildcard
    if length(regexall("\\*\\..+", dvo.domain_name)) > 0
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = data.aws_route53_zone.datamanager_route53.zone_id
}

# This tells terraform to cause the route53 validation to happen
resource "aws_acm_certificate_validation" "certificate_validation" {
  certificate_arn         = aws_acm_certificate.cert.arn
  validation_record_fqdns = [for record in aws_route53_record.cert_validation : record.fqdn]
}

output "aws_acm_certificate_arn" {
  description = "ARN of ACM Certificate"
  value = aws_acm_certificate.cert.arn
}

resource "aws_wafv2_ip_set" "ip_allow_list" {
  name               = "datamanager_ip_allow_list"
  description        = "Data Manager IP Allow List"
  scope              = "REGIONAL"
  ip_address_version = "IPV4"
  addresses          = [local.envs["ALLOWED_IP_ADDRESS"]]
}

resource "aws_wafv2_web_acl" "wafacl" {
  name = "datamanager_waf"

  scope = "REGIONAL"

  default_action {
    block {}
  }

  rule {
    name     = "ip-allowlist"
    priority = 1

    action {
      allow {}
    }

    statement {
      ip_set_reference_statement {
        arn = aws_wafv2_ip_set.ip_allow_list.arn
      }
    }

    visibility_config {
      cloudwatch_metrics_enabled = true
      metric_name                = "AllowlistedIP"
      sampled_requests_enabled   = true
    }
  }

  visibility_config {
    cloudwatch_metrics_enabled = true
    metric_name                = "Blocked"
    sampled_requests_enabled   = true
  }
}

output "aws_wafv2_web_acl_arn" {
  description = "ARN of waf"
  value = aws_wafv2_web_acl.wafacl.arn
}

output "vpc_id" {
  description = "VPC ID"
  value = aws_vpc.data-manager-vpc.id
}

output "vpc_cidr_block" {
  description = "EKS VPC CIDR Block"
  value = aws_vpc.data-manager-vpc.cidr_block
}

output "private_subnet" {
  description = "Private Subnet"
  value = [aws_subnet.data-manager-private-subnet[0], aws_subnet.data-manager-private-subnet[1]]
}