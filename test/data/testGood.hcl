source "openstack" "base_image" {
  flavor              = var.vm_flavor
  floating_ip_network = var.vlan
  image_name          = var.base_image_name
  image_visibility    = var.image_visibility
  image_min_disk      = 8
  insecure            = "true"
  instance_name       = "packer"
  metadata = {
    architecture     = "x86_64"
    base             = "True"
    os               = "linux"
    os_distro        = var.distro_name
    os_version_major = var.distro_version
    build_date       = timestamp()
  }
  networks          = [var.internal_network_id]
  source_image_name = var.source_image_name
  ssh_interface     = "public_ip"
  ssh_ip_version    = "4"
  ssh_username      = var.ssh_username
}

build {
  sources = ["source.openstack.base_image"]

  provisioner "file" {
    destination = "/tmp"
    source      = "/image-build/packer_templates/base_image/configuration/"
  }

  provisioner "shell" {
    execute_command = "echo 'packer' | sudo -S sh -c '{{ .Vars }} {{ .Path }}'"
    inline = [
      "/usr/bin/env bash /tmp/configure.sh",
    ]
  }

}