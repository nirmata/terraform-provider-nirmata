output "agent_script" {
  description = "Nirmata agent install command"
  value       = nirmata_host_group_direct_connect.dc-host-group.curl_script
}
