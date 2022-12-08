provider "aws" {
  shared_config_files      = ["~/.aws/config"]
  shared_credentials_files = ["~/.aws/credentials"]
  profile                  = "personal"
}

data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "data-manager-eks-vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    "Name"                                      = "data-manager-eks-node"
    "kubernetes.io/cluster/${var.cluster-name}" = "shared"
  }
}

resource "aws_subnet" "data-manager-eks-subnet" {
  count = 2
  map_public_ip_on_launch = true
  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = cidrsubnet(aws_vpc.data-manager-eks-vpc.cidr_block, 8, count.index)
  vpc_id            = aws_vpc.data-manager-eks-vpc.id

  tags = {
    "Name"                                      = "data-manager-eks-node"
    "kubernetes.io/cluster/${var.cluster-name}" = "shared"
    "kubernetes.io/role/elb" = "1"
  }
}

resource "aws_internet_gateway" "data-manager-eks-gw" {
  vpc_id = aws_vpc.data-manager-eks-vpc.id

  tags = {
    "Name" = var.cluster-name
  }
  depends_on = [aws_vpc.data-manager-eks-vpc]
}

resource "aws_route_table" "data-manager-eks-rt" {
  vpc_id = aws_vpc.data-manager-eks-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.data-manager-eks-gw.id
  }
  depends_on = [aws_internet_gateway.data-manager-eks-gw]
}

resource "aws_route_table_association" "data-manager-eks-rta" {
  count = 2

  subnet_id      = aws_subnet.data-manager-eks-subnet[count.index].id
  route_table_id = aws_route_table.data-manager-eks-rt.id
}

# This data source looks up the public DNS zone
data "aws_route53_zone" "datamanager_route53" {
  name         = "datamanagertool.com"
  private_zone = false
}

# This creates an SSL certificate
resource "aws_acm_certificate" "cert" {
  domain_name       = "datamanagertool.com"
  validation_method = "DNS"
  lifecycle {
    create_before_destroy = true
  }
  subject_alternative_names = [
    "*.datamanagertool.com"
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
